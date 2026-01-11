package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a temporary file with content.
func createTempFile(t *testing.T, dir, filename, content string) string {
	t.Helper()
	path := filepath.Join(dir, filename)
	err := os.WriteFile(path, []byte(content), 0600)
	require.NoError(t, err, "Failed to create temp file %s", path)
	return path
}

func TestLoadProject(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		tempDir := t.TempDir()
		projectFileContent := `{
			"includes": ["base.tmpl"],
			"outputs": [
				{
					"template": "page.tmpl",
					"input": "data.json",
					"output": "page.html"
				}
			]
		}`
		projectFilePath := createTempFile(t, tempDir, "project.json", projectFileContent)

		project, err := LoadProject(projectFilePath)

		require.NoError(t, err)
		require.NotNil(t, project)

		assert.Equal(t, []string{"base.tmpl"}, project.Includes)
		require.Len(t, project.Outputs, 1)
		assert.Equal(t, "page.tmpl", project.Outputs[0].Template)
		assert.Equal(t, "data.json", project.Outputs[0].Input)
		assert.Equal(t, "page.html", project.Outputs[0].Output)
	})

	t.Run("file not found", func(t *testing.T) {
		project, err := LoadProject("non-existent-file.json")
		require.Error(t, err)
		assert.Nil(t, project)
	})

	t.Run("invalid json", func(t *testing.T) {
		tempDir := t.TempDir()
		projectFilePath := createTempFile(t, tempDir, "project.json", "{ not json }")

		project, err := LoadProject(projectFilePath)
		require.Error(t, err)
		assert.Nil(t, project)
	})
}

func TestProject_Execute(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create test files
		baseTmplPath := createTempFile(t, tempDir, "base.tmpl", `{{define "base"}}Base content: {{template "content" .}}{{end}}`)
		pageTmplPath := createTempFile(t, tempDir, "page.tmpl", `{{template "base" .}}{{define "content"}}Hello, {{.Name}}!{{end}}`)
		inputJSONPath := createTempFile(t, tempDir, "data.json", `{"Name": "World"}`)
		outputPath := filepath.Join(tempDir, "page.html")

		project := &Project{
			Includes: []string{baseTmplPath},
			Outputs: []Output{
				{
					Template: pageTmplPath,
					Input:    inputJSONPath,
					Output:   outputPath,
				},
			},
		}

		err := project.Execute()
		require.NoError(t, err)

		outputContent, err := os.ReadFile(outputPath)
		require.NoError(t, err)
		assert.Equal(t, `Base content: Hello, World!`, string(outputContent))
	})

	t.Run("missing include file", func(t *testing.T) {
		project := &Project{Includes: []string{"non-existent.tmpl"}}
		err := project.Execute()
		require.Error(t, err)
	})

	t.Run("missing output template file", func(t *testing.T) {
		tempDir := t.TempDir()
		baseTmplPath := createTempFile(t, tempDir, "base.tmpl", `{{define "base"}}{{end}}`)
		dummyInputPath := createTempFile(t, tempDir, "input.json", "{}")

		project := &Project{
			Includes: []string{baseTmplPath},
			Outputs: []Output{
				{
					Template: "non-existent.tmpl",
					Input:    dummyInputPath,
					Output:   "output.html",
				},
			},
		}
		err := project.Execute()
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read output template")
	})

	t.Run("missing input file", func(t *testing.T) {
		tempDir := t.TempDir()
		baseTmplPath := createTempFile(t, tempDir, "base.tmpl", `{{define "base"}}{{end}}`)
		pageTmplPath := createTempFile(t, tempDir, "page.tmpl", `Hello`)

		project := &Project{
			Includes: []string{baseTmplPath},
			Outputs: []Output{
				{
					Template: pageTmplPath,
					Input:    "non-existent.json",
					Output:   "output.html",
				},
			},
		}
		err := project.Execute()
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read input file")
	})

	t.Run("invalid input json", func(t *testing.T) {
		tempDir := t.TempDir()
		baseTmplPath := createTempFile(t, tempDir, "base.tmpl", `{{define "base"}}{{end}}`)
		pageTmplPath := createTempFile(t, tempDir, "page.tmpl", `Hello, {{.Name}}`)
		inputJSONPath := createTempFile(t, tempDir, "data.json", `{"Name": "World"`) // Invalid JSON

		project := &Project{
			Includes: []string{baseTmplPath},
			Outputs: []Output{
				{
					Template: pageTmplPath,
					Input:    inputJSONPath,
					Output:   "output.html",
				},
			},
		}

		err := project.Execute()
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse input file")
	})

	t.Run("invalid template syntax", func(t *testing.T) {
		tempDir := t.TempDir()
		baseTmplPath := createTempFile(t, tempDir, "base.tmpl", `{{define "base"}}{{end}}`)
		pageTmplPath := createTempFile(t, tempDir, "page.tmpl", `{{ .Name }`) // Invalid template
		inputJSONPath := createTempFile(t, tempDir, "data.json", `{"Name": "World"}`)

		project := &Project{
			Includes: []string{baseTmplPath},
			Outputs: []Output{
				{
					Template: pageTmplPath,
					Input:    inputJSONPath,
					Output:   "output.html",
				},
			},
		}

		err := project.Execute()
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse template")
	})
}
