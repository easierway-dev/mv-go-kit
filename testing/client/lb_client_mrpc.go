package main

import (
	"context"
	"fmt"
	"gitlab.mobvista.com/voyager/mrpc"
    "github.com/aws/aws-sdk-go/aws/ec2metadata"
    "github.com/aws/aws-sdk-go/aws/session"
	"gitlab.mobvista.com/voyager/protocols/gen/helloworld"
	"gitlab.mobvista.com/voyager/zlog"
        "os/signal"
            "gitlab.mobvista.com/mtech/tracelog"
                "go.opentelemetry.io/otel"
                    "syscall"

    "time"
)

const CONSUL_ADDR = "127.0.0.1:8500"
const SERVER_SERVICE = "lb_server_demo"
const SERVICENAME = "lb_test_test"

func client_call(ctx context.Context) {
	client, err := helloworld.NewGreeterPool(CONSUL_ADDR, SERVER_SERVICE, 20, time.Millisecond*time.Duration(200))
	if err != nil {
		fmt.Println("[fatal] mrpc client pool init err: ", err)
		panic(err)
	}

	req := &helloworld.HelloRequest{Name: "call"}
    az, _ := ec2metadata.New(session.New()).GetMetadata("placement/availability-zone")
    _ = az
	for i := 0; i < 1000; i++ {
        res, err := client.SayHello(ctx, req)
        if err != nil {
            continue
        }
        fmt.Println(req, res)

	}
}

func main() {
	logTest, _ := zlog.NewZLog(&zlog.Ops{
		MaxAge:       100,
		Path:         "",
		Format:       "",
		Level:        "warn",
		ReportCaller: false,
	})

	err := mrpc.InitRpc(CONSUL_ADDR, int64(1000), logTest)
	if err != nil {
		fmt.Println("[fatal] mrpc init err: ", err)
		panic(err)
	}
	err = mrpc.SetClientConfig("mrpc/lb_test")
	if err != nil {
		fmt.Println("[fatal] mrpc client init err: ", err)
		panic(err)
	}

       tracelog.FromConsulConfig(SERVICENAME, "127.0.0.1:8500", "/config/tracelog")
        ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
            defer cancel()
                tr := otel.Tracer("client_context")
                    ctx1, span := tr.Start(ctx, "client_call")
                    defer span.End()
    client_call(ctx1)





}
