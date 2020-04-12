// package schema provides helpers for openapi3 schema

package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// Validator describes openapi3 schema validator.
type Validator struct {
	swagger *openapi3.Swagger
}

// Swagger is a alias for openapi3.Swagger.
type Swagger = openapi3.Swagger

// NewValidatorFromSwagger returns a new Validator from Swagger.
func NewValidatorFromSwagger(s *Swagger) *Validator {
	return &Validator{swagger: s}
}

// NewValidatorFromPath returns a new Validator from spec path.
func NewValidatorFromPath(path string) (*Validator, error) {
	schema, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(path)
	if err != nil {
		return nil, err
	}
	return NewValidatorFromSwagger(schema), nil
}

// Validate validates value by schema name.
func (v *Validator) Validate(schemaName string, value interface{}) error {
	vv, err := toI(value)
	if err != nil {
		return errors.New("value is not a valid json")
	}
	schema, ok := v.swagger.Components.Schemas[schemaName]
	if !ok {
		return fmt.Errorf("schema %s does not exists", schemaName)
	}
	if schema.Value == nil {
		return fmt.Errorf("schema %s does not have Value", schemaName)
	}
	if err := schema.Value.VisitJSON(vv); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

func toI(value interface{}) (interface{}, error) {
	var vv interface{}
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &vv); err != nil {
		return nil, err
	}
	return vv, nil
}
