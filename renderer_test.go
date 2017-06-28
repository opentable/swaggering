package swaggering

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
)

import "testing"

func ExampleSnakeCase() {
	fmt.Println(snakeCase("HTTPtest"))
	fmt.Println(snakeCase("SomethingJSON"))
	// Output:
	// httptest
	// something_json
}

func TestRenderer(t *testing.T) {
	assert := assert.New(t)

	assert.NotPanics(func() {
		NewRenderer("nowhere")
	})
}

func TestRenderModel(t *testing.T) {
	assert := assert.New(t)

	r := NewRenderer("/tmp")

	m := &Struct{}
	m.Name = "TestModel"
	m.Fields = []*Attribute{
		{
			SwaggerName: "test",
			Field: &Field{
				Name:         "Test",
				TypeStringer: &PrimitiveType{Name: "string"},
			},
		},
	}

	bytes, err := r.renderStruct("testing", m)
	if !assert.NoError(err) {
		writeErrfile("/tmp/swaggering-test/brknModel.go", bytes)
	}
}

func TestRenderCodefile(t *testing.T) {
	assert := assert.New(t)
	r := NewRenderer("/tmp")

	a := &CodeFile{}
	a.BasePackageName = "test"
	a.Methods = []*Method{
		{
			Name:    "DtoOp",
			Method:  "GET",
			Path:    "/dto-op",
			HasBody: false,
		},
		{
			Name:   "HasResult",
			Method: "GET",
			Path:   "/dto-op",
			Results: []*Field{
				{
					Name: "response",
					TypeStringer: &Pointer{to: &Struct{
						Package: "dtos",
						Name:    "ResponseThing",
					}},
				},
			},
			HasBody: false,
		},
		{
			Name:   "HasBody",
			Method: "PUT",
			Path:   "/dto-op",
			Params: []*Param{
				{
					ParamType: "path",
					Field: &Field{
						Name:         "pathparm",
						TypeStringer: &PrimitiveType{Name: "string"},
					},
				},
				{
					ParamType: "query",
					Field: &Field{
						Name:         "dto",
						TypeStringer: &PrimitiveType{Name: "int64"},
					},
				},
				{
					ParamType: "body",
					Field: &Field{
						Name: "dto",
						TypeStringer: &Pointer{to: &Struct{
							Package: "dtos",
							Name:    "BodyThing",
						}},
					},
				},
			},
		},
	}

	bytes, err := r.renderCodeFile("testing", a)
	writeErrfile("/tmp/swaggering-test/brknApi.go", bytes)
	if !assert.NoError(err) {
		writeErrfile("/tmp/swaggering-test/brknApi.go", bytes)
	}
}

func writeErrfile(path string, bytes []byte) {
	os.MkdirAll(filepath.Dir(path), os.ModeDir|os.ModePerm)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Print(err)
		return
	}
	f.Write(bytes)
	log.Print("Code written to ", path)
	//log.Print("\n", string(bytes))
}
