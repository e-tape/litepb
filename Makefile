test-compile: build-protoc-gen-litepb
	mkdir -p ./test/generated
	protoc --plugin ./bin/protoc-gen-litepb --litepb_out ./test/generated ./test/proto/test.proto

build-protoc-gen-litepb:
	go build -o ./bin/protoc-gen-litepb ./cmd/protoc-gen-litepb
