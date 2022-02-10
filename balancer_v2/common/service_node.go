package balancer_common

//ServiceNode
type ServiceNode struct {
	Address string
	Host    string
	Port    int
	Zone    string
	Weight  int

	//Cur Weight
	CurWeight int
}
