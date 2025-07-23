package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDifferentIntTypeAliases(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	// Write test enum file with different int-based type aliases
	enumFile := `package test

type Priority int8

const (
    Low Priority = iota
    Medium
    High
)

type Size int64

const (
    Small Size = 1
    Large Size = 100
)

type Level uint

const (
    Beginner Level = iota
    Intermediate
    Advanced
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

	// Test Priority (int8)
	*typeNames = "Priority"
	*output = "priority_enum.go"

	pkgs := loadPackages([]string{"."})
	if len(pkgs) == 0 {
		t.Fatal("no packages loaded")
	}

	generateAll(pkgs, []string{"Priority"}, tmpDir)

	generatedContent, err := os.ReadFile("priority_enum.go")
	if err != nil {
		t.Fatal(err)
	}

	generatedStr := string(generatedContent)
	t.Logf("Priority enum generated content:\n%s", generatedStr)

	// Verify Priority enum uses int8 as base type
	if !strings.Contains(generatedStr, "var _PriorityValues = map[int8]Priority{") {
		t.Error("Generated file should contain Priority values map with int8 keys")
	}

	if !strings.Contains(generatedStr, "func GetPriority(x int8) (Priority, bool) {") {
		t.Error("Generated file should contain GetPriority function with int8 parameter")
	}

	// Test Size (int64)
	*typeNames = "Size"
	*output = "size_enum.go"

	generateAll(pkgs, []string{"Size"}, tmpDir)

	generatedContent, err = os.ReadFile("size_enum.go")
	if err != nil {
		t.Fatal(err)
	}

	generatedStr = string(generatedContent)
	t.Logf("Size enum generated content:\n%s", generatedStr)

	// Verify Size enum uses int64 as base type
	if !strings.Contains(generatedStr, "var _SizeValues = map[int64]Size{") {
		t.Error("Generated file should contain Size values map with int64 keys")
	}

	if !strings.Contains(generatedStr, "func GetSize(x int64) (Size, bool) {") {
		t.Error("Generated file should contain GetSize function with int64 parameter")
	}

	// Verify the map contains the expected mappings
	if !strings.Contains(generatedStr, "1: Small") {
		t.Error("Generated file should map 1 to Small")
	}
	if !strings.Contains(generatedStr, "100: Large") {
		t.Error("Generated file should map 100 to Large")
	}

	// Test Level (uint)
	*typeNames = "Level"
	*output = "level_enum.go"

	generateAll(pkgs, []string{"Level"}, tmpDir)

	generatedContent, err = os.ReadFile("level_enum.go")
	if err != nil {
		t.Fatal(err)
	}

	generatedStr = string(generatedContent)
	t.Logf("Level enum generated content:\n%s", generatedStr)

	// Verify Level enum uses uint as base type
	if !strings.Contains(generatedStr, "var _LevelValues = map[uint]Level{") {
		t.Error("Generated file should contain Level values map with uint keys")
	}

	if !strings.Contains(generatedStr, "func GetLevel(x uint) (Level, bool) {") {
		t.Error("Generated file should contain GetLevel function with uint parameter")
	}
}

func TestStringTypeAliases(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	// Write test enum file with different string-based type aliases
	enumFile := `package test

type HTTPMethod string

const (
    GET    HTTPMethod = "GET"
    POST   HTTPMethod = "POST"
    PUT    HTTPMethod = "PUT"
    DELETE HTTPMethod = "DELETE"
)

type MimeType string

const (
    JSON MimeType = "application/json"
    XML  MimeType = "application/xml"
    HTML MimeType = "text/html"
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

	// Test HTTPMethod
	*typeNames = "HTTPMethod"
	*output = "httpmethod_enum.go"

	pkgs := loadPackages([]string{"."})
	if len(pkgs) == 0 {
		t.Fatal("no packages loaded")
	}

	generateAll(pkgs, []string{"HTTPMethod"}, tmpDir)

	generatedContent, err := os.ReadFile("httpmethod_enum.go")
	if err != nil {
		t.Fatal(err)
	}

	generatedStr := string(generatedContent)
	t.Logf("HTTPMethod enum generated content:\n%s", generatedStr)

	// Verify HTTPMethod enum uses string as base type
	if !strings.Contains(generatedStr, "var _HTTPMethodValues = map[string]HTTPMethod{") {
		t.Error("Generated file should contain HTTPMethod values map with string keys")
	}

	if !strings.Contains(generatedStr, "func GetHTTPMethod(x string) (HTTPMethod, bool) {") {
		t.Error("Generated file should contain GetHTTPMethod function with string parameter")
	}

	// Verify the map contains the expected mappings
	if !strings.Contains(generatedStr, `"GET": GET`) {
		t.Error("Generated file should map \"GET\" to GET")
	}
	if !strings.Contains(generatedStr, `"POST": POST`) {
		t.Error("Generated file should map \"POST\" to POST")
	}

	// Test MimeType
	*typeNames = "MimeType"
	*output = "mimetype_enum.go"

	generateAll(pkgs, []string{"MimeType"}, tmpDir)

	generatedContent, err = os.ReadFile("mimetype_enum.go")
	if err != nil {
		t.Fatal(err)
	}

	generatedStr = string(generatedContent)
	t.Logf("MimeType enum generated content:\n%s", generatedStr)

	// Verify MimeType enum uses string as base type
	if !strings.Contains(generatedStr, "var _MimeTypeValues = map[string]MimeType{") {
		t.Error("Generated file should contain MimeType values map with string keys")
	}

	// Verify the map contains the expected mappings with complex string values
	if !strings.Contains(generatedStr, `"application/json": JSON`) {
		t.Error("Generated file should map \"application/json\" to JSON")
	}
	if !strings.Contains(generatedStr, `"text/html": HTML`) {
		t.Error("Generated file should map \"text/html\" to HTML")
	}
}
