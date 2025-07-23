package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnumGeneration(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	// Write test enum file
	enumFile := `package test

type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
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
	*typeNames = "Status"
	*output = "status_enum.go"

	// Load packages and run enumgen logic using the refactored generateAll function
	pkgs := loadPackages([]string{"."})
	if len(pkgs) == 0 {
		t.Fatal("no packages loaded")
	}

	// Use the refactored generateAll function
	generateAll(pkgs, []string{"Status"}, tmpDir)

	// Verify the generated file exists and contains expected content
	generatedContent, err := os.ReadFile("status_enum.go")
	if err != nil {
		t.Fatal(err)
	}

	generatedStr := string(generatedContent)

	// Debug: print the generated content
	t.Logf("Generated content:\n%s", generatedStr)

	// Check for expected content in the generated file
	expectedContents := []string{
		"package test",
		"_StatusValues",
		"GetStatus",
		"StatusPending",
		"StatusActive",
		"StatusInactive",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(generatedStr, expected) {
			t.Errorf("Generated file missing expected content: %s", expected)
		}
	}

	// Verify the generated code structure matches the refactored output
	// The new format uses the base type name (int) as the key type
	if !strings.Contains(generatedStr, "var _StatusValues = map[int]Status{") {
		t.Error("Generated file should contain Status values map with int keys")
	}

	if !strings.Contains(generatedStr, "func GetStatus(x int) (Status, bool) {") {
		t.Error("Generated file should contain GetStatus function with int parameter")
	}

	// Verify the function contains the correct return statement
	if !strings.Contains(generatedStr, "return _StatusValues[x]") {
		t.Error("Generated file should contain correct return statement")
	}

	// Verify the map contains the expected mappings (values as keys, names as values)
	if !strings.Contains(generatedStr, "0: StatusPending") {
		t.Error("Generated file should map 0 to StatusPending")
	}
	if !strings.Contains(generatedStr, "1: StatusActive") {
		t.Error("Generated file should map 1 to StatusActive")
	}
	if !strings.Contains(generatedStr, "2: StatusInactive") {
		t.Error("Generated file should map 2 to StatusInactive")
	}
}


