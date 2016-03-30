package swaggering

import "testing"

func TestResolveModel(t *testing.T) {
	mod := Model{}
	mod.Id = "TestModel"
	testProp := Property{}
	mod.Properties = make(map[string]*Property)
	mod.Properties["test"] = &testProp

	testProp.Type = "array"
	testProp.Items.Ref = "TestModel"

	testStr := Property{}
	mod.Properties["testStr"] = &testStr
	testStr.Type = "array"
	testStr.Items.Type = "string"

	ctx := Context{}
	ctx.models = append(ctx.models, &mod)

	resolveModel(&mod, &ctx)

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
