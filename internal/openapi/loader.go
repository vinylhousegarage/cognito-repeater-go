package openapi

import (
	_ "embed"
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"go.uber.org/zap"
)

var openapiSpec []byte

func LoadEmbeddedOpenAPISpec(logger *zap.Logger) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(openapiSpec)
	if err != nil {
		logger.Error("failed to load embedded OpenAPI spec", zap.Error(err))
		return nil, ErrFailedToLoadOpenAPISpec
	}
	if err := doc.Validate(context.Background()); err != nil {
		logger.Error("failed to validate embedded OpenAPI spec", zap.Error(err))
		fmt.Printf("VALIDATION ERROR: %+v\n", err)
		return nil, ErrFailedToValidateOpenAPISpec
	}
	return doc, nil
}
