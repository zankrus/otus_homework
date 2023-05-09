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
	defaultTimeout        time.Duration = time.Second * 10
	ErrWrongArgumentCount               = errors.New("wrong arguments count")
	ErrCantConnect                      = errors.New("cannot connect to telnet")
)

func init() {
	flag.DurationVar(&timeout, "timeout", defaultTimeout, "timeout duration")
}

func main() {
	flag.Parse()
	flag.DurationVar(&timeout, "timeout", defaultTimeout, "timeout duration")

	if flag.NArg() != 2 {
		os.Exit(1)
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
