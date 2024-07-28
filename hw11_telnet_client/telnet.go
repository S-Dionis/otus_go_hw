package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	Address    string
	Timeout    time.Duration
	Connection net.Conn
	In         io.ReadCloser
	Out        io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{
		Address:    address,
		Timeout:    timeout,
		Connection: nil,
		In:         in,
		Out:        out,
	}
}

func (t *TelnetClientImpl) Connect() error {
	conn, err := net.DialTimeout("tcp", t.Address, t.Timeout)
	if err != nil {
		return err
	}

	t.Connection = conn

	return nil
}

func (t *TelnetClientImpl) Close() error {
	return t.Connection.Close()
}

func (t *TelnetClientImpl) Send() error {
	reader := bufio.NewReader(t.In)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
			return nil
		}

		_, err = t.Connection.Write([]byte(text))
		if err != nil {
			return err
		}
	}
}

func (t *TelnetClientImpl) Receive() error {
	reader := bufio.NewReader(t.Connection)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
			return nil
		}
		_, err = t.Out.Write([]byte(text))
		if err != nil {
			return err
		}
	}
}
