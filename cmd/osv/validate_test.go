package main

import (
	"os"
	"testing"
)

func TestValidateCommandValidFile(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	rootCmd.SetArgs([]string{"validate", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("validate command failed for valid file: %v", err)
	}
}

func TestValidateCommandInvalidFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "invalid-osv-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("{invalid json")
	tmpFile.Close()

	result := validateFile(tmpFile.Name())
	if result.Valid {
		t.Error("expected invalid file to fail validation")
	}
}

func TestValidateFileMissingID(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "missing-id-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(`{"schema_version":"1.4.0"}`)
	tmpFile.Close()

	result := validateFile(tmpFile.Name())
	if result.Valid {
		t.Error("expected file missing id to fail validation")
	}
	if len(result.Errors) == 0 {
		t.Error("expected at least one error for missing id")
	}
}

func TestValidateFileNonExistent(t *testing.T) {
	result := validateFile("nonexistent-file.json")
	if result.Valid {
		t.Error("expected nonexistent file to fail validation")
	}
}
