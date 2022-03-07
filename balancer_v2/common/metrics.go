package balancer_common

import (
	"github.com/prometheus/client_golang/prometheus"
)

func CreateMetrics(subsystem string) (*prometheus.CounterVec, *prometheus.HistogramVec, *prometheus.HistogramVec, *prometheus.HistogramVec) {
	zoneIpCallCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "balancer",
		Subsystem: subsystem,
		Name:      "zone_ip_call_count",
	}, []string{"zone", "loc_zone", "ip", "service"})
	prometheus.MustRegister(zoneIpCallCounter)

	zoneWeightHistogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "balancer",
		Subsystem: subsystem,
		Name:      "zone_weight",
	}, []string{"zone", "loc_zone", "use_zone", "service"})
	prometheus.MustRegister(zoneWeightHistogramVec)

	ipWeightHistogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "balancer",
		Subsystem: subsystem,
		Name:      "ip_weight",
	}, []string{"ip", "loc_zone", "service"})
	prometheus.MustRegister(ipWeightHistogramVec)

	culWeightHistogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "balancer",
		Subsystem: subsystem,
		Name:      "zone_ip_cul_weight",
	}, []string{"zone", "loc_zone", "ip", "use_zone", "service"})
	prometheus.MustRegister(culWeightHistogramVec)
	return zoneIpCallCounter, zoneWeightHistogramVec, ipWeightHistogramVec, culWeightHistogramVec
}