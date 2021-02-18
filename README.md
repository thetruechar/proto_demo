# install:
1. go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0 
2. go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
3. go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.2.0
4. go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.2.0



# generate

1. protoc -I proto --go_out ./sdk --go-grpc_out ./sdk proto/*.proto  
2. protoc -I proto --grpc-gateway_out ./sdk --grpc-gateway_opt logtostderr=true \
--grpc-gateway_opt allow_repeated_fields_in_body=true \
proto/*.proto
3. protoc -I proto --openapiv2_out ./openapi --openapiv2_opt logtostderr=true \
--openapiv2_opt allow_repeated_fields_in_body=true \
--openapiv2_opt allow_merge=true --openapiv2_opt merge_file_name=proto_demo \
proto/*.proto
