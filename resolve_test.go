package swaggering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveModel(t *testing.T) {
	assert := assert.New(t)

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

	ctx := Context{}
	ctx.models = append(ctx.models, &mod)

	ctx.resolveModel(&mod)

	assert.True(mod.GoUses)

	assert.Equal("Test", testProp.GoName)
	assert.Equal("", testProp.GoTypePrefix)
	assert.Equal("TestModelList", testProp.GoBaseType)

	assert.Equal("StringList", testStr.GoBaseType)
}

func TestResolveProperty_Maps(t *testing.T) {
	assert := assert.New(t)

	ctx := Context{}
	mapStrStr := Property{}
	mapStrStr.Ref = "Map[string,string]"

	ctx.resolveProperty("test", &mapStrStr)
	assert.Equal("map[string]string", mapStrStr.GoBaseType)
	assert.Equal("", mapStrStr.GoTypePrefix)
}

func TestResolveProperty_ListOfModels(t *testing.T) {
	assert := assert.New(t)

	mod := Model{}
	mod.Id = "Thing"
	mod.Properties = make(map[string]*Property)

	ctx := Context{}
	ctx.models = append(ctx.models, &mod)

	mapStrStr := Property{}
	mapStrStr.Ref = "List[Thing]"

	ctx.resolveProperty("test", &mapStrStr)

	assert.Equal("ThingList", mapStrStr.GoBaseType)
	assert.Equal("", mapStrStr.GoTypePrefix)
}
