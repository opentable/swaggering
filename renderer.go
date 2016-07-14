package swaggering

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/imports"
)

//go:generate inlinefiles . templates.go
//go:generate inlinefiles --vfs=Templates . vfs_templates.go

// Renderer is responsible for actually turning the parsed API into go code
type Renderer struct {
	targetDir                   string
	modelTmpl, apiTmpl, dtoTmpl *template.Template
}

// NewRenderer creates a new renderer
func NewRenderer(tgt string) (renderer *Renderer) {
	renderer = &Renderer{targetDir: tgt}
	renderer.apiTmpl = template.Must(loadTemplate(Templates, "api", "defaultApi.tmpl"))
	renderer.modelTmpl = template.Must(loadTemplate(Templates, "model", "defaultModel.tmpl"))
	renderer.apiTmpl = template.Must(loadTemplate(Templates, "dto", "dtoGo.tmpl"))

	return
}

// RenderService performs the rendering of a service
func RenderService(target string, ingester *Context) {
	self := NewRenderer(target)
	self.renderHandlers(ingester)
	for _, model := range ingester.models {
		if model.GoUses {
			log.Print("Model: ", model.GoName)
			path := filepath.Join("dtos", snakeCase(model.GoName))
			self.renderModel(path, model)
		}
	}

	for _, api := range ingester.apis {
		self.renderAPI(apiPath(api.Path), api)
	}
}

func loadTemplate(fs vfs.Opener, tName, fName string) (*template.Template, error) {
	tmplFile, err := fs.Open(fName)
	if err != nil {
		return nil, err
	}
	tmplB := &bytes.Buffer{}
	_, err = tmplB.ReadFrom(tmplFile)
	if err != nil {
		return nil, err
	}
	log.Print("\n", tmplB.String())
	return template.New(tName).Parse(tmplB.String())
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

func (r *Renderer) renderHandlers(ctx *Context) {
	renderCode(r.targetDir, "dtos/dto", r.dtoTmpl, ctx)
}

func (r *Renderer) renderModel(path string, model *Model) {
	renderCode(r.targetDir, path, r.modelTmpl, model)
}

func (r *Renderer) renderAPI(path string, api *Api) {
	renderCode(r.targetDir, path, r.apiTmpl, api)
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
