package openapi

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=openapi.gen.cfg.yml openapi.yml
//go:generate go run github.com/matryer/moq -out openapi_moq.gen.go . ServerInterface
