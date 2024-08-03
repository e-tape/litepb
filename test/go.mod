module test

go 1.22.4

replace github.com/e-tape/litepb => ../

require (
	github.com/google/uuid v1.6.0
	google.golang.org/protobuf v1.34.2
)

require golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
