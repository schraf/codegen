package internal

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

// Output defines the structure for a single output file generation task.
type Output struct {
	// Template is the path to the main template file for this output.
	Template string `json:"template"`
	// Input is the path to the JSON file containing data for the template.
	Input string `json:"input"`
	// Output is the path where the generated file will be saved.
	Output string `json:"output"`
}

// Project defines the structure of the codegen project file (e.g., project.json).
// It contains a list of include files that can be shared across templates,
// and a list of output tasks to be executed.
type Project struct {
	// Includes is a list of file paths to templates that can be included
	// in other templates. These are parsed first and can be referenced
	// by name in the output templates.
	Includes []string `json:"includes"`
	// Outputs is a list of output generation tasks. Each task specifies
	// a template, an input data file, and an output file path.
	Outputs []Output `json:"outputs"`
}

// LoadProject reads and parses a project definition file (in JSON format)
// from the given filename. It returns a Project struct instance.
func LoadProject(filename string) (*Project, error) {
	project := Project{
		Includes: []string{},
		Outputs:  []Output{},
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read project file '%s': %w", filename, err)
	}

	if err := json.Unmarshal(contents, &project); err != nil {
		return nil, fmt.Errorf("failed to parse project file '%s': %w", filename, err)
	}

	return &project, nil
}

// Execute runs the code generation process defined by the Project.
// It parses the include files, then for each output task, it loads the
// specific template and its input data, executes the template, and writes
// the result to the specified output file.
func (p Project) Execute() error {
	baseTemplate, err := template.ParseFiles(p.Includes...)
	if err != nil {
		return fmt.Errorf("failed parsing include files: %w", err)
	}

	for _, output := range p.Outputs {
		//--========================================================--
		//--== LOAD THE INPUT VARIABLES
		//--========================================================--

		inputContents, err := ioutil.ReadFile(output.Input)
		if err != nil {
			return fmt.Errorf("failed to read input file '%s': %w", output.Input, err)
		}

		var input map[string]any

		if err := json.Unmarshal(inputContents, &input); err != nil {
			return fmt.Errorf("failed to parse input file '%s': %w", output.Input, err)
		}

		//--========================================================--
		//--== LOAD THE TEMPLATE
		//--========================================================--

		templateContents, err := ioutil.ReadFile(output.Template)
		if err != nil {
			return fmt.Errorf("failed to read output template '%s': %w", output.Template, err)
		}

		outputTemplate, err := baseTemplate.Clone()
		if err != nil {
			return fmt.Errorf("failed to create template '%s': %w", output.Template, err)
		}

		outputTemplate, err = outputTemplate.Parse(string(templateContents))
		if err != nil {
			return fmt.Errorf("failed to parse template '%s': %w", output.Template, err)
		}

		//--========================================================--
		//--== EXECUTE THE TEMPLATE
		//--========================================================--

		outfile, err := os.Create(output.Output)
		if err != nil {
			return fmt.Errorf("failed to create output file '%s': %w", output.Output, err)
		}
		defer outfile.Close()

		if err := outputTemplate.Execute(outfile, input); err != nil {
			return fmt.Errorf("failed to execute template '%s': %w", output.Output, err)
		}
	}

	return nil
}
