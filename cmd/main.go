package main

import (
	"flag"
	"log"

	"github.com/schraf/codegen/internal"
)

func main() {
	var filename string

	flag.StringVar(&filename, "project", "codegen.json", "codegen project configuration file")
	flag.Parse()

	project, err := internal.LoadProject(filename)
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}

	if err := project.Execute(); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
}
