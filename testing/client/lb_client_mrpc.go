package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.mobvista.com/mtech/tracelog"
	"gitlab.mobvista.com/mtech/tracelog/logevent"
	"gitlab.mobvista.com/voyager/mrpc"
	"gitlab.mobvista.com/voyager/protocols/gen/helloworld"
	"gitlab.mobvista.com/voyager/zlog"
	"go.opentelemetry.io/otel"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"time"
)

const CONSUL_ADDR = "127.0.0.1:8500"
const SERVER_SERVICE = "lb_server_demo"
const SERVICENAME = "lb_test_test"
const CON = 10
const LOCAL_AZ_PREFIX = "6"

var (
	STAT = make([]int, 4)
	cnt  = 0
	succ = 0
	ch   = make(chan int, 1000)
)

func callDetail(serverConn, serverDetail string, err error) {
	idx := strings.Index(serverConn, ":")
    local_az := 0
    succ := 0
	if serverConn[idx+1:12] == "6" {
		local_az = 1
	}
	if err == nil {
		succ = 1
	}

	ch <- local_az*10 + succ
}

func consume() {
	for {
		select {
		case i := <-ch:
			switch i {
			case 0:
				STAT[2] += 1
			case 1:
				STAT[2] += 1
				STAT[3] += 1
			case 10:
				STAT[0] += 1
			case 11:
				STAT[0] += 1
				STAT[1] += 1

			}
		}
	}
}
func report(ctx context.Context) {
	go consume()
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			LOCAL_CNT, LOCAL_SUCC, CROSS_CNT, CROSS_SUCC := STAT[0], STAT[1], STAT[2], STAT[3]
			fmt.Println(STAT)
			STAT = []int{0, 0, 0, 0}
			m := make(map[string]string)
			m["local_az_succ"] = fmt.Sprintf("%d", LOCAL_SUCC)
			m["local_az_fail"] = fmt.Sprintf("%d", LOCAL_CNT-LOCAL_SUCC)
			m["local_az_all"] = fmt.Sprintf("%d", LOCAL_CNT)
			m["cross_az_succ"] = fmt.Sprintf("%d", CROSS_SUCC)
			m["cross_az_fail"] = fmt.Sprintf("%d", CROSS_CNT-CROSS_SUCC)
			m["cross_az_all"] = fmt.Sprintf("%d", CROSS_CNT)
			logevent.WithContext(ctx, "client_detail").WithLabelValues(m).Log("")
		}

	}

}
func client_call(ctx context.Context) {
	client, err := helloworld.NewGreeterPool(CONSUL_ADDR, SERVER_SERVICE, 200, time.Millisecond*time.Duration(200))
	if err != nil {
		fmt.Println("[fatal] mrpc client pool init err: ", err)
		panic(err)
	}

	req := &helloworld.HelloRequest{Name: "call"}
	az, _ := ec2metadata.New(session.New()).GetMetadata("placement/availability-zone")
	_ = az
	for {
		for j := 0; j < CON; j++ {
			var wg sync.WaitGroup
			go func() {
				wg.Add(1)
				defer wg.Done()
				res, err := client.SayHello(ctx, req)
				message := ""
				if err == nil {
					message = res.Message
				}
				callDetail(req.Name, message, err)
			}()
			wg.Wait()
		}
		time.Sleep(10 * time.Millisecond)

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
	go report(ctx1)
	defer span.End()
	client_call(ctx1)

}
