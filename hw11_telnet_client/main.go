package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout    *time.Duration
	host, port string
)

func main() {
	timeout = flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatalln("Usage: go-telnet [--timeout=10s] host port")
	}
	host = flag.Arg(0)
	port = flag.Arg(1)
	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}
	defer client.Close()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()
	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "Send error: %v\n", err)
			cancel()
		}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
			} else {
				fmt.Fprintf(os.Stderr, "Receive error: %v\n", err)
			}
			cancel()
		}
	}()
	<-ctx.Done()
	fmt.Fprintln(os.Stderr, "...Connection closed")
}
