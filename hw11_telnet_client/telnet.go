package main

import (
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

type telnetClient struct {
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
	timeout    time.Duration
	address    string
}

func (c telnetClient) Close() error {
	err := c.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c telnetClient) Send() error {
	_, err := io.Copy(c.connection, c.in)
	if err != nil {
		return err
	}
	return nil
}

func (c telnetClient) Receive() error {
	_, err := io.Copy(c.out, c.connection)
	if err != nil {
		return err
	}
	return nil
}

func (c telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.connection = conn
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return telnetClient{
		connection: nil,
		in:         in,
		out:        out,
		timeout:    timeout,
		address:    address,
	}
}
