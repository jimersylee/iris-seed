package socks5

import (
	. "choco-proxy/app/common"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

const (
	SOCKS5Version = 0x05
	ReservedField = 0x00
)

type Server interface {
	Run() error
}

type ServerSocks5 struct {
	IP     string
	Port   int
	Config *Config
}

type Config struct {
	AuthMethod      Method
	PassWordChecker func(username, passwd string) bool
	PassWordMap     map[string]string
	TCPTimeout      time.Duration
}

func initConfig(config *Config) error {
	if config.AuthMethod == MethodUseAndPassWD && config.PassWordMap == nil {
		return errors.New("no password checker")
	}
	var mutex sync.Mutex

	// 如果密码不为空
	if config.PassWordMap != nil && len(config.PassWordMap) > 0 {
		config.PassWordChecker = func(username, passwd string) bool {
			mutex.Lock()
			defer mutex.Unlock()
			wantPasswd, ok := config.PassWordMap[username]
			if !ok {
				return false
			}
			return wantPasswd == passwd
		}
	}
	return nil
}

func (s *ServerSocks5) Run() error {
	// 初始化配置
	if err := initConfig(s.Config); err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", s.IP, s.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	Log.Infof("run socks5 proxy successful on %s", address)
	for {
		con, err := listen.Accept()
		if err != nil {
			Log.Fatal("connection failure from %s: %s\n", con.RemoteAddr(), err)
		}
		go func() {
			defer con.Close()
			err = s.handleConnection(con, s.Config)
			if err != nil {
				Log.Warn("handle connection error: %s\n", err)
			}
		}()
	}
}

func (s *ServerSocks5) handleConnection(con net.Conn, config *Config) error {
	// 处理 socks5 协议的协商过程
	err := auth(con, config)
	if err != nil {
		return err
	}
	// 请求过程
	return s.request(con)
}

func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait))
	}()
	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait))
	wg.Wait()

	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}

	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) {
		return err1
	}

	return nil
}

func auth(conn net.Conn, config *Config) error {
	// Read client auth message
	clientMessage, err := NewClientAuthMessage(conn)
	if err != nil {
		return err
	}

	// check if the auth method is supported
	var acceptable bool
	for _, method := range clientMessage.Methods {
		if method == config.AuthMethod {
			acceptable = true
		}
	}

	if !acceptable {
		_ = NewServerAuthMessage(conn, MethodNoAcceptable)
		return errors.New("method not supported")
	}

	if err := NewServerAuthMessage(conn, config.AuthMethod); err != nil {
		return err
	}

	if config.AuthMethod == MethodUseAndPassWD {
		cpm, err := NewClientPasswordMessage(conn)
		if err != nil {
			return err
		}

		if !config.PassWordChecker(cpm.Username, cpm.Passwd) {
			// 认证失败
			_ = WriteServerPasswdMessage(conn, PasswordAuthFailure)
			return ErrPasswordAuthFailure
		}
		if err := WriteServerPasswdMessage(conn, PasswordAuthSuccess); err != nil {
			return err
		}
	}
	return nil
}

func (s *ServerSocks5) request(conn net.Conn) error {
	message, err := NewClientRequestMessage(conn)
	if err != nil {
		return err
	}
	// 地址暂时只支持 IPv4 和 DomainName
	if CommandConnect != message.Cmd && CommandUDP != message.Cmd {
		WriteRequestFailureMessage(conn, ReplyCommandNotSupported)
		return ErrCommandNotSupported
	}

	// 暂时只支持 ipv4
	if message.AddressType == TypeIPv6 {
		WriteRequestFailureMessage(conn, ReplyAddressTypeNotSupported)
		return ErrAddressTypeNotSupported
	}

	// 处理 tcp
	if CommandConnect == message.Cmd {
		return s.handleTCP(conn, message)
	} else if CommandUDP == message.Cmd {

		return s.handleUDP()
	} else {
		return ErrCommandNotSupported
	}
}

func (s *ServerSocks5) handleUDP() error {
	return nil
}

func (s *ServerSocks5) handleTCP(conn io.ReadWriter, message *ClientRequestMessage) error {
	address := fmt.Sprintf("%s:%d", message.Address, message.Port)
	// 与目标地址建立 tcp 连接，并且设置超时时间
	targetConn, err := net.DialTimeout("tcp", address, s.Config.TCPTimeout)
	if err != nil {
		Log.Warnf("connection Erro %s", err)
		// 根据错误来写这个错误的返回类型
		WriteRequestFailureMessage(conn, ReplyConnectionRefused)
		return err
	}
	addrValue := targetConn.LocalAddr()
	addr := addrValue.(*net.TCPAddr)
	if err := WriteRequestSuccessMessage(conn, addr.IP, TypeIPv4, uint16(addr.Port)); err != nil {
		return err
	}

	// 转发过程
	return forward(conn, targetConn)
}

func forward(conn io.ReadWriter, targetConn net.Conn) error {
	//defer targetConn.Close()
	//go io.Copy(targetConn, conn)
	//_, err := io.Copy(conn, targetConn)
	//return err
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_, _ = io.Copy(targetConn, conn)
		cancel()
	}()
	go func() {
		_, _ = io.Copy(conn, targetConn)
		cancel()
	}()

	<-ctx.Done()
	return nil
}

var (
	ErrVersionNotSupported     = errors.New("protocol version not supported")
	ErrCommandNotSupported     = errors.New("request command not supported")
	ErrInvalidReservedField    = errors.New("invalid reserved field")
	ErrAddressTypeNotSupported = errors.New("address Type not supported")

	ErrMethodVersionNotSupport = errors.New("sub-negotiation method version not supported")
	ErrPasswordAuthFailure     = errors.New("error authenticating username/passewd")
)
