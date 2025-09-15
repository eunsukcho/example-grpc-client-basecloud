module example-cloud-agent

go 1.23.0

toolchain go1.23.4

require (
	example-receiver v0.0.0
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
)

replace example-receiver => ../example-receiver

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
)
