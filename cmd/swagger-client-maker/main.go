package main

import (
	"fmt"
	"log"
	"path/filepath"

	docopt "github.com/docopt/docopt-go"
	"github.com/opentable/swaggering"
)

func main() {
	log.SetFlags(0)
	doc := cleanWS(`
	Usage: swagger-client-maker [options] <apidoc_dir> <go_project_dir>

	Options:
	  --import-name=<name>    Force the import name of the package (i.e. 'import github.com/import/name'), instead of infering from go_project_dir
	  --client-package=<name>  Force the client package name (i.e. 'package client-package'), instead of infering from go_project_dir

	`)

	parsed, err := docopt.Parse(doc, nil, true, "", false)
	if err != nil {
		log.Fatal(err)
	}

	serviceSource := parsed["<apidoc_dir>"].(string)
	renderTarget := parsed["<go_project_dir>"].(string)

	importName, ok := parsed["--import-name"]
	if !ok || importName == nil {
		importName, err = swaggering.InferPackage(renderTarget)
		if err != nil {
			err = fmt.Errorf("%v - try using --client-package", err)
			log.Fatal(err)
		}
	}
	in := importName.(string)

	packageName := filepath.Base(in)
	if png := parsed["--client-package"]; png != nil {
		packageName = png.(string)
	}

	swaggering.Process(in, packageName, serviceSource, renderTarget)
}
