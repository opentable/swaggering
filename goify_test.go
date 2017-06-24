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

	dt := DataType{}
	json.Unmarshal([]byte(dtJSON), &dt)

	tf, err := dt.goPrimitiveType()
	if err != nil {
		t.Fatalf("Error should be nil, was: %v", err)
	}

	if tf != "string" {
		t.Fatalf("Formatted type should be 'string' was %v", t)
	}
}

func TestGoifyMapToMap(t *testing.T) {

	dtJSON := `{
		"$ref": "Map[int,Map[string,string]]",
		"description": "Map of environment variable overrides for specific task instances."
	}`

	ctx := Context{}
	dt := DataType{}
	json.Unmarshal([]byte(dtJSON), &dt)

	err := dt.findGoType(&ctx)
	assert.NoError(t, err)
	assert.Equal(t, `map[int]map[string]string`, dt.Type)

}
