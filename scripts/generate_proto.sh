protoc --proto_path=./proto --go_out=proto/ --go_opt=paths=source_relative \
    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
    coordinator.proto kvs.proto status.proto node.proto slot.proto storage.proto