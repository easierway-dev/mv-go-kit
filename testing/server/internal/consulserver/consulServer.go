package consulserver

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
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
	ServerProperty *ServerProperty
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
	fmt.Printf("Received: %v\n", in.GetName())
	if rand.Float64() < s.ServerProperty.ErrRate {
		return &helloworld.HelloReply{s.registerConfig.StringInfo()}, errors.New("make failed")
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
	s.registerConfig = registerConfig
	fmt.Println("[applyProperty]", s.registerConfig.StringInfo())
	s.registerConfig.Deregister()
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
