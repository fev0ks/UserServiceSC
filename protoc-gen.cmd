protoc --proto_path=api/proto --proto_path=./ --proto_path=third_party --go-grpc_out=pkg/api user_service.proto

protoc --proto_path=api/proto --proto_path=third_party --go_out=pkg/api user_service.proto