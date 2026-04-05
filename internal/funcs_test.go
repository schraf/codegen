package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitWords(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"foo", []string{"foo"}},
		{"foo_bar_baz", []string{"foo", "bar", "baz"}},
		{"foo-bar-baz", []string{"foo", "bar", "baz"}},
		{"foo.bar.baz", []string{"foo", "bar", "baz"}},
		{"foo/bar/baz", []string{"foo", "bar", "baz"}},
		{"fooBarBaz", []string{"foo", "Bar", "Baz"}},
		{"FooBarBaz", []string{"Foo", "Bar", "Baz"}},
		{"XMLParser", []string{"XML", "Parser"}},
		{"someID", []string{"some", "ID"}},
		{" mixed_Case_with-UPPER", []string{"mixed", "Case", "with", "UPPER"}},
		{"user2Name", []string{"user2", "Name"}},
		{"user25Name", []string{"user25", "Name"}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, splitWords(tt.input))
		})
	}
}

func TestCaseTransformations(t *testing.T) {
	tests := []struct {
		input          string
		camel          string
		pascal         string
		snake          string
		kebab          string
		screamingSnake string
		screamingKebab string
		dot            string
		path           string
	}{
		{
			input:          "foo_bar_baz",
			camel:          "fooBarBaz",
			pascal:         "FooBarBaz",
			snake:          "foo_bar_baz",
			kebab:          "foo-bar-baz",
			screamingSnake: "FOO_BAR_BAZ",
			screamingKebab: "FOO-BAR-BAZ",
			dot:            "foo.bar.baz",
			path:           "foo/bar/baz",
		},
		{
			input:          "XMLParser",
			camel:          "xmlParser",
			pascal:         "XmlParser",
			snake:          "xml_parser",
			kebab:          "xml-parser",
			screamingSnake: "XML_PARSER",
			screamingKebab: "XML-PARSER",
			dot:            "xml.parser",
			path:           "xml/parser",
		},
		{
			input:          "someID",
			camel:          "someId",
			pascal:         "SomeId",
			snake:          "some_id",
			kebab:          "some-id",
			screamingSnake: "SOME_ID",
			screamingKebab: "SOME-ID",
			dot:            "some.id",
			path:           "some/id",
		},
		{
			input:          "",
			camel:          "",
			pascal:         "",
			snake:          "",
			kebab:          "",
			screamingSnake: "",
			screamingKebab: "",
			dot:            "",
			path:           "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.camel, toCamelCase(tt.input), "camelCase")
			assert.Equal(t, tt.pascal, toPascalCase(tt.input), "PascalCase")
			assert.Equal(t, tt.snake, toSnakeCase(tt.input), "snake_case")
			assert.Equal(t, tt.kebab, toKebabCase(tt.input), "kebab-case")
			assert.Equal(t, tt.screamingSnake, toScreamingSnake(tt.input), "SCREAMING_SNAKE")
			assert.Equal(t, tt.screamingKebab, toScreamingKebab(tt.input), "SCREAMING-KEBAB")
			assert.Equal(t, tt.dot, toDotCase(tt.input), "dot.case")
			assert.Equal(t, tt.path, toPathCase(tt.input), "path/case")
		})
	}
}

func TestCapitalize(t *testing.T) {
	assert.Equal(t, "Foo", capitalize("foo"))
	assert.Equal(t, "Foo", capitalize("FOO"))
	assert.Equal(t, "F", capitalize("f"))
	assert.Equal(t, "", capitalize(""))
}
