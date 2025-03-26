PHONY: genproto

genproto:
	@protoc -I pkg/pb/proto/ pkg/pb/proto/oten-test.proto \
		--go_out=pkg/pb/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/pb/proto \
		--go-grpc_opt=paths=source_relative
