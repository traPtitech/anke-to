package openapi

//go:generate go tool oapi-codegen --package openapi --generate server -o server.go ../docs/swagger/swagger.yaml
//go:generate go tool oapi-codegen --package openapi --generate types -o types.go ../docs/swagger/swagger.yaml
//go:generate go tool oapi-codegen --package openapi --generate spec -o spec.go ../docs/swagger/swagger.yaml
