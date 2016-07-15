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

	m := &Model{}
	m.GoName = "TestModel"
	m.Properties = make(map[string]*Property)
	m.Properties["test"] = &Property{
		SwaggerName: "testSw",
		GoName:      "TestGo",
	}
	m.Properties["test"].GoBaseType = "string"

	bytes, err := renderCode("testing", r.modelTmpl, m)
	if !assert.NoError(err) {
		writeErrfile("/tmp/swaggering-test/brknModel.go", bytes)
	}
}

func TestRenderApi(t *testing.T) {
	assert := assert.New(t)
	r := NewRenderer("/tmp")

	a := &Api{}
	a.BasePackageName = "test"
	a.Operations = append(a.Operations, &Operation{
		Nickname:     "dtoOp",
		Method:       "GET",
		Path:         "/dto-op",
		GoMethodName: "DtoOp",
		HasBody:      false,
	})

	bytes, err := renderCode("testing", r.apiTmpl, a)
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
