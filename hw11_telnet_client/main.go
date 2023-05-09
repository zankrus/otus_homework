package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout               time.Duration
	ErrWrongArgumentCount = errors.New("wrong arguments count")
	ErrCantConnect        = errors.New("cannot connect to telnet")
)

func main() {
	flag.DurationVar(&timeout, "timeout", 10, "timeout duration")
	flag.Parse()

	if flag.NArg() != 2 {
		log.Panic(ErrWrongArgumentCount)
	}

	client := NewTelnetClient(net.JoinHostPort(flag.Arg(0), flag.Arg(1)), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Panic(ErrCantConnect)
	}
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	go func() {
		defer cancel()
		err := client.Receive()
		if err != nil {
			return
		}
	}()

	go func() {
		defer cancel()
		err := client.Send()
		if err != nil {
			return
		}
	}()

	<-ctx.Done()
}
