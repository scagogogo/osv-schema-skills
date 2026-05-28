package main

import (
	"bytes"
	"os"
	"testing"
)

func TestParseCommand(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"parse", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("parse command failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("GHSA-vxv8-r8q2-63xw")) {
		t.Errorf("expected output to contain vulnerability ID, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("CVE-2022-35981")) {
		t.Errorf("expected output to contain CVE alias, got: %s", output)
	}
}

func TestParseCommandJSON(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"parse", "-o", "json", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("parse command with json output failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte(`"id"`)) {
		t.Errorf("expected JSON output to contain id field, got: %s", output)
	}
}

func TestParseCommandFileNotFound(t *testing.T) {
	rootCmd.SetArgs([]string{"parse", "nonexistent.json"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}
