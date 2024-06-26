test-compile: build-protoc-gen-litepb
	mkdir -p ./test/generated
	rm -rf ./test/generated/*
	go mod init -C ./test/generated generated
	protoc --plugin ./bin/protoc-gen-litepb --litepb_out ./ --proto_path=./test/proto/ ./test/proto/test.proto
	#protoc --plugin ./bin/protoc-gen-litepb --litepb_out ./test/ --proto_path=./test/proto/ ./test/proto/test4.proto

test-compile-for-bench: build-protoc-gen-litepb
	mkdir -p test/bench/proto/google/
	protoc --proto_path=./test/proto/bench/ --proto_path=/usr/local/include/ --proto_path=./ --go_out test/bench/proto/google/ ./test/proto/bench/bench.proto
	mkdir -p test/bench/proto/gogo/
	protoc --proto_path=./test/proto/bench/ --proto_path=/usr/local/include/ --proto_path=./ --gogofast_out test/bench/proto/gogo/ ./test/proto/bench/bench.proto
	mkdir -p test/bench/proto/litepb/
	protoc --plugin ./bin/protoc-gen-litepb \
 		--proto_path=./test/proto/bench/ \
 		--proto_path=./ \
 		--proto_path=/usr/local/include/ \
 		--litepb_out test/bench/proto/litepb/ \
 		--litepb_opt test/proto/bench/litepb.yaml \
 		./test/proto/bench/bench.proto

build-protoc-gen-litepb:
	go build -tags debug -o ./bin/protoc-gen-litepb ./cmd/protoc-gen-litepb

build-protos: build-protoc-gen-litepb
	find ./proto -name '*.proto' | while read file; \
		do protoc \
			--plugin ./bin/protoc-gen-litepb \
			--proto_path=./proto \
			--litepb_out=./proto/ \
			--litepb_opt=proto/litepb.yaml \
			$${file} || exit 1; \
		done
uu:
	find ./proto -name '*.proto' | while read file; \
		do protoc \
			--plugin ./bin/protoc-gen-litepb \
			--proto_path=./proto \
			--go_out=./proto/ \
			--go_opt=paths=source_relative \
			$${file} || exit 1; \
		done