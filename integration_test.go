package swaggering

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update .golden files")

func TestMapStringString(t *testing.T) {
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

	require.Contains(t, ctx.structs, "SingularityDeploy")
	dep := ctx.structs["SingularityDeploy"]
	envField := dep.findField("Env")
	require.NotNil(t, envField)

	assert.Equal(t, "map[string]string", envField.TypeString(""))

	targetDir := filepath.Join(os.TempDir(), "swaggering_test")
	os.RemoveAll(targetDir)
	os.MkdirAll(targetDir, os.ModePerm)
	os.MkdirAll("testdata", os.ModePerm)
	log.Printf("Integration test dir: %s", targetDir)

	RenderService(targetDir, ctx)

	fileRendersCorrectly := func(name string) {
		path := filepath.Join(targetDir, name)
		_, err := os.Stat(path)
		require.NoError(t, err)
		actual, err := ioutil.ReadFile(path)
		require.NoError(t, err)
		subName := strings.ReplaceAll(name, string(os.PathSeparator), "_")
		golden := filepath.Join("testdata", fmt.Sprintf("%s-%s.golden", t.Name(), subName))

		if *update {
			ioutil.WriteFile(golden, actual, 0644)
		}
		expected, err := ioutil.ReadFile(golden)
		require.NoError(t, err)

		actualLines := bytes.Split(actual, []byte{'\n'})
		expectedLines := bytes.Split(expected, []byte{'\n'})

		assert.ElementsMatch(t, actualLines, expectedLines)
	}

	fileRendersCorrectly("deploys.go")
	fileRendersCorrectly("dtos/singularity_deploy.go")
}
