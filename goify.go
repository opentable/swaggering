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

func (self *Parameter) findGoType(context *Context) (err error) {
	if self.ParamType == "body" {
		err = context.modelFor(self.Type, &self.DataType)
	} else {
		err = self.DataType.findGoType(context)
	}

	return
}

func (self *Operation) findGoType(context *Context) (err error) {
	switch self.Type {
	case "void":
		self.GoBaseType = ""
	case "array":
		self.Collection.findGoType(context)
	case "":
		/* Singularity's swagger has some bugs... */
		self.Type = "array"
		self.Collection.findGoType(context)
	case "string", "bool", "integer", "number":
		typeName, err := self.goPrimitiveType()
		self.setGoType(typeName, err)
	default:
		err = context.modelFor(self.Type, &self.DataType)
	}

	return
}

func (self *Collection) findGoType(context *Context) (err error) {
	if self.Type == "array" {
		err = findGoType(context, &self.Items, &self.DataType)
		if err == nil {
			if self.GoModel {
				self.GoTypePrefix = ""
				self.GoBaseType = self.GoBaseType + "List"
			} else if self.GoBaseType == "string" {
				self.GoBaseType = "swaggering.StringList"
				self.GoTypePrefix = ""
				self.GoModel = true
			} else {
				self.GoTypePrefix = "[]" + self.GoTypePrefix
			}
		}
	} else {
		err = findGoType(context, &self.DataType, &self.DataType)
	}
	return
}

func (self *DataType) findGoType(context *Context) (err error) {
	return findGoType(context, self, self)
}

func findGoType(context *Context, from, to *DataType) (err error) {
	var typeName string

	if len(from.Enum) > 0 {
		to.EnumDesc = Enum{Name: from.Ref, Values: from.Enum}
		to.setGoType(from.Ref, nil)
	} else if from.Type == "" {
		if err = context.aggregateType(from.Ref, to); err != nil {
			err = context.modelFor(from.Ref, to)
		}
	} else {
		typeName, err = from.goPrimitiveType()
		to.setGoType(typeName, err)
	}
	return
}

func (context *Context) aggregateType(typeDesc string, to *DataType) (err error) {
	if matches := mapRE.FindStringSubmatch(typeDesc); matches != nil {
		var keys string
		keys, err = goPrimitiveType(matches[1])

		_, values, terr := goPrimitiveOrModel(context, matches[2])
		if terr != nil {
			err = terr
		}

		to.setGoType(fmt.Sprintf("map[%s]%s", keys, values), err)
		return
	}

	if matches := listRE.FindStringSubmatch(typeDesc); matches != nil {
		var values string
		var prim bool

		prim, values, err = goPrimitiveOrModel(context, matches[1])
		if prim {
			to.setGoType(fmt.Sprintf("[]%s", values), err)
		} else {
			to.setGoType(fmt.Sprintf("%sList", values), err)
		}
		return
	}

	return fmt.Errorf("Not recognized as an aggregate type: %s", typeDesc)
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

func (self *DataType) goPrimitiveType() (t string, err error) {
	return goPrimitiveFormattedType(self.Type, self.Format)
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

func (self *DataType) setGoType(typeName string, err error) {
	if err != nil {
		self.GoTypeInvalid = true
	}
	self.GoBaseType = typeName
}
