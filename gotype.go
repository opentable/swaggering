package swaggering

import "fmt"

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
		TypeString() string
		Valid() bool
	}

	invalidity  bool
	alwaysValid struct{}

	// PrimitiveType is a primitive type.
	PrimitiveType struct {
		alwaysValid
		Name string
	}

	// EnumType is an enum type.
	EnumType struct {
		alwaysValid
		Name      string
		Values    []string
		HostModel string
	}

	// MapType is an map[] type.
	MapType struct {
		alwaysValid
		keys, values TypeStringer
	}

	// SliceType is a []slice type.
	SliceType struct {
		alwaysValid
		items TypeStringer
	}

	// Struct describes a Go struct that will be build from a swagger API.
	Struct struct {
		invalidity
		Package, Name string
		Fields        []*Field
		Enums         []*EnumType
	}
	// Pointer describes a pointer.
	Pointer struct {
		alwaysValid
		to TypeStringer
	}
	/*
		// GoType represents a datatype to be rendered as Go code.
		GoType struct {
			Prefix, Package, BaseType string
			Invalid, Model            bool
			EnumDesc                  Enum
		}
	*/

	// Method describes the Go method that will be build from a swagger API.
	Method struct {
		invalidity
		Name                string
		Params              []*Field
		Results             []*Field
		Method, Path        string
		HasBody, DTORequest bool
	}

	// Field describes a Go field that will be build from a swagger API.
	Field struct {
		Name string
		Type TypeStringer
	}
)

func (v alwaysValid) Valid() bool {
	return true
}

func (v invalidity) Valid() bool {
	return !bool(v)
}

// TypeString implements TypeStringer on PrimitiveType.
func (t *PrimitiveType) TypeString() string {
	return t.Name
}

// TypeString implements TypeStringer on PrimitiveType.
func (t *EnumType) TypeString() string {
	return t.Name
}

// TypeString implements TypeStringer on MapType.
func (t *MapType) TypeString() string {
	return fmt.Sprintf("map[%s]%s", t.keys.TypeString(), t.values.TypeString())
}

// TypeString implements TypeStringer on Pointer.
func (t *Pointer) TypeString() string {
	return fmt.Sprintf("*%s", t.to.TypeString())
}

// TypeString implements TypeStringer on Struct.
func (t *Struct) TypeString() string {
	if t.Package == "" {
		return fmt.Sprintf("%s", t.Name)
	}
	return fmt.Sprintf("%s.%s", t.Package, t.Name)
}

/*
	if err == nil {
		if c.GoModel {
			c.GoTypePrefix = ""
			c.GoPackage = "dtos"
			c.GoBaseType = c.GoBaseType + "List"
		} else if c.GoBaseType == "string" {
			c.GoBaseType = "StringList"
			c.GoPackage = "swaggering"
			c.GoTypePrefix = ""
			c.GoModel = false
		} else {
			c.GoTypePrefix = "[]" + c.GoTypePrefix
		}
	}
*/
// TypeString implements TypeStringer on SliceType.
func (t *SliceType) TypeString() string {
	if isPrimitive(t.items) {
		return fmt.Sprintf("[]%s", t.items.TypeString())
	}
	return fmt.Sprintf("%sList", t.items.TypeString())
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

func (t *Struct) findField(name string) *Field {
	for _, f := range t.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// MakesResult indicates that the result value for this operation should be
// allocated in Go using make() as opposed to new()
func (method *Method) MakesResult() bool {
	if len(method.Results) == 0 {
		return false
	}
	switch method.Results[0].Type.(type) {
	default:
		return true
	case *Pointer:
		return false
	}
}
