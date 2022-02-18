package discover

import (
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

type TestingDiscover struct {
	discoverNode string
	lastIndex    uint64
	nodes        []*balancer_common.ServiceNode
	interval     time.Duration
	notify       DiscoverNotify

	logger balancer_common.Logger
}

func NewTestingDiscover(interval time.Duration, notify DiscoverNotify, logger balancer_common.Logger) (Discover, error) {
	discover := &TestingDiscover{
		interval: interval,
		notify:   notify,
		logger:   logger,
	}
	//start timer
	if err := discover.Start(); err != nil {
		return nil, err
	}
	return discover, nil
}

func NewTestServiceNode() []*api.ServiceEntry {
	service1 := &api.AgentService{
		Address: "192.168.1.1",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "local_zone", "__weight": "120"},
	}
	service2 := &api.AgentService{
		Address: "192.168.1.2",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "local_zone", "__weight": "100"},
	}
	service3 := &api.AgentService{
		Address: "192.168.1.3",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "local_zone", "__weight": "100"},
	}
	service4 := &api.AgentService{
		Address: "10.0.0.1",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "other_zone1", "__weight": "100"},
	}
	service5 := &api.AgentService{
		Address: "10.0.2.3",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "other_zone2", "__weight": "100"},
	}

	entrys := []*api.ServiceEntry{}
	entrys = append(entrys,
		&api.ServiceEntry{Service: service1},
		&api.ServiceEntry{Service: service2},
		&api.ServiceEntry{Service: service3},
		&api.ServiceEntry{Service: service4},
		&api.ServiceEntry{Service: service5})
	return entrys
}

func (discover *TestingDiscover) Start() error {
	go func() {
		entrys := NewTestServiceNode()
		discover.UpdateNodes(entrys)
		ticker := time.NewTicker(discover.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				discover.UpdateNodes(entrys)
			}
		}
	}()
	return nil
}

func (discover *TestingDiscover) UpdateNodes(services []*api.ServiceEntry) {
	//get nodes
	nodes := make([]*balancer_common.ServiceNode, 0, len(services))
	for _, service := range services {
		if service.Service.Address == "" || service.Service.Port == 0 {
			continue
		}
		zone := "empty"
		if zoneStr, ok := service.Service.Meta["__zone_id"]; ok {
			zone = zoneStr
		}
		weight := 100
		if weightStr, ok := service.Service.Meta["__weight"]; ok {
			weightInt, err := strconv.Atoi(weightStr)
			if err == nil && weightInt > 0 {
				weight = weightInt
			}
		}
		node := &balancer_common.ServiceNode{
			Address: net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port)),
			Host:    service.Service.Address,
			Port:    service.Service.Port,
			Zone:    zone,
			Weight:  weight,
		}
		nodes = append(nodes, node)
	}
	discover.nodes = nodes
	//notify
	if discover.notify != nil {
		discover.notify.UpdateServicesNotify(nodes)
	}
}
