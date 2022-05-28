package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"word-of-wisdom/internal/server"
)

const defaultServerAddr = "127.0.0.1:8081"

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	rand.Seed(time.Now().UnixNano())

	addr := defaultServerAddr
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	if err := server.Run(ctx, addr); err != nil {
		fmt.Printf("failed to run server: %v", err)
	}

	return

}
