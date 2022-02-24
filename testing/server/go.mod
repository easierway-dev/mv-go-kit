module main

go 1.15

require internal/helloworld v1.0.0

replace internal/helloworld => ./internal/helloworld

require internal/resource v1.0.0

replace internal/resource => ./internal/resource

require internal/consulutils v1.0.0

replace internal/consulutils => ./internal/consulutils

replace google.golang.org/grpc => google.golang.org/grpc v1.29.0

require (
	github.com/BurntSushi/toml v1.0.0
	github.com/hashicorp/consul/api v1.12.0
	github.com/pkg/errors v0.8.1
	google.golang.org/grpc v1.25.1
)
