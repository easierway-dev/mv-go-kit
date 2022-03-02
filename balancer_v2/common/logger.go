package balancer_common

type Logger interface {
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}
