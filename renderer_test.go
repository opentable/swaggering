package swaggering

import (
	"fmt"

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
	err := r.renderModel("testing", m)
	assert.NoError(err)
}
