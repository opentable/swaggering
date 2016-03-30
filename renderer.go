package swaggering

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/tools/imports"
)

const (
	defaultModelTmpl = `
{{- /**/ -}}
package dtos

import "io"

type {{.GoName}} struct {
{{range $name, $prop := .Properties}}
{{- if $prop.GoTypeInvalid}}//{{end}}	{{$prop.GoName}} {{$prop.GoTypePrefix}}{{$prop.GoBaseType}}
{{end}}}

func (self *{{.GoName}}) Populate(jsonReader io.ReadCloser) (err error) {
	err = ReadPopulate(jsonReader, self)
	return
}

func (self *{{.GoName}}) FormatText() string {
	return FormatText(self)
}

func (self *{{.GoName}}) FormatJSON() string {
	return FormatJSON(self)
}

type {{.GoName}}List []*{{.GoName}}

func (list {{.GoName}}List) Populate(jsonReader io.ReadCloser) (err error) {
	return ReadPopulate(jsonReader, list)
}

func (list {{.GoName}}List) FormatText() string {
	text := []byte{}
	for _, dto := range list {
		text = append(text, (*dto).FormatText()...)
		text = append(text, "\n"...)
	}
	return string(text)
}

func (list {{.GoName}}List) FormatJSON() string {
	return FormatJSON(list)
}
{{/**/ -}}
`
	defaultApiTmpl = `
{{- /**/ -}}
package client

import "github.com/opentable/sous-singularity/client/dtos"

{{range .Operations}}
{{- if not .GoTypeInvalid -}}
func (client *Client) {{.GoMethodName}}(
	{{- range .Parameters}}{{.Name}} {{.GoTypePrefix -}}
	{{if .GoModel}}dtos.{{end -}}{{.GoBaseType}}, {{end -}}
) ({{ if not (eq .GoBaseType "") }}response {{.GoTypePrefix -}}
  {{if .GoModel}}dtos.{{end -}}
{{.GoBaseType}}, {{end}} err error) {
	pathParamMap := map[string]interface{}{
		{{range .Parameters -}}
		{{if eq "path" .ParamType -}}
		  "{{.Name}}": {{.Name}},
	  {{- end }}
		{{- end }}
	}
	queryParamMap := map[string]interface{}{
		{{range .Parameters -}}
		{{if eq "query" .ParamType -}}
		  "{{.Name}}": {{.Name}},
	  {{- end }}
		{{- end }}
	}

	{{if .GoModel -}}
	{{if eq .GoTypePrefix ""}}
		response = make({{.GoTypePrefix}}dtos.{{.GoBaseType}}, 0)
		err = client.DTORequest(response, "{{.Method}}", "{{.Path}}", pathParamMap, queryParamMap
		{{- if .HasBody -}}
		, body
		{{- end -}})
	{{else}}
		response = new(dtos.{{.GoBaseType}})
		err = client.DTORequest(response, "{{.Method}}", "{{.Path}}", pathParamMap, queryParamMap
		{{- if .HasBody -}}
		, body
		{{- end -}})
	{{end}}
	{{else if (eq .GoBaseType "")}}
	_, err = client.Request("{{.Method}}", "{{.Path}}", pathParamMap, queryParamMap
	{{- if .HasBody -}}
	, body
  {{- end -}})
	{{else if eq .GoBaseType "string"}}
	resBody, err := client.Request("{{.Method}}", "{{.Path}}", pathParamMap, queryParamMap
	{{- if .HasBody -}}
	, body
  {{- end -}})
	readBuf := bytes.Buffer{}
	readBuf.ReadFrom(resBody)
	response = string(readBuf.Bytes())
	{{- end}}
	return
}
{{end}}
{{end}}
{{- /**/ -}}
	`
)

type Renderer struct {
	targetDir          string
	modelTmpl, apiTmpl *template.Template
}

func NewRenderer(tgt string) (renderer *Renderer) {
	renderer = &Renderer{targetDir: tgt}
	renderer.apiTmpl = template.Must(template.New("api").Parse(defaultApiTmpl))
	renderer.modelTmpl = template.Must(template.New("model").Parse(defaultModelTmpl))

	return
}

func RenderService(target string, ingester *Context) {
	defer log.SetFlags(log.Flags())
	log.SetFlags(log.Flags() | log.Lshortfile)
	self := NewRenderer(target)
	for _, model := range ingester.models {
		if model.GoUses {
			log.Print("Model: ", model.GoName)
			path := filepath.Join("dtos", snakeCase(model.GoName))
			self.renderModel(path, model)
		}
	}

	for _, api := range ingester.apis {
		self.renderApi(apiPath(api.Path), api)
	}
}

func apiPath(urlPath string) string {
	slashes := regexp.MustCompile("/+")
	unders := regexp.MustCompile("^_+")
	brackets := regexp.MustCompile("[}{]")
	urlPath = slashes.ReplaceAllString(urlPath, "_")
	urlPath = unders.ReplaceAllString(urlPath, "")
	urlPath = brackets.ReplaceAllString(urlPath, "")
	return urlPath
}

func snakeCase(symbol string) string {
	write := make([]byte, 0, len(symbol)*2)
	read := []byte(symbol)
	re := regexp.MustCompile("[A-Z]+")

	for {
		found := re.FindIndex(read)

		if found == nil {
			write = append(write, bytes.ToLower(read)...)
			return string(write)
		}

		write = append(write, read[0:found[0]]...)
		if found[0] != 0 {
			write = append(write, byte("_"[0]))
		}
		write = append(write, bytes.ToLower(read[found[0]:found[1]])...)
		read = read[found[1]:]
	}
}

func (self *Renderer) renderModel(path string, model *Model) {
	renderCode(self.targetDir, path, self.modelTmpl, model)
}

func (self *Renderer) renderApi(path string, api *Api) {
	renderCode(self.targetDir, path, self.apiTmpl, api)
}

func renderCode(dir, path string, tmpl *template.Template, context interface{}) (err error) {
	fullpath := filepath.Join(dir, path) + ".go"
	log.Print("Rendering to ", fullpath)
	targetBuf := bytes.Buffer{}

	err = tmpl.Execute(&targetBuf, context)
	if err != nil {
		log.Fatal("Problem rendering ", tmpl.Name, " to ", fullpath, ": ", err)
	}

	formattedBytes, err := imports.Process(fullpath, targetBuf.Bytes(), nil)
	if err != nil {
		log.Print("Problem formatting ", tmpl.Name, " to ", fullpath, ": ", err)
		formattedBytes = targetBuf.Bytes()
	}
	if len(formattedBytes) > 0 {
		targetFile, err := os.Create(fullpath)
		defer targetFile.Close()
		if err != nil {
			log.Fatal("Problem creating file ", fullpath, " to render ", tmpl.Name, "into: ", err)
		}

		targetFile.Write(formattedBytes)
	}

	return
}
