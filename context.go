package swaggering

import (
	"fmt"
	"log"
)

type Context struct {
	packageName, importName string
	swaggers                []*Swagger
	openModels              []*Model

	codefiles []*CodeFile
	structs   map[string]*Struct
}

func NewContext(packageName, importName string) (context *Context) {
	context = &Context{packageName: packageName, importName: importName}
	context.swaggers = make([]*Swagger, 0)

	return
}

func (context *Context) Resolve() {

	/*
		for _, swagger := range context.swaggers {
			for adx := range swagger.Apis {
				context.apis = append(context.apis, &swagger.Apis[adx])
			}
			for _, model := range swagger.Models {
				context.models = append(context.models, model)
			}
		}

		log.Printf("  Found %d apis", len(context.apis))
		log.Printf("  Found %d models", len(context.models))

		context.openModels = make([]*Model, 0, len(context.models))
	*/

	context.resolveApis()
	context.resolveModels()
}

func (context *Context) resolveApis() {
	for _, swagger := range context.swaggers {
		file := CodeFile{}
		context.codefiles = append(context.codefiles, &file)

		file.BasePackageName = context.packageName
		file.PackageImportName = context.importName
		file.Name = swagger.name

		for _, api := range swagger.Apis {

			for _, op := range api.Operations {
				method := Method{}
				file.Methods = append(file.Methods, &method)

				method.Path = api.Path
				method.Name = capitalize(op.Nickname)

				mtype, err := op.findGoType(context)
				//context.responseType(op)
				if err != nil {
					logErr(err, "Operation %s invalid: %v", op.Nickname)
					method.invalidity = true
				}

				if mtype != nil {
					method.DTORequest = !isPrimitive(mtype)
					method.Results = append(method.Results, &Field{Name: "response", TypeStringer: mtype})
				}

				for _, parm := range op.Parameters {
					field := Field{Name: parm.Name}
					prm := Param{Field: &field, ParamType: parm.ParamType}
					method.Params = append(method.Params, &prm)

					t, err := parm.findGoType(context)
					logErr(err, "Operation %s invalid: parameter %s: %v", op.Nickname, parm.Name)

					field.TypeStringer = t

					if !t.Valid() {
						method.invalidity = true
					}

					if parm.Name == "body" {
						method.HasBody = true
					}
				}
			}
		}
	}
}

func (context *Context) responseType(op *Operation) (TypeStringer, error) {
	var rms string
	for _, rm := range op.ResponseMessages {
		if rm.ResponseModel != "" {
			if rms != "" {
				log.Printf("Operation %q has multiple response models - using %q, (not %q) but probably swaggering needs to be extended.",
					op.Nickname, rms, rm.ResponseModel)
				continue
			}
			rms = rm.ResponseModel
		}
	}

	if rms != "" {
		rmst := SwaggerType{Ref: rms}
		return findGoType(context, &rmst)
	}

	return nil, nil
}

func (context *Context) modelFor(typeName string) (TypeStringer, error) {
	t, err := context.modelUsed(typeName)
	if err != nil {
		return nil, err
	}

	return &Pointer{to: t}, nil
}

func (context *Context) modelUsed(name string) (TypeStringer, error) {
	for _, swagger := range context.swaggers {
		for _, model := range swagger.Models { // XXX it's a map - simply use Models[name] ?
			if model.Id == name {
				if !model.resolved {
					context.openModels = append(context.openModels, model)
				}

				return context.getStruct(name)
			}
		}
	}
	return nil, fmt.Errorf("model %q doesn't appear in known models", name)
}

func (context *Context) getStruct(name string) (*Struct, error) {
	if context.structs == nil {
		context.structs = map[string]*Struct{}
	}
	if s, has := context.structs[name]; has {
		return s, nil
	}
	context.structs[name] = &Struct{
		Package: "dtos",
		Name:    name,
	}
	return context.structs[name], nil
}

func (context *Context) resolveModels() {
	var cur *Model

	for len(context.openModels) > 0 {
		cur, context.openModels = context.openModels[0], context.openModels[1:]
		if cur.resolved {
			continue
		}
		context.resolveModel(cur)
	}
}

func (context *Context) resolveModel(model *Model) *Struct {
	model.resolved = true

	s, err := context.getStruct(model.Id)
	logErr(err, "when getting struct by name: %q: %v", model.Id)

	for name, prop := range model.Properties {
		field, err := context.resolveProperty(name, prop)
		logErr(err, "when resolving property type: %v")

		if field == nil {
			continue
		}
		attr := Attribute{
			Field:       field,
			SwaggerName: prop.SwaggerName,
		}
		s.Fields = append(s.Fields, &attr)

		switch enum := field.TypeStringer.(type) {
		case nil:
		case *EnumType:
			enum.HostModel = model.Id

			exists := false
			for _, e := range s.Enums {
				if e.Name == enum.Name {
					exists = true
					break
				}
			}

			if !exists {
				s.Enums = append(s.Enums, enum)
			}
		}
	}

	return s
}

func (context *Context) resolveProperty(name string, prop *Property) (*Field, error) {
	t, err := prop.findGoType(context) //uses embedded Collection's impl
	if err != nil {
		return nil, err
	}
	return &Field{Name: capitalize(name), TypeStringer: t}, nil
}
