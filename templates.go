// This file was automatically generated based on the contents of *.tmpl
// If you need to update this file, change the contents of those files
// (or add new ones) and run 'go generate'

package swaggering 

const (
defaultApiTmpl = "package singularity\n\nimport \"{{.BasePackageName}}/dtos\"\n\n{{range .Operations}}\n{{- if not .GoTypeInvalid -}}\nfunc (client *Client) {{.GoMethodName}}(\n	{{- range .Parameters}}{{.Name}} {{.GoTypePrefix -}}\n	{{if .GoModel}}dtos.{{end -}}{{.GoBaseType}}, {{end -}}\n) ({{ if not (eq .GoBaseType \"\") }}response {{.GoTypePrefix -}}\n  {{if .GoModel}}dtos.{{end -}}\n{{.GoBaseType}}, {{end}} err error) {\n	pathParamMap := map[string]interface{}{\n		{{range .Parameters -}}\n		{{if eq \"path\" .ParamType -}}\n		  \"{{.Name}}\": {{.Name}},\n	  {{- end }}\n		{{- end }}\n	}\n\n  queryParamMap := map[string]interface{}{\n		{{range .Parameters -}}\n		{{if eq \"query\" .ParamType -}}\n		  \"{{.Name}}\": {{.Name}},\n	  {{- end }}\n		{{- end }}\n	}\n\n	{{if .GoModel -}}\n	{{if eq .GoTypePrefix \"\"}}\n		response = make({{.GoTypePrefix}}dtos.{{.GoBaseType}}, 0)\n		err = client.DTORequest(&response, \"{{.Method}}\", \"{{.Path}}\", pathParamMap, queryParamMap\n		{{- if .HasBody -}}\n		, body\n		{{- end -}})\n	{{else}}\n		response = new(dtos.{{.GoBaseType}})\n		err = client.DTORequest(response, \"{{.Method}}\", \"{{.Path}}\", pathParamMap, queryParamMap\n		{{- if .HasBody -}}\n		, body\n		{{- end -}})\n	{{end}}\n	{{else if (eq .GoBaseType \"\")}}\n	_, err = client.Request(\"{{.Method}}\", \"{{.Path}}\", pathParamMap, queryParamMap\n	{{- if .HasBody -}}\n	, body\n  {{- end -}})\n	{{else if eq .GoBaseType \"string\"}}\n	resBody, err := client.Request(\"{{.Method}}\", \"{{.Path}}\", pathParamMap, queryParamMap\n	{{- if .HasBody -}}\n	, body\n  {{- end -}})\n	readBuf := bytes.Buffer{}\n	readBuf.ReadFrom(resBody)\n	response = string(readBuf.Bytes())\n	{{- end}}\n	return\n}\n{{end}}\n{{end}}\n"

defaultModelTmpl = "package dtos\n\nimport \"io\"\n\n{{range $enum := .Enums}}\ntype {{$enum.Name}} string\n\nconst (\n  {{- range $value := $enum.Values}}\n  {{$enum.Name}}{{$value}} {{$enum.Name}} = \"{{$value}}\"\n  {{- end}}\n)\n{{end}}\n\ntype {{.GoName}} struct {\n{{range $name, $prop := .Properties}}\n  {{- if $prop.GoTypeInvalid}}//{{end}}	{{$prop.GoName}} {{$prop.GoTypePrefix}}{{$prop.GoBaseType}} `json:\"{{$prop.SwaggerName}}\"`\n{{end}}\n}\n\nfunc (self *{{.GoName}}) Populate(jsonReader io.ReadCloser) (err error) {\n	err = ReadPopulate(jsonReader, self)\n	return\n}\n\nfunc (self *{{.GoName}}) FormatText() string {\n	return FormatText(self)\n}\n\nfunc (self *{{.GoName}}) FormatJSON() string {\n	return FormatJSON(self)\n}\n\ntype {{.GoName}}List []*{{.GoName}}\n\nfunc (list *{{.GoName}}List) Populate(jsonReader io.ReadCloser) (err error) {\n	return ReadPopulate(jsonReader, list)\n}\n\nfunc (list *{{.GoName}}List) FormatText() string {\n	text := []byte{}\n	for _, dto := range *list {\n		text = append(text, (*dto).FormatText()...)\n		text = append(text, \"\\n\"...)\n	}\n	return string(text)\n}\n\nfunc (list *{{.GoName}}List) FormatJSON() string {\n	return FormatJSON(list)\n}\n"

)
