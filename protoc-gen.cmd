protoc --proto_path=api/proto/v1 --proto_path=./ --proto_path=third_party --go-grpc_out=pkg/api/proto_gen user_service.proto

protoc --proto_path=api/proto/v1 --go_out=pkg/api/proto_gen user_message.proto