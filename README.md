# Swaggering
Swagger 1.2 code generator

(For Swagger 2.0, please see https://github.com/go-swagger/go-swagger)

Swaggering is a library for generating code based on Swagger 1.2 JSON description files.

You'll very likely want to look at opentable/swagger-client-maker to actually make use of this code,
but you're encouraged to review and reuse this library.


# Building

Note that the default templates live in defaultApi.tmpl and defaultModel.tmpl -
if you update those files, you need to run `go generate` to update templates.go.
