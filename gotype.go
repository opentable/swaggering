package swaggering

import (
	"fmt"
	"strings"
)

type (
	// CodeFile describes a go code file, which will contain methods
	CodeFile struct {
		Name              string
		BasePackageName   string
		PackageImportName string
		Methods           []*Method
	}

	// A TypeStringer implements TypeString, which returns the string that describes a go type.
	TypeStringer interface {
		TypeString(pkg string) string
		Valid() bool
    IsConcrete() bool
	}

	invalidity  bool
	alwaysValid struct{}

  isConcrete struct{}
  isReference struct{}

	// PrimitiveType is a primitive type.
	PrimitiveType struct {
		alwaysValid
    isConcrete
		Name string
	}

	// EnumType is an enum type.
	EnumType struct {
		alwaysValid
    isConcrete
		Name      string
		Values    []string
		HostModel string
	}

	// MapType is an map[] type.
	MapType struct {
		alwaysValid
    isConcrete
		keys, values TypeStringer
	}

	// SliceType is a []slice type.
	SliceType struct {
		alwaysValid
    isConcrete
		items TypeStringer
	}

	// Struct describes a Go struct that will be build from a swagger API.
	Struct struct {
		invalidity
    isConcrete
		Package, Name string
		Fields        []*Attribute
		Enums         []*EnumType
	}
	// Pointer describes a pointer.
	Pointer struct {
		TypeStringer
	}

	// Method describes the Go method that will be build from a swagger API.
	Method struct {
		invalidity
		hostPackage         string
		Name                string
		Params              []*Param
		Results             []*Field
		Method, Path        string
		HasBody, DTORequest bool
	}

	// Field describes a Go field that will be build from a swagger API.
	Field struct {
		Name string
		TypeStringer
	}

	// Param describes a parameter to a method.
	Param struct {
		*Field

		ParamType string
	}

	// An Attribute is a field in a Swagger Struct, which includes information about JSON serialization.
	Attribute struct {
		*Field

		SwaggerName string
	}
)

func (alwaysValid) Valid() bool {
	return true
}

func (v invalidity) Valid() bool {
	return !bool(v)
}

func (isConcrete) IsConcrete() bool {
  return true
}

func (isReference) IsConcrete() bool {
  return false
}

// TypeString implements TypeStringer on PrimitiveType.
func (t *PrimitiveType) TypeString(pkg string) string {
	return t.Name
}

// TypeString implements TypeStringer on PrimitiveType.
func (t *EnumType) TypeString(pkg string) string {
	return t.HostModel + t.Name
}

// TypeString implements TypeStringer on MapType.
func (t *MapType) TypeString(pkg string) string {
	return fmt.Sprintf("map[%s]%s", t.keys.TypeString(pkg), t.values.TypeString(pkg))
}

// TypeString implements TypeStringer on Pointer.
func (t *Pointer) TypeString(pkg string) string {
	return fmt.Sprintf("*%s", t.TypeStringer.TypeString(pkg))
}

// IsConcrete implements TypeStringer on Pointer.
func (t *Pointer) IsConcrete() bool {
  return false
}

// TypeString implements TypeStringer on Struct.
func (t *Struct) TypeString(pkg string) string {
	if t.Package == "" || t.Package == pkg {
		return fmt.Sprintf("%s", t.Name)
	}
	return fmt.Sprintf("%s.%s", t.Package, t.Name)
}

// TypeString implements TypeStringer on SliceType.
func (t *SliceType) TypeString(pkg string) string {
	if st, is := t.items.(*PrimitiveType); is {
		if st.Name == "string" {
			return "swaggering.StringList"
		}
		return fmt.Sprintf("[]%s", t.items.TypeString(pkg))
	}
	return fmt.Sprintf("%sList", t.items.TypeString(pkg))
}

// IsPrimitive implements TypeStringer on PrimitiveType.
func isPrimitive(t TypeStringer) bool {
	switch t.(type) {
	case *PrimitiveType:
		return true
	default:
		return false
	}
}

func (t *Struct) findField(name string) *Attribute {
	for _, f := range t.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// consider templatehelpers.go

// MakesResult indicates that the result value for this operation should be
// allocated in Go using make() as opposed to new()
func (method *Method) MakesResult() bool {
	if len(method.Results) == 0 {
		return false
	}
	switch method.Results[0].TypeStringer.(type) {
	default:
		return true
	case *Pointer:
		return false
	}
}

// HasResult indicates that the Method wraps an API call that produces a result value.
func (method *Method) HasResult() bool {
	return len(method.Results) > 0
}

// ResultTypeString is a shortcut for returning the typestring of a result value.
func (method *Method) ResultTypeString(pkg string) string {
	if !method.HasResult() {
		return "NO RESULT STRING"
	}
	return method.Results[0].TypeString(pkg)
}

func (method *Method) BaseResultTypeString(pkg string) string {
	switch res := method.Results[0].TypeStringer.(type) {
	default:
		return res.TypeString(pkg)
	case *Pointer:
		return res.TypeStringer.TypeString(pkg)
	}
}

func (method Method) ResourceName() string {
	return method.hostPackage + "-" + strings.ToLower(method.Name)
}

func (attr *Attribute) Omittable() bool {
	switch t := attr.TypeStringer.(type) {
	default:
		return false
	case *PrimitiveType:
		return t.Name == "string"
	}
}
