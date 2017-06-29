package swaggering

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveModel(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)

	assert := assert.New(t)
	require := require.New(t)

	mod := Model{}
	mod.Id = "TestModel"
	mod.Properties = make(map[string]*Property)

	testProp := Property{}
	testProp.Type = "array"
	testProp.Items.Ref = "TestModel"
	mod.Properties["test"] = &testProp

	testStr := Property{}
	testStr.Type = "array"
	testStr.Items.Type = "string"
	mod.Properties["testStr"] = &testStr

	swg := Swagger{
		Models: map[string]*Model{
			"TestModel": &mod,
		},
	}

	ctx := Context{
		swaggers:   []*Swagger{&swg},
		openModels: []*Model{&mod},
	}
	//ctx.models = append(ctx.models, &mod)

	strct := ctx.resolveModel(&mod)

	testArray := strct.findField("Test")
	testString := strct.findField("TestStr")

	require.NotNil(testArray)
	require.NotNil(testString)

	assert.Equal("dtos.TestModelList", testArray.TypeString(""))
	assert.Equal("TestModelList", testArray.TypeString("dtos"))
	assert.Equal("swaggering.StringList", testString.TypeString(""))
}

func TestResolveProperty_Maps(t *testing.T) {
	assert := assert.New(t)

	ctx := Context{}
	mapStrStr := Property{}
	mapStrStr.Ref = "Map[string,string]"

	prop, err := ctx.resolveProperty("test", &mapStrStr)
	if assert.NoError(err) {
		assert.Equal("map[string]string", prop.TypeString(""))
	}
}

func TestResolveProperty_DeepMaps(t *testing.T) {
	assert := assert.New(t)

	ctx := Context{}
	m := Property{}
	m.Ref = "Map[string,Map[string,string]]"

	prop, err := ctx.resolveProperty("test", &m)
	if assert.NoError(err) {
		assert.Equal("map[string]map[string]string", prop.TypeString(""))
	}
}

func TestResolveProperty_ListOfModels(t *testing.T) {
	assert := assert.New(t)

	mod := Model{}
	mod.Id = "Thing"
	mod.Properties = make(map[string]*Property)

	swagger := Swagger{Models: map[string]*Model{}}
	swagger.Models["Thing"] = &mod

	ctx := Context{}
	ctx.swaggers = append(ctx.swaggers, &swagger)

	mapStrStr := Property{}
	mapStrStr.Ref = "List[Thing]"

	prop, err := ctx.resolveProperty("test", &mapStrStr)
	if assert.NoError(err) {
		assert.Equal("dtos.ThingList", prop.TypeString(""))
	}
}

func TestResolveProperty_Enum(t *testing.T) {
	assert := assert.New(t)

	mod := Model{}
	mod.Id = "Thing"
	mod.Properties = make(map[string]*Property)

	swagger := Swagger{Models: map[string]*Model{}}
	swagger.Models["Thing"] = &mod

	ctx := Context{}
	ctx.swaggers = append(ctx.swaggers, &swagger)

	enum := Property{}
	enum.Ref = "EnumKind"
	enum.Enum = []string{"A", "B", "C"}

	other := Property{}
	other.Ref = "EnumKind"
	other.Enum = []string{"A", "B", "C"}

	mod.Properties["enummy"] = &enum
	mod.Properties["other"] = &enum

	strct := ctx.resolveModel(&mod)

	f := strct.findField("Enummy")
	if assert.NotNil(f) {
		//assert.Equal(false, enum.GoTypeInvalid)
		assert.Equal("ThingEnumKind", f.TypeString(""))
	}

	assert.Equal(1, len(strct.Enums))
	assert.Equal("ThingEnumKind", strct.Enums[0].TypeString(""))
}

func TestResolveOperation_GetPendingDeploys(t *testing.T) { // from Singularity
	glbcJSON := `
		{
			"method": "GET",
			"summary": "Retrieve the list of current pending deploys",
			"notes": "",
			"type": "array",
			"items": {
				"$ref": "SingularityPendingDeploy"
			},
			"nickname": "getPendingDeploys",
			"parameters": []
		} `

	op := Operation{}
	json.Unmarshal([]byte(glbcJSON), &op)

	mod := Model{}
	mod.Id = "SingularityPendingDeploy"
	mod.Properties = make(map[string]*Property)

	swagger := Swagger{Models: map[string]*Model{"SingularityPendingDeploy": &mod}}

	ctx := Context{swaggers: []*Swagger{&swagger}}
	method := ctx.resolveOperation(&op)

	assert.True(t, method.HasResult())
	assert.Equal(t, "dtos.SingularityPendingDeployList", method.ResultTypeString(""))
}

func TestResolveOperation_GetLBCleanup(t *testing.T) { // from Singularity
	glbcJSON := `
		{
			"method": "GET",
			"summary": "Retrieve the list of tasks being cleaned from load balancers.",
			"notes": "",
			"items": {
				"type": "string"
			},
			"nickname": "getLbCleanupRequests",
			"parameters": [
				{
					"name": "useWebCache",
					"required": false,
					"type": "boolean",
					"paramType": "query",
					"allowMultiple": false
				}
			]
		}
	`

	op := Operation{}
	json.Unmarshal([]byte(glbcJSON), &op)

	ctx := Context{}
	method := ctx.resolveOperation(&op)

	assert.True(t, method.Valid())
	assert.True(t, method.HasResult())
	assert.False(t, method.HasBody)
	assert.Equal(t, method.Method, "GET")
	assert.Equal(t, method.Name, "GetLbCleanupRequests")
	assert.Equal(t, method.ResultTypeString(""), "swaggering.StringList")
	assert.Len(t, method.Params, 1)
	wc := method.Params[0]
	assert.Equal(t, wc.Name, "useWebCache")
	assert.Equal(t, wc.TypeString(""), "bool")
	assert.Equal(t, wc.ParamType, "query")
}

func TestResolveModel_SingularityDockerInfo(t *testing.T) { // from Singularity
	sdiJSON := `
	{
		"id": "SingularityDockerInfo",
		"required": [
			"image",
			"privileged"
		],
		"properties": {
			"image": {
				"type": "string",
				"description": "Docker image name"
			},
			"privileged": {
				"type": "boolean",
				"description": "Controls use of the docker --privleged flag"
			},
			"network": {
				"$ref": "SingularityDockerNetworkType",
				"description": "Docker netowkr type. Value can be BRIDGE, HOST, or NONE",
				"enum": [
					"HOST",
					"BRIDGE",
					"NONE"
				]
			},
			"portMappings": {
				"type": "array",
				"description": "List of port mappings",
				"items": {
					"$ref": "SingularityDockerPortMapping"
				}
			},
			"forcePullImage": {
				"type": "boolean",
				"description": "Always run docker pull even if the image already exists locally"
			},
			"parameters": {
				"$ref": "Map[string,string]"
			},
			"dockerParameters": {
				"type": "array",
				"description": "Other docker run command line options to be set",
				"items": {
					"$ref": "SingularityDockerParameter"
				}
			}
		}
	}
	`

	mod := Model{}
	err := json.Unmarshal([]byte(sdiJSON), &mod)
	assert.NoError(t, err)

	ctx := Context{}
	strct := ctx.resolveModel(&mod)

	assert.True(t, strct.Valid())
	assert.Equal(t, strct.Name, "SingularityDockerInfo")
	assert.Equal(t, strct.Package, "dtos")
	assert.Len(t, strct.Fields, 7)
	assert.Len(t, strct.Enums, 1)

}
