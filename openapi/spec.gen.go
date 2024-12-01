// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
