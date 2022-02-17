package discover

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

// ConsulDiscover
type ConsulDiscover struct {
	client       *api.Client
	discoverNode string
	lastIndex    uint64
	nodes        []*balancer_common.ServiceNode
	interval     time.Duration
	notify       DiscoverNotify
	stop         bool

	logger balancer_common.Logger
}

//new discover
func NewConsulDiscover(address string, discoverNode string,
	interval time.Duration, notify DiscoverNotify, logger balancer_common.Logger) (*ConsulDiscover, error) {
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	//init discover
	discover := &ConsulDiscover{
		client:       client,
		discoverNode: discoverNode,
		interval:     interval,
		notify:       notify,
		stop:         false,
		logger:       logger,
	}
	//start timer
	if err := discover.Start(); err != nil {
		return nil, err
	}
	return discover, nil
}

//Update Notify
func (discover *ConsulDiscover) UpdateNotify(notify DiscoverNotify) {
	discover.notify = notify
}

//Start timer
func (discover *ConsulDiscover) Start() error {
	if err := discover.updateService(); err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(discover.interval)
		defer ticker.Stop()
		for {
			if discover.stop {
				break
			}
			select {
			case <-ticker.C:
				if err := discover.updateService(); err != nil && discover.logger != nil {
					discover.logger.Warnf("update service zone failed. err: [%v]", err)
				}
			}
		}
	}()
	return nil
}

//update service
func (discover *ConsulDiscover) updateService() error {
	//get services list
	services, metainfo, err := discover.client.Health().Service(discover.discoverNode, "", true, &api.QueryOptions{
		WaitIndex:  discover.lastIndex,
		AllowStale: true,
	})
	if err != nil {
		return fmt.Errorf("error retrieving instances from Consul: %v", err)
	}
	discover.lastIndex = metainfo.LastIndex
	discover.UpdateNodes(services)
	return nil
}

func (discover *ConsulDiscover) UpdateNodes(services []*api.ServiceEntry) {
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

func (discover *ConsulDiscover) Stop() {
	discover.stop = true
}
