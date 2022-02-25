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
	github.com/gfremex/logrus-kafka-hook v0.0.0-20180109031623-f62e125fcbfe // indirect
	github.com/hashicorp/consul/api v1.12.0
	github.com/pkg/errors v0.8.1
	gitlab.mobvista.com/mtech/tracelog v1.2.2
	go.opentelemetry.io/otel v1.0.0-RC3
	google.golang.org/grpc v1.25.1
	gopkg.in/sohlich/elogrus.v7 v7.0.0 // indirect

)
