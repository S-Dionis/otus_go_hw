package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 1000, "connection timeout")
}

func main() {
	flag.Parse()

	args := flag.Args()

	if flag.NArg() < 2 {
		fmt.Println("Usage: go run main.go [-timeout 1000]")
	}

	host := args[0]
	port := args[1]

	address := net.JoinHostPort(host, port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	defer client.Close()

	connect(client)

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for range sigCh {
			cancel()
		}
	}()

	go func() {
		if err := client.Send(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}

func connect(client TelnetClient) {
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
}
