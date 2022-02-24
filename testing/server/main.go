package main

import (
	"context"
	"fmt"
	"main/internal/consulserver"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()
	sm := consulserver.NewServerManager()
	go consulserver.RunTask(sm)
	sm.Serve(ctx)
	// graceful shutdown
	fmt.Println("graceful shutdown")
	sm.Clear()
}
