package consulserver

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"internal/consulutils"
	"internal/helloworld"
	"internal/resource"
	"math/rand"
	"net"
	"time"
)

var (
	KvNotFound        = errors.New("kv not found")
	GetConsulKvFailed = errors.New("get consul kv info failed")
)

const CONSULKEY = "config/loadbalancer/test/serverprototypeconfig"
const SERVICE = "lb_server_demo"

type Ops struct {
	Type     string
	Address  string
	Path     string
	Interval time.Duration
	TryTimes int
	OnChange func(value interface{}, err error) bool
}

type Server struct {
	server         *grpc.Server
	serverProperty *ServerProperty
	registerConfig *consulutils.RegisterConfig
}

func NewServer(port int) *Server {
	fmt.Println("[port]", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	ss := &Server{server: s}
	helloworld.RegisterGreeterServer(s, ss)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
			fmt.Printf("failed to serve: %v\n", err)
		}
	}()
	return ss
}

// SayHello implements api.HelloServiceServer
func (s *Server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	time.Sleep(time.Millisecond)

	if rand.Float64() < s.serverProperty.ErrRate {
		return &helloworld.HelloReply{s.registerConfig.StringInfo()}, status.Errorf(400, "fail")
	}
	return &helloworld.HelloReply{s.registerConfig.StringInfo()}, nil
}
func (s *Server) applyProperty(sp *ServerProperty) error {
	registerConfig, err := consulutils.NewRegisterConfig(consulutils.WithServiceName(SERVICE),
		consulutils.WithAddressAndPort(resource.IP(), sp.Port),
		consulutils.WithMeta("__zone_id", sp.AvailabilityZone))
	if err != nil {
		return err
	}
	s.serverProperty = sp
	s.registerConfig = registerConfig
	fmt.Println("[applyProperty]", s.registerConfig.StringInfo())
	err = s.registerConfig.Register()

	return err
}
func (s *Server) destroy() {
	// consulserver deregister
	// consulserver shutdown
	s.registerConfig.Deregister()
	s.server.Stop()
	return
}
