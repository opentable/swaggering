package swaggering

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapStringString(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	ctx := NewContext("test", "github.com/test/test")

	serviceContents := `
{
  "apiVersion" : "0.4.12-SNAPSHOT",
  "swaggerVersion" : "1.2",
  "basePath" : null,
  "resourcePath" : "/api/deploys",
  "produces" : [ "application/json" ],
  "apis" : [ {
    "path" : "/api/deploys/update",
    "operations" : [ {
      "method" : "POST",
      "summary" : "Update the target active instance count for a pending deploy",
      "notes" : "",
      "type" : "SingularityDeploy",
      "nickname" : "updatePendingDeploy",
      "parameters" : [ ],
      "responseMessages" : [ {
        "code" : 400,
        "message" : "Deploy is not in the pending state pending or is not not present"
      } ]
    } ]
  } ],
  "models" : {
    "SingularityDeploy" : {
      "id" : "SingularityDeploy",
      "required" : [ "requestId", "id" ],
      "properties" : {
        "requestId" : {
          "type" : "string",
          "description" : "Singularity Request Id which is associated with this deploy."
        },
        "id" : {
          "type" : "string",
          "description" : "Singularity deploy id."
        },
        "env" : {
          "$ref" : "Map[string,string]",
          "description" : "Map of environment variable definitions."
        }
      }
    }
  }
}
	`

	ctx.IngestApi("deploys", "deploys.json", bytes.NewBufferString(serviceContents))

	ctx.Resolve()

	require.Contains(ctx.structs, "SingularityDeploy")
	dep := ctx.structs["SingularityDeploy"]
	envField := dep.findField("Env")
	require.NotNil(envField)

	assert.Equal("map[string]string", envField.TypeString())
}
