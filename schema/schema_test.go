package schema_test

import (
	"testing"

	"github.com/johejo/go-openapi3-helper/schema"
)

type Test struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

type TestRef struct {
	Test Test `json:"test"`
	ID   int  `json:"id"`
}

func TestValidator_Validate(t *testing.T) {
	v, err := schema.NewValidatorFromPath("./testdata/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		schemaName string
		value      interface{}
		wantErr    bool
	}{
		{"valid struct", "Test", Test{Foo: "foo", Bar: 2}, false},
		{"invalid struct", "Test", Test{Foo: "foo", Bar: 6}, true},
		{"valid map", "Test", map[string]interface{}{"foo": "foo", "bar": 2}, false},
		{"invalid map", "Test", map[string]interface{}{"foo": "foo", "bar": 8}, true},
		{"no schema", "Nothing", nil, true},
		{"valid struct with ref", "TestRef", TestRef{Test: Test{Foo: "foo", Bar: 1}, ID: 1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.schemaName, tt.value); tt.wantErr != (err != nil) {
				t.Fatal(err)
			}
		})
	}
}
