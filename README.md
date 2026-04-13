# Codegen

A simple, template-based code generation tool written in Go. It uses Go's built-in `html/template` package to generate any kind of text-based file from JSON data sources.

## How it Works

The tool is driven by a `project.json` file (or any other name passed as an argument) that defines the code generation tasks. This file specifies:

-   **Includes:** A list of reusable, shared template files.
-   **Outputs:** A list of files to generate, each with its own main template, input data file, and output path.

## Usage

1.  Create a `project.json` file to define your code generation tasks.
2.  Run the tool, passing the project file as an argument.

```bash
go run ./cmd/main.go project.json
```

## The `project.json` Format

The project file has two main sections: `includes` and `outputs`.

-   `includes`: An array of strings, where each string is a path to a template file that can be included and reused in other templates. These are parsed first.
-   `outputs`: An array of objects, where each object defines a file to be generated.

### Output Object Structure

Each object in the `outputs` array has the following properties:

-   `template`: The path to the main template file for this output.
-   `input`: The path to the JSON file containing the data to be used in the template.
-   `output`: The path where the generated file will be saved.

### Example `project.json`

```json
{
  "includes": [
    "templates/helpers.tmpl"
  ],
  "outputs": [
    {
      "template": "templates/service.go.tmpl",
      "input": "data/user_service.json",
      "output": "gen/user_service.go"
    }
  ]
}
```

## Template Functions

The engine includes a wide variety of built-in functions that you can use in your templates.

### String Case Conversions

These are especially useful for converting variable names between different conventions:

| Function | Example Usage | Result for "user_id" |
| :--- | :--- | :--- |
| `toCamelCase` | `{{.Name \| toCamelCase}}` | `userId` |
| `toPascalCase` | `{{.Name \| toPascalCase}}` | `UserId` |
| `toSnakeCase` | `{{.Name \| toSnakeCase}}` | `user_id` |
| `toKebabCase` | `{{.Name \| toKebabCase}}` | `user-id` |
| `toScreamingSnake` | `{{.Name \| toScreamingSnake}}` | `USER_ID` |
| `toScreamingKebab` | `{{.Name \| toScreamingKebab}}` | `USER-ID` |
| `toDotCase` | `{{.Name \| toDotCase}}` | `user.id` |
| `toPathCase` | `{{.Name \| toPathCase}}` | `user/id` |
| `toLower` | `{{.Name \| toLower}}` | `user_id` |
| `toUpper` | `{{.Name \| toUpper}}` | `USER_ID` |
| `capitalize` | `{{.Name \| capitalize}}` | `User_id` |

*Note: The case conversion functions are smart and can handle inputs in snake_case, camelCase, PascalCase, kebab-case, or mixed cases seamlessly.*

### String Manipulation

| Function | Example Usage | Description |
| :--- | :--- | :--- |
| `trimSpace` | `{{.Name \| trimSpace}}` | Removes leading and trailing white space. |
| `trimPrefix` | `{{.Name \| trimPrefix "user_"}}` | Removes the specified prefix. |
| `trimSuffix` | `{{.Name \| trimSuffix "_id"}}` | Removes the specified suffix. |
| `hasPrefix` | `{{hasPrefix "user_" .Name}}` | Checks if a string starts with a prefix. |
| `hasSuffix` | `{{hasSuffix "_id" .Name}}` | Checks if a string ends with a suffix. |
| `contains` | `{{contains "id" .Name}}` | Checks if a string contains a substring. |
| `replace` | `{{replace .Name "old" "new"}}` | Replaces all occurrences of a substring. |
| `split` | `{{split .Name ","}}` | Splits a string into an array of strings. |
| `join` | `{{join .List ","}}` | Joins an array of strings into a single string. |
| `indent` | `{{.Content \| indent 4}}` | Indents each non-empty line with 4 spaces. |
| `indentTab` | `{{.Content \| indentTab 2}}` | Indents each non-empty line with 2 tabs. |
| `regexMatch` | `{{regexMatch "^[a-z]+$" .Name}}` | Reports whether the string matches the regex. |
| `regexReplace` | `{{regexReplace "a+" "b" .Name}}` | Replaces regex matches in a string. |

### Math & Numeric

| Function | Example Usage | Result |
| :--- | :--- | :--- |
| `add` | `{{add 1 2}}` | `3` |
| `sub` | `{{sub 5 2}}` | `3` |
| `mul` | `{{mul 2 3}}` | `6` |
| `div` | `{{div 6 2}}` | `3` |
| `mod` | `{{mod 7 3}}` | `1` |
| `hexString` | `{{hexString 255}}` | `"ff"` |

### File & System Utilities

| Function | Example Usage | Description |
| :--- | :--- | :--- |
| `fileSize` | `{{fileSize "path/to/file"}}` | Returns the file size in bytes. |
| `readFile` | `{{readFile "path/to/file"}}` | Reads and returns the file contents as a string. |
| `env` | `{{env "USER"}}` | Retrieves the value of an environment variable. |

### Encoding & Hashing

| Function | Example Usage | Description |
| :--- | :--- | :--- |
| `hash` | `{{hash "my string"}}` | Returns the xxHash 64-bit digest of the string. |
| `fileHash` | `{{fileHash "path/to/file"}}` | Returns the xxHash 64-bit digest of the file contents. |
| `base64Encode` | `{{base64Encode "data"}}` | Returns the base64 encoding of the string. |
| `base64Decode` | `{{base64Decode "ZGF0YQ=="}}` | Decodes a base64 encoded string. |

### General Utilities

| Function | Example Usage | Description |
| :--- | :--- | :--- |
| `default` | `{{.Name \| default "unnamed"}}` | Returns the default value if the input is considered empty. |
| `seq` | `{{seq 1 10 2}}` | Generates a sequence of integers for iterating (e.g. `[1, 3, 5, 7, 9]`). |

## Example

Let's say you want to generate a simple Go struct.

### 1. Project File

**`project.json`**
```json
{
  "includes": [],
  "outputs": [
    {
      "template": "struct.tmpl",
      "input": "struct_data.json",
      "output": "person.go"
    }
  ]
}
```

### 2. Input Data

**`struct_data.json`**
```json
{
  "PackageName": "main",
  "StructName": "Person",
  "Fields": [
    { "Name": "Name", "Type": "string" },
    { "Name": "Age", "Type": "int" }
  ]
}
```

### 3. Template File

**`struct.tmpl`**
```go-template
// Code generated by codegen. DO NOT EDIT.
package {{.PackageName}}

type {{.StructName}} struct {
{{- range .Fields}}
    {{.Name}} {{.Type}}
{{- end}}
}
```

### 4. Run the tool

```bash
go run ./cmd/main.go project.json
```

### 5. Generated Output

**`person.go`**
```go
// Code generated by codegen. DO NOT EDIT.
package main

type Person struct {
    Name string
    Age int
}
```

## Building from Source

### Prerequisites

-   [Go](https://golang.org/doc/install) (version 1.24 or later)

### Build Command

To build the executable, run:
```bash
go build -o codegen ./cmd/main.go
```
You can then run the tool directly: `./codegen project.json`.

### Makefile

A `Makefile` is provided with common targets:

-   `make build`: Build the project.
-   `make test`: Run tests.
-   `make vet`: Vet the Go source code for issues.
-   `make fmt`: Format the Go source code.
-   `make all`: Run `deps`, `vet`, and `test`.

## Using with CMake

You can integrate `codegen` into your C/C++ projects using the provided `cmake/codegen.cmake` module.

### 1. Include the module

Copy `cmake/codegen.cmake` to your project and include it in your `CMakeLists.txt`:

```cmake
include(cmake/codegen.cmake)
```

### 2. Initialize the tool

Call `codegen_init` with a specific Git tag or branch to download and build the `codegen` tool:

```cmake
codegen_init(main)
```

### 3. Generate code for a target

Use `codegen_project` to link a `project.json` file to your target. This will automatically run the generation tool during the build process and add the generated files to your target's source list.

```cmake
add_executable(my_app main.c)

codegen_project(my_app
    PROJECT_FILE "codegen_project.json"
    OUTPUTS "generated_source.c" "generated_header.h"
    DEPENDS "templates/template.tmpl" "data/data.json"
)
```

- `PROJECT_FILE`: The path to your `codegen` project configuration.
- `OUTPUTS`: The list of files that `codegen` will generate.
- `DEPENDS`: (Optional) Additional files that should trigger a re-generation if modified (e.g., templates or data files).
