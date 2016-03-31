package swaggering

type Context struct {
	packageName string
	swaggers    []*Swagger
	apis        []*Api
	models      []*Model
	openModels  []*Model
}

func NewContext(packageName string) (context *Context) {
	context = &Context{packageName: packageName}
	context.swaggers = make([]*Swagger, 0)

	return
}
