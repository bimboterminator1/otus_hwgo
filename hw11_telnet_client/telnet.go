package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *TelnetClient {
	return &TelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *TelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	tc.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", tc.address)
	return nil
}

func (tc *TelnetClient) Close() error {
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}

func (tc *TelnetClient) Send() error {
	scanner := bufio.NewScanner(tc.in)
	for scanner.Scan() {
		_, err := tc.conn.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return err
		}
	}
	if scanner.Err() == nil {
		return io.EOF
	}
	return scanner.Err()
}

func (tc *TelnetClient) Receive() error {
	reader := bufio.NewReader(tc.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(tc.out, line)
		if err != nil {
			return err
		}
	}
}
