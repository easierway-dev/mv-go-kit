package balancer_common

// balance type
const (
	RandomSelect       = 1
	RoundRobin         = 2
	WeightedRoundRobin = 3
	ConsistentHash     = 4
)

// discover type
const (
	TestingDiscover = 0
	ConsulDiscover  = 1
)

//Event
const (
	Success = 1
	Failed  = 2
	Timeout = 3
)

//Const Time
const (
	MaxTimeGap = 60
)
