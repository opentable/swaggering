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
