package openapi

import (
	"context"
	"os"

	"github.com/getkin/kin-openapi/openapi3"

	"go.uber.org/zap"
)

func LoadOpenAPISpec(path string, logger *zap.Logger) (*openapi3.T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Failed to read OpenAPI spec", zap.Error(err))
		return nil, ErrFailedToReadOpenAPISpec
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromData(data)
	if err != nil {
		logger.Error("Failed to load OpenAPI spec", zap.Error(err))
		return nil, ErrFailedToLoadOpenAPISpec
	}

	if err := doc.Validate(context.Background()); err != nil {
		logger.Error("Failed to validate OpenAPI spec", zap.Error(err))
		return nil, ErrFailedToValidateOpenAPISpec
	}

	return doc, nil
}
