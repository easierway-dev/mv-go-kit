package consulserver

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"internal/consulutils"
	"internal/helloworld"
	"internal/resource"
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

/*
	[ServerPrototype.az_1]
	Counts = 100
	StartPort = 44000
	AvailabilityZone = "az_1"
	ErrRate = 0.01

	[ServerPrototype.az_2]
	Counts = 100
	StartPort = 44000
	AvailabilityZone = "az_1"
	ErrRate = 0.1
*/
//
//type ServerManager struct {
//	sc      *ServersConfig
//	servers map[int]*Server // key: port
//	hashTag string          // current config tag
//	status  bool
//}
//type ServersConfig struct {
//	Servers map[string]*ServerPrototypeConfig
//	ports   []int
//	hashTag string
//}
//type ServerPrototypeConfig struct {
//	Counts           int     `toml:"Counts"`
//	StartPort        int     `toml:"StartPort"`
//	AvailabilityZone string  `toml:"AvailabilityZone"`
//	ErrRate          float64 `toml:"ErrRate"`
//}
//type ServerProperty struct {
//	Port             int
//	AvailabilityZone string  `toml:"AvailabilityZone"`
//	ErrRate          float64 `toml:"ErrRate"`
//}
type Server struct {
	server         *grpc.Server
	registerConfig *consulutils.RegisterConfig
}

//func (sc *ServersConfig) GetServerConfigs() map[int]*ServerProperty {
//	serverProperties := make(map[int]*ServerProperty)
//	for _, serverConf := range sc.Servers {
//		for i := 0; i < serverConf.Counts; i++ {
//			port := serverConf.StartPort + i
//			sp := &ServerProperty{port, serverConf.AvailabilityZone, serverConf.ErrRate}
//			serverProperties[port] = sp
//		}
//	}
//	return serverProperties
//}
//
//func GetServersConfigFromLocal() *ServersConfig {
//	a := `
//			[Servers.az_1]
//			Counts = 10
//			StartPort = 64000
//			AvailabilityZone = "az_1"
//			ErrRate = 0.01
//
//			[Servers.az_2]
//			Counts = 10
//			StartPort = 65000
//			AvailabilityZone = "az_2"
//			ErrRate = 0.1
//		`
//	var sc ServersConfig
//
//	if _, err := toml.Decode(a, &sc); err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	// 把所有的port取出来
//	ports := make([]int, 0)
//	for _, serverConf := range sc.Servers {
//		for i := 0; i < serverConf.Counts; i++ {
//			ports = append(ports, i+serverConf.StartPort)
//		}
//	}
//	sc.ports = ports
//	sc.hashTag = MD5(a)
//	return &sc
//}

//func NewServerManager() *ServerManager {
//	fmt.Println("Create ServerManager")
//	// sc := GetServersConfigFromConsul()
//	//sc := GetServersConfigFromLocal()
//	sc, _ := FromConsulConfig("47.252.4.203:8500", "/jianjilong")
//	sm, err := NewServerManagerWithConfig(sc)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	return sm
//}
//
//func (sm *ServerManager) Serve(ctx context.Context) {
//	// todo: clear all service node
//	select {
//	case <-ctx.Done():
//		return
//	}
//}
//func NewServerManagerWithConfig(sc *ServersConfig) (*ServerManager, error) {
//	fmt.Println("Create ServerManager From Config")
//	sm := &ServerManager{
//		sc:      sc,
//		servers: make(map[int]*Server),
//	}
//	err := sm.manage()
//	if err != nil {
//		fmt.Println(err)
//		return nil, err
//	}
//	return sm, nil
//}
//func (sm *ServerManager) manage() error {
//	sm.Clear()
//	ticker := time.NewTicker(time.Second)
//	sm.sync()
//	go func() {
//		for {
//			<-ticker.C
//			sm.sync()
//		}
//	}()
//	return nil
//}
//func (sm *ServerManager) Clear() {
//	fmt.Println("deregister all services", SERVICE)
//}
//
//func (sm *ServerManager) sync() {
//	/*
//		将consul的配置同步到真正的服务上
//	*/
//	fmt.Println("start sync")
//	if sm.hashTag == sm.sc.hashTag {
//		// 配置没变，啥也不干
//		return
//	}
//	// 配置有问题, 啥也不干
//	serverConfigs := sm.sc.GetServerConfigs()
//	if len(serverConfigs) == 0 {
//		sm.status = false
//		return
//	}
//
//	// create or update and register consulserver
//	for port, serverProperty := range serverConfigs {
//		// 不在servers这个map中，创建一个server
//		if _, ok := sm.servers[port]; !ok {
//			sm.servers[port] = NewServer(port)
//		}
//		consulserver := sm.servers[port]
//		if err := consulserver.applyProperty(serverProperty); err != nil {
//			fmt.Println(err)
//			sm.status = false
//			return
//		}
//	}
//	// remove and deregister consulserver
//	for port, consulserver := range sm.servers {
//		if _, ok := serverConfigs[port]; !ok {
//			consulserver.destroy()
//			delete(sm.servers, port)
//		}
//	}
//	sm.status = true
//	sm.hashTag = sm.sc.hashTag
//}

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
	time.Sleep(50 * time.Millisecond)

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
