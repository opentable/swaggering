package swaggering

// DataType represents an abstract datatype as described by swagger
type DataType struct {
	GoTypePrefix, GoPackage, GoBaseType string
	GoTypeInvalid, GoModel              bool
	Type, Format                        string
	Ref                                 string `json:"$ref"`
	Enum                                []string
	EnumDesc                            Enum
}

func (dt *DataType) findGoType(context *Context) (err error) {
	if len(ds.Enum) > 0 {
		dt.EnumDesc = Enum{Name: dt.Ref, Values: dt.Enum}
	}
	return dt.setGoType(findGoType(context, dt))
}

func (dt *DataType) goPrimitiveType() (t string, err error) {
	return goPrimitiveFormattedType(self.Type, self.Format)
}

func (dt *DataType) setGoType(typeName string, err error) {
	if err != nil {
		log.Printf("Invalid type: %s %v", typeName, err)
		self.GoTypeInvalid = true
	}
	self.GoBaseType = typeName
}
