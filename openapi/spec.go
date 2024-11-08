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
