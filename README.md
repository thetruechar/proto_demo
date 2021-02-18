# install:
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest



# generate

protoc -I proto --go_out ./sdk --go-grpc_out ./sdk proto/*.proto  

protoc -I proto --grpc-gateway_out ./sdk --grpc-gateway_opt logtostderr=true \
--grpc-gateway_opt allow_repeated_fields_in_body=true \
proto/*.proto

protoc -I proto --openapiv2_out ./openapi --openapiv2_opt logtostderr=true \
--openapiv2_opt allow_repeated_fields_in_body=true \
--openapiv2_opt allow_merge=true --openapiv2_opt merge_file_name=proto_demo \
proto/*.proto
