package balancer_common

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ZoneIpCallCounter *prometheus.CounterVec
)

func init() {
	ZoneIpCallCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "balancer",
		Subsystem: "v2",
		Name:      "zone_ip_call_count",
	}, []string{"zone", "loc_zone", "remote_ip", "service"})
	prometheus.MustRegister(ZoneIpCallCounter)
}
