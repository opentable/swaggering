package swaggering

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringGoify(t *testing.T) {
	// From Singularity "create disaster"
	dtJSON := `{
		"name": "type",
		"required": true,
		"type": "string",
		"paramType": "path",
		"allowMultiple": false,
		"enum": [
			"EXCESSIVE_TASK_LAG",
			"LOST_SLAVES",
			"LOST_TASKS",
			"USER_INITIATED"
		]
	}`

	dt := SwaggerType{}
	json.Unmarshal([]byte(dtJSON), &dt)

	tf, err := goPrimitiveType(dt.Type)
	if err != nil {
		t.Fatalf("Error should be nil, was: %v", err)
	}

	if tf.TypeString() != "string" {
		t.Fatalf("Formatted type should be 'string' was %v", t)
	}
}

func TestGoifyMapToMap(t *testing.T) {

	dtJSON := `{
		"$ref": "Map[int,Map[string,string]]",
		"description": "Map of environment variable overrides for specific task instances."
	}`

	ctx := Context{}
	dt := SwaggerType{}
	json.Unmarshal([]byte(dtJSON), &dt)

	typ, err := findGoType(&ctx, &dt)
	assert.NoError(t, err)
	if assert.NotNil(t, typ) {
		assert.Equal(t, `map[int64]map[string]string`, typ.TypeString())
	}

}
