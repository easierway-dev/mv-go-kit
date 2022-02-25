package main

import (
	"context"
	"fmt"
	"gitlab.mobvista.com/voyager/mrpc"
	"gitlab.mobvista.com/voyager/protocols/gen/helloworld"
	"gitlab.mobvista.com/voyager/zlog"
	"time"
)

const CONSUL_ADDR = "127.0.0.1:8500"
const SERVICE = "lb_server_demo"

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
	client, err := helloworld.NewGreeterPool(CONSUL_ADDR, SERVICE, 20, time.Millisecond*time.Duration(100))
	if err != nil {
		fmt.Println("[fatal] mrpc client pool init err: ", err)
		panic(err)
	}

	req := &helloworld.HelloRequest{Name: "call"}
	for i := 0; i < 1000; i++ {
		fmt.Println(client.SayHello(context.Background(), req))

	}
}
