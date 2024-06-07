test-compile: build-protoc-gen-litepb
	mkdir -p ./test/generated
	rm -rf ./test/generated/*
	go mod init -C ./test/generated generated
	protoc --plugin ./bin/protoc-gen-litepb --litepb_out ./test/ --proto_path=./test/proto/ ./test/proto/test.proto
	#protoc --plugin ./bin/protoc-gen-litepb --litepb_out ./test/ --proto_path=./test/proto/ ./test/proto/test4.proto

test-compile-google:
	protoc --proto_path=./test/proto/ --go_out ./test/ ./test/proto/test.proto

build-protoc-gen-litepb:
	go build -tags debug -o ./bin/protoc-gen-litepb ./cmd/protoc-gen-litepb

build-plugin-proto:
	protoc \
	 --proto_path=pkg/plugin\
	 --go_out ./pkg/plugin/\
	 --go_opt=paths=source_relative \
	 ./pkg/plugin/plugin.proto