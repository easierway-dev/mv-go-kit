package main

import (
	"context"
	"fmt"
	"main/internal/consulserver"
	"os/signal"
	"gitlab.mobvista.com/mtech/tracelog"
	"go.opentelemetry.io/otel"
	"syscall"
)

const SERVICENAME = "lb_test_server"
func main() {
	tracelog.FromConsulConfig(SERVICENAME, "127.0.0.1:8500", "/config/tracelog")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()
	tr := otel.Tracer("server_context")
	ctx1, span := tr.Start(ctx, "server_serve")
    defer span.End()


	sm := consulserver.NewServerManager()
	go consulserver.RunTask(sm)
	sm.Serve(ctx1)
	// graceful shutdown
	fmt.Println("graceful shutdown")
	sm.Clear()
}
