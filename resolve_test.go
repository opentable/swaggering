package swaggering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveModel(t *testing.T) {
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

	if !mod.GoUses {
		t.Error("Model isn't marked as used")
	}

	if testProp.GoName != "Test" {
		t.Error("Property name should be Test, was ", testProp.GoName)
	}

	if testProp.GoBaseType != "TestModelList" {
		t.Error("Property GoBaseType should be TestModelList, was ", testProp.GoBaseType)
	}

	if testProp.GoTypePrefix != "" {
		t.Error("Property GoTypePrefix should be '', was ", testProp.GoTypePrefix)
	}

	if testStr.GoBaseType != "StringList" {
		t.Error("Property GoBaseType should be StringList, was", testStr.GoBaseType)
	}
}

func TestResolveProperty(t *testing.T) {
	assert := assert.New(t)

	ctx := Context{}
	mapStrStr := Property{}
	mapStrStr.Ref = "Map[string,string]"

	ctx.resolveProperty("test", &mapStrStr)
	assert.Equal("map[string]string", mapStrStr.GoBaseType)
	assert.Equal("", mapStrStr.GoTypePrefix)
}
