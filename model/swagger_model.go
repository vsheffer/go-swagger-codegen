// Package to parse and validate a Swagger JSON document into Go struct.
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/validator.v2"
	"reflect"
	"strings"
)

var (
	validSchemes []string = []string{
		"http",
		"https",
		"ws",
		"wss",
	}
	validFormats []string = []string{
		"int32",
		"int64",
		"float",
		"double",
		"string",
		"byte",
		"boolean",
		"date",
		"date-time",
	}
	validTypes []string = []string{
		"string",
		"number",
		"integer",
		"boolean",
		"array",
		"file",
	}
)

func isValidValue(value string, validValues []string) bool {
	for i := 0; i < len(validValues); i++ {
		if value == validValues[i] {
			return true
		}
	}
	return false
}

func validScheme(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.String:
		if isValidValue(st.String(), validSchemes) {
			return nil
		}
		return fmt.Errorf("Invalid url scheme: '%s'", st.String())
	case reflect.Slice:
		var badSchemes []string
		for i := 0; i < st.Len(); i++ {
			value := st.Index(i).String()
			if !isValidValue(value, validSchemes) {
				badSchemes = append(badSchemes, value)
			}
		}

		if len(badSchemes) > 0 {
			return fmt.Errorf("Invalid schemes: [%s]", strings.Join(badSchemes, ","))
		}
		return nil
	}
	return errors.New("validScheme only validates slices of strings or strings.")
}

func validFormat(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.String:
		if isValidValue(st.String(), validFormats) {
			return nil
		}
		return fmt.Errorf("Invalid type format: '%s'", st.String())
	}
	return errors.New("validFormat only validates slices of strings or strings.")
}

func validType(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.String:
		if isValidValue(st.String(), validTypes) {
			return nil
		}
		return fmt.Errorf("Invalid parameter type: '%s'", st.String())
	}
	return errors.New("validFormat only validates slices of strings or strings.")
}

type Swagger struct {
	Swagger     string              `json:"swagger" validate:"nonzero"`
	Info        *Info               `json:"info"`
	Host        string              `json:"host"`
	BasePath    string              `json:"basePath"`
	Schemes     []string            `json:"schemes"`
	Consumes    []string            `json:"consumes"`
	Produces    []string            `json:"produces"`
	Paths       []map[string]Paths  `json:"path" validate:"nonzero"`
	Definitions []map[string]Schema `json:"definitions"`
}

type Info struct {
	Title          string   `json:"title" validate:"nonzero"`
	Description    string   `json:"description"`
	TermsOfService string   `json:"termsOfService"`
	Contact        *Contact `json:"contact"`
	Version        string   `json:"version" validate:"nonzero"`
}

type Contact struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Email string `json:"email"`
}

type License struct {
	Name string `json:"name" validate:"required"`
	Url  string `json:"url"`
}

type Paths struct {
	Ref       string `json:"$ref"`
	Operation map[string]Operation
}

type Operation struct {
	Tags                  []string               `json:"tags"`
	Summary               string                 `json:"summary"`
	Description           string                 `json:"description"`
	ExternalDocumentation *ExternalDocumentation `json:"externalDocumentation"`
	Operationid           string                 `json:"operationId"`
	Consumes              []string               `json:"consumes"`
	Produces              []string               `json:"produces"`
	Parameters            []Parameter            `json:"parameters"`
}

type Parameter struct {
	Name             string  `json:"name" validate:"nonzero"`
	In               string  `json:"in" validate:"nonzero"`
	Description      string  `json:"description"`
	Required         bool    `json:"required"`
	Schema           *Schema `json:"schema"`
	Type             string  `json:"type" validate:"validTypes"`
	Format           string  `json:"format" validate:"validFormat"`
	Items            string  `json:"format" validate:"validType"`
	CollectionFormat string  `json:collectionFormat`
	Default          string  `json:"default"`
}

type Schema struct {
	Ref           string `json:"$ref"`
	Format        string `json:"format" validate:"validFormat"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Default       string `json:"default"`
	Discriminator string `json:"discriminator"`
	ReadOnly      bool   `json:"readOnly"`
	Xml           *Xml   `json:"xml"`
}

type Xml struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Prefix    string `json:"prefix"`
	Attribute bool   `json:"attribute"`
	Wrapped   bool   `json:"wrapped"`
}

type ExternalDocumentation struct {
	Description string `json:"description"`
	Url         string `json:"url" validate:"nonzero"`
}

func NewSwagger(jsonBytes []byte) (*Swagger, error) {
	var swagger *Swagger
	err := json.Unmarshal(jsonBytes, swagger)
	if err != nil {
		return nil, err
	}

	err = validator.SetValidationFunc("validScheme", validScheme)
	if err != nil {
		return nil, err
	}

	err = validator.SetValidationFunc("validFormat", validFormat)
	if err != nil {
		return nil, err
	}

	return swagger, nil
}
