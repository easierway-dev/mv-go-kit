package balancer_v2

import (
	"errors"
	"net"
	"regexp"
	"time"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/balancer"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/discover"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/weight_cal"
)

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
}

func NewBalancerResolver(balancerType, discoverType int, zoneName string,
	address string, discoverNode string, interval time.Duration, logger balancer_common.Logger, options ...Option) (*BalancerResolver, error) {
	//create resolver
	resolver := &BalancerResolver{}
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
	resolver.localAddress, _ = GetLocalIp()
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
	return resolver, nil
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
	//open zone cul
	useZoneCul := CheckOpenZoneWeight(nodes, resolver.localZone)
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
		balancer_common.ZoneWeightHistogramVec.WithLabelValues(node.Zone, resolver.localAddress, useZoneCulStr, resolver.discoverNode).Observe(zoneWeight)
		balancer_common.IpWeightHistogramVec.WithLabelValues(node.Address, resolver.localAddress, resolver.discoverNode).Observe(serviceWeight)
		balancer_common.CulWeightHistogramVec.WithLabelValues(node.Zone, resolver.localAddress, node.Address, useZoneCulStr, resolver.discoverNode).Observe(weight)
	}
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
	balancer_common.ZoneIpCallCounter.WithLabelValues(node.Zone, resolver.localAddress, node.Address, resolver.discoverNode).Inc()
	return node, nil
}

func (resolver *BalancerResolver) GetNode() (string, error) {
	node, err := resolver.DiscoverNode()
	if err != nil {
		return "", err
	}
	return node.Address, nil
}

func CheckOpenZoneWeight(nodes []*balancer_common.ServiceNode, localZoneName string) bool {
	localZoneNum := 0
	otherZoneNum := 0
	if len(localZoneName) != 0 {
		for _, node := range nodes {
			if localZoneName == node.Zone {
				localZoneNum += 1
			} else {
				otherZoneNum += 1
			}
		}
	}
	if localZoneNum > 0 && otherZoneNum > 0 {
		return true
	} else {
		return false
	}
}

var ip4Reg = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

func GetLocalIp() (string, error) {
	addr, err := localIPv4s()
	if err != nil {
		return "", err
	}

	if len(addr) == 0 {
		return "", errors.New("get local ip error")
	}
	return addr[0], nil
}

func localIPv4s() ([]string, error) {
	var ips []string
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}
	for _, a := range addr {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			if ip4Reg.MatchString(ipnet.IP.String()) {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}
