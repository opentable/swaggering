// This file was automatically generated based on the contents of *.tmpl
// If you need to update this file, change the contents of those files
// (or add new ones) and run 'go generate'

package swaggering 

const (
defaultApiTmpl = "package singularity\n\nimport \"{{.BasePackageName}}/dtos\"\n\n{{range .Operations}}\n{{- if not .GoTypeInvalid -}}\nfunc (client *Client) {{.GoMethodName}}(\n	{{- range .Parameters}}{{.Name}} {{.GoTypePrefix -}}\n	{{if .GoModel}}dtos.{{end -}}{{.GoBaseType}}, {{end -}}\n) ({{ if not (eq .GoBaseType \"\") }}response {{.GoTypePrefix -}}\n  {{if .GoModel}}dtos.{{end -}}\n{{.GoBaseType}}, {{end}} err error) {\n	pathParamMap := map[string]interface{}{\n		{{range .Parameters -}}\n		{{if eq \"path\" .ParamType -}}\n		  \"{{.Name}}\": {{.Name}},\n	  {{- end }}\n		{{- end }}\n	}\n\n  queryParamMap := map[string]interface{}{\n		{{range .Parameters -}}\n		{{if eq \"query\" .ParamType -}}\n		  \"{{.Name}}\": {{.Name}},\n	  {{- end }}\n		{{- end }}\n	}\n\n	{{if .GoModel -}}\n	{{if eq .GoTypePrefix \"\"}}\n		response = make({{.GoTypePrefix}}dtos.{{.GoBaseType}}, 0)\n		err = client.DTORequest(&response, \"{{.Method}}\", \"{{.Path}}\", pathParamMap, queryParamMap\n		{{- if .HasBody -}}\n		, body\n		{{- end -}})\n	{{else}}\n		response = new(dtos.{{.GoBaseType}})\n		err = client.DTORequest(response, \"{{.Method}}\", \"{{.Path}}\", pathParamMap, queryParamMap\n		{{- if .HasBody -}}\n		, body\n		{{- end -}})\n	{{end}}\n	{{else if (eq .GoBaseType \"\")}}\n	_, err = client.Request(\"{{.Method}}\", \"{{.Path}}\", pathParamMap, queryParamMap\n	{{- if .HasBody -}}\n	, body\n  {{- end -}})\n	{{else if eq .GoBaseType \"string\"}}\n	resBody, err := client.Request(\"{{.Method}}\", \"{{.Path}}\", pathParamMap, queryParamMap\n	{{- if .HasBody -}}\n	, body\n  {{- end -}})\n	readBuf := bytes.Buffer{}\n	readBuf.ReadFrom(resBody)\n	response = string(readBuf.Bytes())\n	{{- end}}\n	return\n}\n{{end}}\n{{end}}\n"

defaultModelTmpl = "package dtos\n\nimport (\n  \"fmt\"\n  \"io\"\n)\n\n\n{{range $enum := .Enums}}\ntype {{$enum.Name}} string\n\nconst (\n  {{- range $value := $enum.Values}}\n  {{$enum.Name}}{{$value}} {{$enum.Name}} = \"{{$value}}\"\n  {{- end}}\n)\n{{end}}\n\ntype {{.GoName}} struct {\n  present map[string]bool\n{{range $name, $prop := .Properties}}\n  {{- if $prop.GoTypeInvalid}}//{{end}}	{{$prop.GoName}} {{$prop.GoTypePrefix}}{{$prop.GoBaseType}} `json:\"{{$prop.SwaggerName}}\n  {{- if eq $prop.GoBaseType \"string\" -}}\n  ,omitempty\n  {{- end -}}\n  \"`\n{{end}}\n}\n\nfunc (self *{{.GoName}}) Populate(jsonReader io.ReadCloser) (err error) {\n	return ReadPopulate(jsonReader, self)\n}\n\nfunc (self *{{.GoName}}) MarshalJSON() ([]byte, error) {\n  return MarshalJSON(self)\n}\n\nfunc (self *{{.GoName}}) FormatText() string {\n	return FormatText(self)\n}\n\nfunc (self *{{.GoName}}) FormatJSON() string {\n	return FormatJSON(self)\n}\n\nfunc (self *{{.GoName}}) FieldsPresent() []string {\n  return presenceFromMap(self.present)\n}\n\nfunc (self *{{.GoName}}) SetField(name string, value interface{}) error {\n  if self.present == nil {\n    self.present = make(map[string]bool)\n  }\n  switch name {\n  default:\n    return fmt.Errorf(\"No such field %s on {{.GoName}}\", name)\n  {{range $name, $prop := .Properties}}\n    {{ if not $prop.GoTypeInvalid }}\n    case \"{{$prop.SwaggerName}}\", \"{{$prop.GoName}}\":\n      v, ok := value.({{$prop.GoTypePrefix}}{{$prop.GoBaseType}})\n      if ok {\n        self.{{$prop.GoName}} = v\n        self.present[\"{{$prop.SwaggerName}}\"] = true\n        return nil\n      } else {\n        return fmt.Errorf(\"Field {{$prop.SwaggerName}}/{{$prop.GoName}}: value %v couldn't be cast to type {{$prop.GoTypePrefix}}{{$prop.GoBaseType}}\", value)\n      }\n    {{end}}\n  {{end}}\n  }\n}\n\nfunc (self *{{.GoName}}) GetField(name string) (interface{}, error) {\n  switch name {\n  default:\n    return nil, fmt.Errorf(\"No such field %s on {{.GoName}}\", name)\n  {{range $name, $prop := .Properties}}\n    {{ if not $prop.GoTypeInvalid }}\n    case \"{{$prop.SwaggerName}}\", \"{{$prop.GoName}}\":\n    if self.present != nil {\n      if _, ok := self.present[\"{{$prop.SwaggerName}}\"]; ok {\n        return self.{{$prop.GoName}}, nil\n      }\n    }\n    return nil, fmt.Errorf(\"Field {{$prop.GoName}} no set on {{.GoName}} %+v\", self)\n    {{end}}\n  {{end}}\n  }\n}\n\nfunc (self *{{.GoName}}) ClearField(name string) error {\n  if self.present == nil {\n    self.present = make(map[string]bool)\n  }\n  switch name {\n  default:\n    return fmt.Errorf(\"No such field %s on {{.GoName}}\", name)\n  {{range $name, $prop := .Properties}}\n    {{ if not $prop.GoTypeInvalid }}\n  case \"{{$prop.SwaggerName}}\", \"{{$prop.GoName}}\":\n    self.present[\"{{$prop.SwaggerName}}\"] = false\n    {{end}}\n  {{end}}\n  }\n\n  return nil\n}\n\nfunc (self *{{.GoName}}) LoadMap(from map[string]interface{}) error {\n  return loadMapIntoDTO(from, self)\n}\n\ntype {{.GoName}}List []*{{.GoName}}\n\nfunc (list *{{.GoName}}List) Populate(jsonReader io.ReadCloser) (err error) {\n	return ReadPopulate(jsonReader, list)\n}\n\nfunc (list *{{.GoName}}List) FormatText() string {\n	text := []byte{}\n	for _, dto := range *list {\n		text = append(text, (*dto).FormatText()...)\n		text = append(text, \"\\n\"...)\n	}\n	return string(text)\n}\n\nfunc (list *{{.GoName}}List) FormatJSON() string {\n	return FormatJSON(list)\n}\n"

dtoGoTmpl = "package dtos\n\nimport (\n	\"bytes\"\n	\"encoding/json\"\n	\"errors\"\n	\"fmt\"\n	\"io\"\n	\"strings\"\n)\n\ntype DTO interface {\n	Populate(io.ReadCloser) error\n	FormatText() string\n	FormatJSON() string\n	FieldsPresent() []string\n	GetField(string) (interface{}, error)\n	SetField(string, interface{}) error\n	ClearField(string) error\n}\n\ntype StringList []string\n\nfunc (list StringList) Populate(jsonReader io.ReadCloser) (err error) {\n	return ReadPopulate(jsonReader, list)\n}\n\nfunc (list StringList) FormatText() string {\n	return strings.Join(list, \"\\n\")\n}\n\nfunc (list StringList) FormatJSON() string {\n	return FormatJSON(list)\n}\n\nfunc ReadPopulate(jsonReader io.ReadCloser, target interface{}) (err error) {\n	data := make([]byte, 0, 1024)\n	chunk := make([]byte, 1024)\n	for {\n		var count int\n		count, err = jsonReader.Read(chunk)\n		data = append(data, chunk[:count]...)\n\n		if err == io.EOF {\n			jsonReader.Close()\n			break\n		}\n		if err != nil {\n			return\n		}\n	}\n\n	if len(data) == 0 {\n		err = nil\n		return\n	}\n\n	err = json.Unmarshal(data, target)\n	return\n}\n\nfunc MarshalJSON(dto DTO) (buf []byte, err error) {\n	data := make(map[string]interface{})\n	for _, name := range dto.FieldsPresent() {\n		data[name] = dto.GetField(name)\n	}\n	return json.Marshal(data)\n}\n\nfunc presenceFromMap(m map[string]bool) []string {\n	presence := make([]string, 0)\n	for name, present := range m {\n		if present {\n			presence = append(presence, name)\n		}\n	}\n	return presence\n}\n\n\nfunc loadMapIntoDTO(from map[string]interface{}, dto DTO) error {\n	errs := make([]error)\n	for name, value := range from {\n		if err := dto.SetField(name, value); err != nil {\n			errs = append(errs, err)\n		}\n	}\n	if len(errs) > 0 {\n		return errors.New(strings.Join(errs))\n	}\n	return nil\n}\n\nfunc FormatText(dto interface{}) string {\n	return fmt.Sprintf(\"%+v\", dto)\n}\n\nfunc FormatJSON(dto interface{}) string {\n	str, err := json.Marshal(dto)\n	if err != nil {\n		return \"<<XXXX>>\"\n	} else {\n		buf := bytes.Buffer{}\n		json.Indent(&buf, str, \"\", \"  \")\n		return buf.String()\n	}\n}\n// vim: set ft=go:\n"

)
