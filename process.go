package swaggering

func Process(packageName, serviceSource, renderTarget string) {
	context := NewContext(packageName)

	ProcessService(serviceSource, context)
	ResolveService(context)
	RenderService(renderTarget, context)
}
