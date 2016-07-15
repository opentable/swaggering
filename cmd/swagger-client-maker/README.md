# swagger-client-maker

Executable to generate swagger 1.2 clients
(For Swagger 2.0, please see https://github.com/go-swagger/go-swagger)

This program makes use of https://github.com/opentable/swaggering -
if you need more functionality, check it out.

As a code generator, we're assuming some familiarity with Go.

# Install

     > go get github.com/opentable/swagger-client-maker

# Usage

Get your `swagger.json` and friends into an `api-docs/` directory.

```
> swagger-client-maker api-docs/ .
```

The complete signature is `swagger-client-maker <source dir> <target dir>`

# Notes

c.f. the swaggering library for details
