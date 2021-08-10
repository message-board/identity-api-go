#Path to .proto files
PROTO_PATH := api/proto

# Output directories.
GRPC_OUT := pkg

.PHONY: openapi
openapi:
	oapi-codegen -generate types -o internal/interfaces/rest/openapi_types.gen.go -package rest api/openapi/identity.yaml
	oapi-codegen -generate chi-server -o internal/interfaces/rest/openapi_api.gen.go -package rest api/openapi/identity.yaml

.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:internal/common/genproto/trainer -I api/proto api/proto/identity.proto