package balancer_common

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ZoneIpCallCounter      *prometheus.CounterVec
	ZoneWeightHistogramVec *prometheus.HistogramVec
	IpWeightHistogramVec   *prometheus.HistogramVec
	CulWeightHistogramVec  *prometheus.HistogramVec
)

func init() {
	ZoneIpCallCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "balancer",
		Subsystem: "v2",
		Name:      "zone_ip_call_count",
	}, []string{"zone", "loc_ip", "ip", "service"})
	ZoneWeightHistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "balancer",
		Subsystem: "v2",
		Name:      "zone_weight",
	}, []string{"zone", "loc_ip", "use_zone", "service"})
	IpWeightHistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "balancer",
		Subsystem: "v2",
		Name:      "ip_weight",
	}, []string{"ip", "loc_ip", "service"})
	CulWeightHistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "balancer",
		Subsystem: "v2",
		Name:      "zone_ip_cul_weight",
	}, []string{"zone", "loc_ip", "ip", "use_zone", "service"})
}
