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

func (self *DataType) findModel(context *Context) (err error) {
	return context.modelFor(self.Type, self)
}

func (self *Parameter) findGoType(context *Context) (err error) {
	if self.ParamType == "body" {
		err = self.findModel(context)
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
		self.findModel(context)
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
				self.GoBaseType = "StringList"
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

	if from.Type == "" {
		if isAggregate(from.Ref) {
			err = context.aggregateType(from.Ref, to)
		} else {
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
		keys, err := goPrimitiveType(matches[1], "none")
		if err != nil {
			return err
		}
		values, err := goPrimitiveType(matches[2], "none")
		if err != nil {
			return err
		}

		to.setGoType(fmt.Sprintf("map[%s]%s", keys, values), nil)
		return nil
	}

	if matches := listRE.FindStringSubmatch(typeDesc); matches != nil {
		values, err := goPrimitiveType(matches[1], "none")
		if err != nil {
			return err
		}
		to.setGoType(fmt.Sprintf("[]%s", values), nil)
	}

	return fmt.Errorf("Not recognized as an aggregate type: %s", typeDesc)
}

func (self *DataType) goPrimitiveType() (t string, err error) {
	return goPrimitiveType(self.Type, self.Format)
}

func goPrimitiveType(sType, format string) (t string, err error) {
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
