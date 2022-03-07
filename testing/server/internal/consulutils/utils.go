package consulutils

import (
	"fmt"
)

type RegisterConfig struct {
	service string
	address string
	port    int
	meta    map[string]string
}

type RegisterOption func(*RegisterConfig)

// 设置serviceName
func WithServiceName(service string) RegisterOption {
	return func(c *RegisterConfig) {
		c.service = service
	}
}

// 设置address port
func WithAddressAndPort(address string, port int) RegisterOption {
	return func(c *RegisterConfig) {
		c.address = address
		c.port = port
	}
}

// 设置meta
func WithMeta(key, value string) RegisterOption {
	return func(c *RegisterConfig) {
		if len(c.meta) == 0 {
			c.meta = make(map[string]string)
		}
		c.meta[key] = value
	}
}

func (c *RegisterConfig) valid() error {
	return nil
}

func NewRegisterConfig(opts ...RegisterOption) (*RegisterConfig, error) {
	var c RegisterConfig
	for _, opt := range opts {
		opt(&c)
	}
	return &c, c.valid()
}

func (c *RegisterConfig) ID() string {
	return fmt.Sprintf("%s@%s:%d", c.service, c.address, c.port)
}
func (c *RegisterConfig) StringInfo() string {
	return fmt.Sprintf("%s@%s:%d\t%s\t%s", c.service, c.address, c.port, c.meta["__zone_id"], c.meta["__weight"])
}

func (c *RegisterConfig) Register() error {
	fmt.Printf("%s %v:%v register!\n", c.service, c.address, c.port)
	registerService(c)
	return nil
}

func (c *RegisterConfig) Deregister() error {
	fmt.Printf("%s %v:%v deregister!\n", c.service, c.address, c.port)
	deregisterService(c)
	return nil
}
