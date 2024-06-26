module bench

go 1.22.4

replace (
	github.com/e-tape/litepb => ../..
	//github.com/golang/protobuf => ./proto/litepb/github.com/golang/protobuf
)

require (
	github.com/e-tape/litepb v0.0.0-00010101000000-000000000000
	github.com/envoyproxy/protoc-gen-validate v1.0.4
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.4
	google.golang.org/protobuf v1.34.2
)
