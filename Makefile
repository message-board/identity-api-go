#Path to .proto files
PROTO_PATH := api/proto

# Output directories.
GRPC_OUT := pkg

.PHONY: openapi
openapi:
	oapi-codegen -generate types -o ports/openapi_types.gen.go -package ports api/openapi/identity.yml
	oapi-codegen -generate chi-server -o ports/openapi_api.gen.go -package ports api/openapi/trainings.yml

.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:internal/common/genproto/trainer -I api/proto api/proto/identity.proto