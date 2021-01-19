package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

type ProjectOutput struct {
	Template string
	Input    string
	Output   string
}

type ProjectConfig struct {
	Includes []string
	Outputs  []ProjectOutput
}

var verbose bool
var config ProjectConfig

func main() {
	var err error
	var filename string
	var buffer []byte
	var baseTemplate *template.Template

	flag.BoolVar(&verbose, "verbose", false, "turns on verbose logging")
	flag.StringVar(&filename, "project", "codegen.proj", "project configuration file")
	flag.Parse()

	if buffer, err = ioutil.ReadFile(filename); err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(buffer, &config); err != nil {
		log.Fatal(err)
	}

	if len(config.Includes) > 0 {
		if baseTemplate, err = template.ParseFiles(config.Includes...); err != nil {
			log.Fatal(err)
		}
	}

	for _, output := range config.Outputs {
		var outputTemplate *template.Template
		var outfile *os.File
		var input map[string]interface{}

		if buffer, err = ioutil.ReadFile(output.Input); err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(buffer, &input); err != nil {
			log.Fatal(err)
		}

		if outputTemplate, err = baseTemplate.Clone(); err != nil {
			log.Fatal(err)
		}

		if buffer, err = ioutil.ReadFile(output.Template); err != nil {
			log.Fatal(err)
		}

		if outputTemplate, err = outputTemplate.Parse(string(buffer)); err != nil {
			log.Fatal(err)
		}

		if outfile, err = os.Create(output.Output); err != nil {
			log.Fatal(err)
		}

		defer outfile.Close()

		if err = outputTemplate.Execute(outfile, input); err != nil {
			log.Fatal(err)
		}
	}
}

