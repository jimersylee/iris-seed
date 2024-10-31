package socks5

import (
	"io"
)

// https://datatracker.ietf.org/doc/html/rfc1928
type ClientAuthMessage struct {
	Version  byte
	NMethods byte
	Methods  []Method
}

type ClientPasswordMessage struct {
	Username string
	Passwd   string
}

type Method = byte

const (
	MethodNoAuth       Method = 0x00
	MethodGSSAPI       Method = 0x01
	MethodUseAndPassWD Method = 0x02
	MethodNoAcceptable Method = 0xff
)

const (
	UsernameMethodVersion = 0x01
	PasswordAuthSuccess   = 0x00
	PasswordAuthFailure   = 0x01
)

func NewClientAuthMessage(conn io.Reader) (*ClientAuthMessage, error) {
	// Read version and nMethods
	buff := make([]byte, 2)
	// 一定读取 buff 的字节数
	_, err := io.ReadFull(conn, buff)
	if err != nil {
		return nil, err
	}

	// 验证 version
	if buff[0] != SOCKS5Version {
		return nil, ErrVersionNotSupported
	}
	nmethods := buff[1]
	buff = make([]byte, nmethods)
	_, err = io.ReadFull(conn, buff)
	if err != nil {
		return nil, err
	}

	return &ClientAuthMessage{
		Version:  SOCKS5Version,
		NMethods: nmethods,
		Methods:  buff,
	}, nil
}

func NewServerAuthMessage(conn io.Writer, method Method) error {
	buff := []byte{SOCKS5Version, method}
	_, err := conn.Write(buff)
	return err
}

func NewClientPasswordMessage(conn io.Reader) (*ClientPasswordMessage, error) {
	// read version and username length
	buff := make([]byte, 2)
	if _, err := io.ReadFull(conn, buff); err != nil {
		return nil, err
	}

	version, usernameLen := buff[0], buff[1]
	if version != UsernameMethodVersion {
		return nil, ErrMethodVersionNotSupport
	}

	// read username, passwd length 1 个字节
	buff = make([]byte, usernameLen+1)
	if _, err := io.ReadFull(conn, buff); err != nil {
		return nil, err
	}
	username, passwdLen := string(buff[:len(buff)-1]), buff[len(buff)-1]

	buff = make([]byte, passwdLen)
	if _, err := io.ReadFull(conn, buff[:passwdLen]); err != nil {
		return nil, err
	}

	passwd := string(buff)

	return &ClientPasswordMessage{
		Username: username,
		Passwd:   passwd,
	}, nil
}

func WriteServerPasswdMessage(conn io.Writer, status byte) error {
	_, err := conn.Write([]byte{MethodUseAndPassWD, status})
	return err
}
