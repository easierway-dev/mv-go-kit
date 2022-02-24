package consulserver

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type ServersConfig struct {
	Servers map[string]*ServerPrototypeConfig
	ports   []int
	hashTag string
}
type ServerPrototypeConfig struct {
	Counts           int     `toml:"Counts"`
	StartPort        int     `toml:"StartPort"`
	AvailabilityZone string  `toml:"AvailabilityZone"`
	ErrRate          float64 `toml:"ErrRate"`
}
type ServerProperty struct {
	Port             int
	AvailabilityZone string  `toml:"AvailabilityZone"`
	ErrRate          float64 `toml:"ErrRate"`
}

func (sc *ServersConfig) GetServerConfigs() map[int]*ServerProperty {
	serverProperties := make(map[int]*ServerProperty)
	for _, serverConf := range sc.Servers {
		for i := 0; i < serverConf.Counts; i++ {
			port := serverConf.StartPort + i
			sp := &ServerProperty{port, serverConf.AvailabilityZone, serverConf.ErrRate}
			serverProperties[port] = sp
		}
	}
	return serverProperties
}

func GetServersConfigFromLocal() *ServersConfig {
	a := `
			[Servers.az_1]
			Counts = 10
			StartPort = 34000
			AvailabilityZone = "az_1"
			ErrRate = 0.01
			
			[Servers.az_2]
			Counts = 10
			StartPort = 65000
			AvailabilityZone = "az_2"
			ErrRate = 0.1
		`
	var sc ServersConfig

	if _, err := toml.Decode(a, &sc); err != nil {
		fmt.Println(err)
		return nil
	}
	// 把所有的port取出来
	ports := make([]int, 0)
	for _, serverConf := range sc.Servers {
		for i := 0; i < serverConf.Counts; i++ {
			ports = append(ports, i+serverConf.StartPort)
		}
	}
	sc.ports = ports
	sc.hashTag = MD5(a)
	return &sc
}

// 通过传入服务名,consul地址,consul_key,获取config配置
func FromConsulConfig(consul_addr string, consul_key string) (*ServersConfig, error) {
	// 定义ServersConfig配置
	var sc ServersConfig
	// 初始化Ops配置，传入配置文件格式,consul地址,key
	tomlFormat := &Ops{Type: "toml", Address: consul_addr, Path: consul_key}
	// 获取配置文件并初始化ServersConfig
	getTomlConfig(tomlFormat, &sc)
	// 把所有的port取出来
	ports := make([]int, 0)
	for _, serverConf := range sc.Servers {
		for i := 0; i < serverConf.Counts; i++ {
			ports = append(ports, i+serverConf.StartPort)
		}
	}
	sc.ports = ports
	return &sc, nil
}
