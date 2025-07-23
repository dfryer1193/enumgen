package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSingleValueEnum(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	// Write test enum file with single value enum
	enumFile := `package test

type State int

const (
    Active State = 42
)
`

	err := os.WriteFile(filepath.Join(tmpDir, "enum.go"), []byte(enumFile), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create go.mod file for the test package
	goModFile := `module test

go 1.21
`
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModFile), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Save current directory and change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Set up flags for enumgen
	*typeNames = "State"
	*output = "state_enum.go"

	// Load packages and run enumgen logic
	pkgs := loadPackages([]string{"."})
	if len(pkgs) == 0 {
		t.Fatal("no packages loaded")
	}

	generateAll(pkgs, []string{"State"}, tmpDir)

	// Verify the generated file exists and contains expected content
	generatedContent, err := os.ReadFile("state_enum.go")
	if err != nil {
		t.Fatal(err)
	}

	generatedStr := string(generatedContent)
	t.Logf("Generated content:\n%s", generatedStr)

	// Verify single value enum works correctly
	if !strings.Contains(generatedStr, "var _StateValues = map[int]State{") {
		t.Error("Generated file should contain State values map with int keys")
	}

	if !strings.Contains(generatedStr, "42: Active") {
		t.Error("Generated file should map 42 to Active")
	}

	if !strings.Contains(generatedStr, "func GetState(x int) (State, bool) {") {
		t.Error("Generated file should contain GetState function")
	}
}

func TestComplexStringValues(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	// Write test enum file with complex string values
	enumFile := `package test

type ErrorCode string

const (
    NotFound     ErrorCode = "NOT_FOUND"
    Unauthorized ErrorCode = "UNAUTHORIZED"
    ServerError  ErrorCode = "INTERNAL_SERVER_ERROR"
)
`

	err := os.WriteFile(filepath.Join(tmpDir, "enum.go"), []byte(enumFile), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create go.mod file for the test package
	goModFile := `module test

go 1.21
`
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModFile), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Save current directory and change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Set up flags for enumgen
	*typeNames = "ErrorCode"
	*output = "errorcode_enum.go"

	// Load packages and run enumgen logic
	pkgs := loadPackages([]string{"."})
	if len(pkgs) == 0 {
		t.Fatal("no packages loaded")
	}

	generateAll(pkgs, []string{"ErrorCode"}, tmpDir)

	// Verify the generated file exists and contains expected content
	generatedContent, err := os.ReadFile("errorcode_enum.go")
	if err != nil {
		t.Fatal(err)
	}

	generatedStr := string(generatedContent)
	t.Logf("Generated content:\n%s", generatedStr)

	// Verify complex string values work correctly
	if !strings.Contains(generatedStr, "var _ErrorCodeValues = map[string]ErrorCode{") {
		t.Error("Generated file should contain ErrorCode values map with string keys")
	}

	if !strings.Contains(generatedStr, `"NOT_FOUND": NotFound`) {
		t.Error("Generated file should map \"NOT_FOUND\" to NotFound")
	}

	if !strings.Contains(generatedStr, `"UNAUTHORIZED": Unauthorized`) {
		t.Error("Generated file should map \"UNAUTHORIZED\" to Unauthorized")
	}

	if !strings.Contains(generatedStr, `"INTERNAL_SERVER_ERROR": ServerError`) {
		t.Error("Generated file should map \"INTERNAL_SERVER_ERROR\" to ServerError")
	}

	if !strings.Contains(generatedStr, "func GetErrorCode(x string) (ErrorCode, bool) {") {
		t.Error("Generated file should contain GetErrorCode function")
	}
}

func TestMixedIntValues(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	// Write test enum file with mixed int values (not sequential)
	enumFile := `package test

type Code int

const (
    Success Code = 200
    NotFound Code = 404
    ServerError Code = 500
)
`

	err := os.WriteFile(filepath.Join(tmpDir, "enum.go"), []byte(enumFile), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create go.mod file for the test package
	goModFile := `module test

go 1.21
`
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModFile), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Save current directory and change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Set up flags for enumgen
	*typeNames = "Code"
	*output = "code_enum.go"

	// Load packages and run enumgen logic
	pkgs := loadPackages([]string{"."})
	if len(pkgs) == 0 {
		t.Fatal("no packages loaded")
	}

	generateAll(pkgs, []string{"Code"}, tmpDir)

	// Verify the generated file exists and contains expected content
	generatedContent, err := os.ReadFile("code_enum.go")
	if err != nil {
		t.Fatal(err)
	}

	generatedStr := string(generatedContent)
	t.Logf("Generated content:\n%s", generatedStr)

	// Verify mixed int values work correctly
	if !strings.Contains(generatedStr, "var _CodeValues = map[int]Code{") {
		t.Error("Generated file should contain Code values map with int keys")
	}

	if !strings.Contains(generatedStr, "200: Success") {
		t.Error("Generated file should map 200 to Success")
	}

	if !strings.Contains(generatedStr, "404: NotFound") {
		t.Error("Generated file should map 404 to NotFound")
	}

	if !strings.Contains(generatedStr, "500: ServerError") {
		t.Error("Generated file should map 500 to ServerError")
	}

	if !strings.Contains(generatedStr, "func GetCode(x int) (Code, bool) {") {
		t.Error("Generated file should contain GetCode function")
	}
}
