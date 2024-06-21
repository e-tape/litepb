test-compile: build-protoc-gen-litepb
	mkdir -p ./test/generated
	rm -rf ./test/generated/*
	go mod init -C ./test/generated generated
	protoc --plugin ./bin/protoc-gen-litepb --litepb_out ./ --proto_path=./test/proto/ ./test/proto/test.proto
	#protoc --plugin ./bin/protoc-gen-litepb --litepb_out ./test/ --proto_path=./test/proto/ ./test/proto/test4.proto

test-compile-for-bench: build-protoc-gen-litepb
	mkdir -p test/bench/proto/litepb/
	protoc --plugin ./bin/protoc-gen-litepb --proto_path=./test/proto/ --litepb_out test/bench/proto/litepb/ ./test/proto/bench/bench.proto
	mkdir -p test/bench/proto/google/
	protoc --proto_path=./test/proto/ --go_out test/bench/proto/google/ ./test/proto/bench/bench.proto
	mkdir -p test/bench/proto/gogo/
	protoc --proto_path=./test/proto/ --gogofast_out test/bench/proto/gogo/ ./test/proto/bench/bench.proto

build-protoc-gen-litepb:
	go build -tags debug -o ./bin/protoc-gen-litepb ./cmd/protoc-gen-litepb

build-plugin-proto:
	protoc \
	 --proto_path=pkg/plugin\
	 --go_out ./pkg/plugin/\
	 --go_opt=paths=source_relative \
	 ./pkg/plugin/plugin.proto