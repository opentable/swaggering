package swaggering

import (
	"fmt"
	"regexp"
	"strings"
)

func capitalize(word string) string {
	firstRE := regexp.MustCompile(`^.`)
	return firstRE.ReplaceAllStringFunc(word, func(match string) string {
		return strings.ToTitle(match)
	})
}

func goName(name string) string {
	if name == "Object" {
		return "interface{}"
	}
	return name
}

var listRE *regexp.Regexp = regexp.MustCompile(`^List\[([^,]*)]`)
var mapRE *regexp.Regexp = regexp.MustCompile(`^Map\[([^,]*),([^,]*)]`)

func isAggregate(kind string) bool {
	return mapRE.FindStringSubmatch(kind) != nil || listRE.FindStringSubmatch(kind) != nil
}

func (p *Parameter) findGoType(context *Context) (err error) {
	if p.ParamType == "body" {
		err = context.modelFor(p.Type, &p.DataType)
	} else {
		err = p.DataType.findGoType(context)
	}

	return
}

func findGoType(context *Context, from *DataType) (string, error) {
	var typeName string

	switch {
	default:
		return from.goPrimitiveType()
	case len(from.Enum) > 0:
		return from.Ref, nil
	case from.Type == "":
		return refType(context, from.Ref)
	}
	return
}

func refType(context *Context, typeDesc string) (string, error) {
	t, err := context.aggregateType(from.Ref, to)
	if err != nil {
		return from.Ref, context.modelFor(from.Ref, to)
	}
	return t, err
}

func aggregateItemType(context *Context, typeDesc string) (string, error) {
	t, err := context.aggregateType(from.Ref, to)
	if err != nil {
		_, t, err := goPrimitiveOrModel(context, typeDesc)
		return t, err
	}
	return t, err
}

func (context *Context) aggregateType(typeDesc string) (string, error) {
	if matches := mapRE.FindStringSubmatch(typeDesc); matches != nil {
		var keys string
		keys, err = goPrimitiveType(matches[1])

		_, values, terr := goPrimitiveOrModel(context, matches[2])
		if terr != nil {
			err = terr
		}

		return fmt.Sprintf("map[%s]%s", keys, values), err
	}

	if matches := listRE.FindStringSubmatch(typeDesc); matches != nil {
		var values string
		var prim bool

		prim, values, err = goPrimitiveOrModel(context, matches[1])
		if prim {
			return fmt.Sprintf("[]%s", values), err
		}
		return fmt.Sprintf("%sList", values), err
	}

	return "", fmt.Errorf("Not recognized as an aggregate type: %s", typeDesc)
}

func goPrimitiveOrModel(context *Context, name string) (prim bool, t string, err error) {
	t, err = goPrimitiveType(name)
	if err == nil {
		prim = true
		return
	}

	t = name
	err = context.modelUsed(name)

	return
}

func goPrimitiveType(sType string) (t string, err error) {
	switch sType {
	default:
		err = fmt.Errorf("Unrecognized primitive type: %s", sType)
	case "boolean":
		t = "bool"
	case "integer":
		t = "int64"
	case "number":
		t = "float64"
	case "string":
		t = "string"
	}
	return
}

func goPrimitiveFormattedType(sType, format string) (t string, err error) {
	switch sType {
	default:
		err = fmt.Errorf("Unrecognized primitive type: %s", sType)
	case "boolean":
		t = "bool"
	case "integer":
		t = format
	case "number":
		switch format {
		case "float", "none":
			t = "float32"
		case "double":
			t = "float64"
		default:
			err = fmt.Errorf("Invalid number format: %s", format)
		}
	case "string":
		switch format {
		case "", "byte", "none":
			t = "string"
		case "date", "data-time":
			t = "time.Time"
		default:
			err = fmt.Errorf("Invalid string format: %s", format)
		}
	}
	return
}
