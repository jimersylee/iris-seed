package socks5

import (
	"encoding/binary"
	"io"
	"net"
)

type ClientRequestMessage struct {
	Cmd         Command
	AddressType AddressType
	Address     string
	Port        uint16
}

type Command = byte

const (
	CommandConnect Command = 0x01 // tcp
	CommandBind    Command = 0x02
	CommandUDP     Command = 0x03
)

type AddressType = byte

const (
	TypeIPv4   AddressType = 0x01
	TypeDomain AddressType = 0x03
	TypeIPv6   AddressType = 0x04
)

/*
REP    Reply field:
             o  X'00' succeeded
             o  X'01' general SOCKS server failure
             o  X'02' connection not allowed by ruleset
             o  X'03' Network unreachable
             o  X'04' Host unreachable
             o  X'05' Connection refused
             o  X'06' TTL expired
             o  X'07' Command not supported
             o  X'08' Address type not supported
             o  X'09' to X'FF' unassigned
*/

type ReplyType = byte

const (
	ReplySuccess ReplyType = iota
	ReplyServerFailure
	ReplyConnectionNotAllowed
	ReplyNetworkUnreachable
	ReplyHostUnreachable
	ReplyConnectionRefused
	ReplyTTLExpired
	ReplyCommandNotSupported
	ReplyAddressTypeNotSupported
)

func NewClientRequestMessage(conn io.Reader) (*ClientRequestMessage, error) {
	// Read Version, CMD, RSV, AddressType
	buff := make([]byte, 4)
	_, err := io.ReadFull(conn, buff)
	if err != nil {
		return nil, err
	}
	version, command, reserved, addrType := buff[0], buff[1], buff[2], buff[3]
	if version != SOCKS5Version {
		return nil, ErrVersionNotSupported
	}
	if command != CommandConnect && command != CommandBind && command != CommandUDP {
		return nil, ErrCommandNotSupported
	}
	if reserved != ReservedField {
		return nil, ErrInvalidReservedField
	}
	if addrType != TypeIPv4 && addrType != TypeDomain && addrType != TypeIPv6 {
		return nil, ErrAddressTypeNotSupported
	}

	message := &ClientRequestMessage{
		Cmd:         command,
		AddressType: addrType,
	}

	// 解析 address
	switch addrType {
	case TypeIPv6:
		buff = make([]byte, 6)
		// 继续执行
		fallthrough
	case TypeIPv4:
		// 复用了 buff，因为 ipv4 也是 4 个字节
		if _, err := io.ReadFull(conn, buff); err != nil {
			return nil, err
		}
		ip := net.IP(buff)
		message.Address = ip.String()
	case TypeDomain:
		// 读取 1 个字节, 表示域名的长度
		if _, err := io.ReadFull(conn, buff[:1]); err != nil {
			return nil, err
		}
		domainLength := buff[0]
		domainBuff := make([]byte, domainLength)
		if _, err := io.ReadFull(conn, domainBuff); err != nil {
			return nil, err
		}
		message.Address = string(domainBuff)
	}
	// 读取剩下两个字节，表示端口
	if _, err := io.ReadFull(conn, buff[:2]); err != nil {
		return nil, err
	}
	message.Port = binary.BigEndian.Uint16(buff[:2])
	return message, nil
}

func WriteRequestSuccessMessage(conn io.Writer, ip net.IP, addressType AddressType, port uint16) error {
	if _, err := conn.Write([]byte{SOCKS5Version, ReplySuccess, ReservedField, addressType}); err != nil {
		return err
	}

	// Write bind IP(IPv4/IPv6)
	if _, err := conn.Write(ip); err != nil {
		return err
	}

	// Write bind port
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, port)
	_, err := conn.Write(buf)
	return err
}

func WriteRequestFailureMessage(conn io.Writer, replyType ReplyType) error {
	_, err := conn.Write([]byte{SOCKS5Version, replyType, ReservedField, TypeIPv4, 0, 0, 0, 0, 0})
	return err
}
