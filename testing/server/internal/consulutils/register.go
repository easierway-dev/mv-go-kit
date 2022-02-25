package consulutils

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
)

var (
	consulClient           *api.Client
	ErrConsulClientNotInit = errors.New("consul client not init")
	ErrConflictService     = errors.New("service conflict")
	ErrConsulAddrConflict  = errors.New("consul addr is conflict")
	ADDR                   = "127.0.0.1:8500"
)

func init() {
	fmt.Println("init consul client")
	conf := api.DefaultConfig()
	conf.Address = ADDR
	client, err := api.NewClient(conf)
	if err != nil {
		fmt.Println("consul client init failed")
		panic(err)
	}
	consulClient = client
}

func registerService(c *RegisterConfig) error {
	if consulClient == nil {
		return ErrConsulClientNotInit
	}
	service := &api.AgentServiceRegistration{
		ID:      c.ID(),
		Name:    c.service,
		Tags:    []string{"lb_test"},
		Address: c.address,
		Port:    c.port,
		Meta:    c.meta,
		Check: &api.AgentServiceCheck{
			Name:                           c.service + " tcp checker",
			TCP:                            fmt.Sprintf("%s:%d", c.address, c.port),
			Interval:                       "5s",
			Timeout:                        "10s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}
	err := consulClient.Agent().ServiceRegister(service)
	fmt.Println("consul register:", service)
	if err != nil {
		fmt.Println(err)
    	return err
	}
	return nil
}

func deregisterService(c *RegisterConfig) error {
	return consulClient.Agent().ServiceDeregister(c.ID())
}
