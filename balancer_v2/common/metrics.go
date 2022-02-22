package balancer_common

import (
	"github.com/prometheus/client_golang/prometheus"
)

func CreateMetrics(sybsystem string) (*prometheus.CounterVec, *prometheus.HistogramVec, *prometheus.HistogramVec, *prometheus.HistogramVec) {
	zoneIpCallCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "balancer",
		Subsystem: sybsystem,
		Name:      "zone_ip_call_count",
	}, []string{"zone", "loc_ip", "ip", "service"})
	zoneWeightHistogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "balancer",
		Subsystem: sybsystem,
		Name:      "zone_weight",
	}, []string{"zone", "loc_ip", "use_zone", "service"})
	ipWeightHistogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "balancer",
		Subsystem: sybsystem,
		Name:      "ip_weight",
	}, []string{"ip", "loc_ip", "service"})
	culWeightHistogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "balancer",
		Subsystem: sybsystem,
		Name:      "zone_ip_cul_weight",
	}, []string{"zone", "loc_ip", "ip", "use_zone", "service"})
	return zoneIpCallCounter, zoneWeightHistogramVec, ipWeightHistogramVec, culWeightHistogramVec
}
