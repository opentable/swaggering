package swaggering

import (
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

	assert.Equal("*dtos.TestModelList", testArray.Type.TypeString())
	assert.Equal("swaggering.StringList", testString.Type.TypeString())
}

func TestResolveProperty_Maps(t *testing.T) {
	assert := assert.New(t)

	ctx := Context{}
	mapStrStr := Property{}
	mapStrStr.Ref = "Map[string,string]"

	prop, err := ctx.resolveProperty("test", &mapStrStr)
	if assert.NoError(err) {
		assert.Equal("map[string]string", prop.Type.TypeString())
	}
}

func TestResolveProperty_DeepMaps(t *testing.T) {
	assert := assert.New(t)

	ctx := Context{}
	m := Property{}
	m.Ref = "Map[string,Map[string,string]]"

	prop, err := ctx.resolveProperty("test", &m)
	if assert.NoError(err) {
		assert.Equal("map[string]map[string]string", prop.Type.TypeString())
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
		assert.Equal("dtos.ThingList", prop.Type.TypeString())
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
		assert.Equal("ThingEnumKind", f.Type.TypeString())
	}

	assert.Equal(1, len(strct.Enums))
	assert.Equal("ThingEnumKind", strct.Enums[0].TypeString())
}
