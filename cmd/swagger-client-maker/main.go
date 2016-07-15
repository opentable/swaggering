package main

import (
	"fmt"
	"log"

	docopt "github.com/docopt/docopt-go"
	"github.com/opentable/swaggering"
)

func main() {
	log.SetFlags(0)
	doc := cleanWS(`
	Usage: swagger-client-maker [--client-package=<name>] <apidoc_dir> <go_project_dir>

	Options:
	  --client-package Force the client package name, instead of infering from go_project_dir

	`)

	parsed, err := docopt.Parse(doc, nil, true, "", false)
	if err != nil {
		log.Fatal(err)
	}

	serviceSource := parsed["<apidoc_dir>"].(string)
	renderTarget := parsed["<go_project_dir>"].(string)

	packageName, ok := parsed["--client-package"]
	if !ok || packageName == nil {
		packageName, err = swaggering.InferPackage(renderTarget)
		if err != nil {
			err = fmt.Errorf("%v - try using --client-package", err)
			log.Fatal(err)
		}
	}

	swaggering.Process(packageName.(string), serviceSource, renderTarget)
}
