package model

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestValidation(t *testing.T) {
	b, err := ioutil.ReadFile("./test_swagger.json")

	if err != nil {
		t.Fatalf("Can't open test file")
	}
	swagger, err := NewSwagger(b)
	if err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			t.Fatalf("Error : %+v, offset = %d", e, e.Offset)
		case *json.InvalidUnmarshalError:
			t.Fatalf("Error : %+v, type = %s", e, e.Type.String())
		}
	}

	t.Log("swagger = %+v", swagger)
}
