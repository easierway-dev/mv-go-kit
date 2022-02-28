package balancer_v2

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/balancer"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/discover"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/weight_cal"
)

type BalancerMetrics struct {
	ZoneIpCallCounter      *prometheus.CounterVec
	ZoneWeightHistogramVec *prometheus.HistogramVec
	IpWeightHistogramVec   *prometheus.HistogramVec
	CulWeightHistogramVec  *prometheus.HistogramVec
}

type BalancerResolver struct {
	discover discover.Discover //service discover
	balancer balancer.Balancer //balancer

	discoverNode string
	localZone    string

	localAddress string

	zoneAdjuster    *weight_cal.WeightAdjuster
	serviceAdjuster *weight_cal.WeightAdjuster

	serviceStep float64
	zoneStep    float64
	beta        float64

	lastUpdateTime int64
	mutex          sync.Mutex
	nodes          []*balancer_common.ServiceNode

	interval        time.Duration
	balancerMetrics BalancerMetrics
}

func NewBalancerResolver(balancerType, discoverType int, zoneName string, address string,
	discoverNode string, interval time.Duration, logger balancer_common.Logger, subsystem string, options ...Option) (*BalancerResolver, error) {
	//create resolver
	resolver := &BalancerResolver{
		serviceStep: 0.02,
		zoneStep:    0.05,
		beta:        0.9,
	}
	//init metrics
	resolver.InitMetrics(subsystem)
	//init options
	for _, option := range options {
		option(resolver)
	}
	//init serviceStep
	if resolver.serviceStep < 0.01 || resolver.serviceStep > 0.5 {
		resolver.serviceStep = 0.02
	}
	//init zone step
	if resolver.zoneStep < 0.01 || resolver.zoneStep > 0.5 {
		resolver.zoneStep = 0.05
	}
	//init local address
	resolver.localAddress, _ = balancer_common.GetLocalIp()
	//create balancer
	balancer, err := balancer.NewBalancer(balancerType, zoneName, discoverNode)
	if err != nil {
		return nil, err
	}
	resolver.discoverNode = discoverNode
	resolver.localZone = zoneName
	resolver.balancer = balancer
	//create zone && service adjuster
	resolver.zoneAdjuster = weight_cal.NewWeightAdjuster(resolver.beta)
	resolver.zoneAdjuster.ClearEmptyCounter(time.Duration(2*60*60) * time.Second)
	resolver.serviceAdjuster = weight_cal.NewWeightAdjuster(resolver.beta)
	resolver.serviceAdjuster.ClearEmptyCounter(time.Duration(2*60*60) * time.Second)
	//create discover
	discover, err := discover.NewDiscover(discoverType, address, discoverNode, interval, resolver, logger)
	if err != nil {
		return nil, err
	}
	resolver.discover = discover
	//Start
	resolver.interval = interval
	resolver.Start()
	return resolver, nil
}

func (resolver *BalancerResolver) Start() {
	go func() {
		ticker := time.NewTicker(resolver.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now().Unix()
				if now-resolver.lastUpdateTime > int64(resolver.interval) {
					resolver.UpdateServicesNotify(resolver.nodes)
				}
			}
		}
	}()
}

func (resolver *BalancerResolver) Notify(address string, zone string, event int) {
	if resolver.zoneAdjuster != nil {
		resolver.zoneAdjuster.Notify(zone, event)
	}
	if resolver.serviceAdjuster != nil {
		resolver.serviceAdjuster.Notify(address, event)
	}
}

func (resolver *BalancerResolver) UpdateServicesNotify(nodes []*balancer_common.ServiceNode) {
	//lock
	resolver.mutex.Lock()
	defer resolver.mutex.Unlock()
	//update lastUpdateTime
	resolver.lastUpdateTime = time.Now().Unix()
	//open zone cul
	useZoneCul := balancer_common.CheckOpenZoneWeight(nodes, resolver.localZone)
	useZoneCulStr := "0"
	if useZoneCul {
		useZoneCulStr = "1"
	}
	//cal CurWeight
	for _, node := range nodes {
		//cul zone Weight
		serviceWeight := weight_cal.GetServiceWeight(resolver.serviceAdjuster, node.Address, resolver.serviceStep)
		weight := float64(node.Weight) * serviceWeight
		zoneWeight := weight_cal.GetZoneWeight(resolver.zoneAdjuster, resolver.localZone, node.Zone, resolver.zoneStep)
		if useZoneCul {
			weight *= zoneWeight
		}
		node.CurWeight = int(weight)
		//local zone min weight = 1
		if node.CurWeight == 0 && node.Zone == resolver.localZone {
			node.CurWeight = 1
		}
		//add metrics
		resolver.balancerMetrics.ZoneWeightHistogramVec.WithLabelValues(node.Zone, resolver.localAddress, useZoneCulStr, resolver.discoverNode).Observe(zoneWeight)
		resolver.balancerMetrics.IpWeightHistogramVec.WithLabelValues(node.Address, resolver.localAddress, resolver.discoverNode).Observe(serviceWeight)
		resolver.balancerMetrics.CulWeightHistogramVec.WithLabelValues(node.Zone, resolver.localAddress, node.Address, useZoneCulStr, resolver.discoverNode).Observe(weight)
	}
	//set nodes
	resolver.nodes = nodes
	//update weight_cal
	if resolver.balancer != nil {
		resolver.balancer.UpdateServices(nodes)
	}
}

func (resolver *BalancerResolver) DiscoverNode() (*balancer_common.ServiceNode, error) {
	node, err := resolver.balancer.DiscoverNode()
	if err != nil {
		return nil, err
	}
	//add metrics
	resolver.balancerMetrics.ZoneIpCallCounter.WithLabelValues(node.Zone, resolver.localAddress, node.Address, resolver.discoverNode).Inc()
	return node, nil
}

func (resolver *BalancerResolver) GetNode() (string, error) {
	node, err := resolver.DiscoverNode()
	if err != nil {
		return "", err
	}
	return node.Address, nil
}

func (resolver *BalancerResolver) InitMetrics(subsystem string) {
	resolver.balancerMetrics.ZoneIpCallCounter, resolver.balancerMetrics.ZoneWeightHistogramVec,
		resolver.balancerMetrics.IpWeightHistogramVec, resolver.balancerMetrics.CulWeightHistogramVec = balancer_common.CreateMetrics(subsystem)
}
