package balancer_v2

type Option func(*BalancerResolver)

func Beta(beta float64) Option {
	return func(resolver *BalancerResolver) {
		resolver.beta = beta
	}
}

func ZoneStep(zoneStep float64) Option {
	return func(resolver *BalancerResolver) {
		resolver.zoneStep = zoneStep
	}
}

func ServiceStep(serviceStep float64) Option {
	return func(resolver *BalancerResolver) {
		resolver.serviceStep = serviceStep
	}
}

func OpenZoneWeight(openZoneWeight bool) Option {
	return func(resolver *BalancerResolver) {
		resolver.openZoneWeight = openZoneWeight
	}
}
