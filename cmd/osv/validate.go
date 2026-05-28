package main

import (
	"encoding/json"
	"fmt"
	"os"

	osv_schema "github.com/scagogogo/osv-schema"
	"github.com/spf13/cobra"
)

type ValidationResult struct {
	File    string   `json:"file"`
	Valid   bool     `json:"valid"`
	Errors  []string `json:"errors,omitempty"`
	ID      string   `json:"id,omitempty"`
	Version string   `json:"schema_version,omitempty"`
}

var validateCmd = &cobra.Command{
	Use:   "validate <file> [file...]",
	Short: "Validate OSV JSON files against the schema",
	Long:  "Validate one or more OSV JSON files to check if they can be parsed and contain required fields (id, schema_version).",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	results := make([]*ValidationResult, 0, len(args))
	hasError := false

	for _, filePath := range args {
		result := validateFile(filePath)
		results = append(results, result)
		if !result.Valid {
			hasError = true
		}
	}

	out := cmd.OutOrStdout()

	if outputFormat == "json" {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(results)
	} else {
		for _, r := range results {
			if r.Valid {
				fmt.Fprintf(out, "✓ %s (id=%s, schema_version=%s)\n", r.File, r.ID, r.Version)
			} else {
				fmt.Fprintf(out, "✗ %s\n", r.File)
				for _, e := range r.Errors {
					fmt.Fprintf(out, "  - %s\n", e)
				}
			}
		}
	}

	if hasError {
		os.Exit(1)
	}
	return nil
}

func validateFile(filePath string) *ValidationResult {
	result := &ValidationResult{File: filePath}

	raw, err := os.ReadFile(filePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("cannot read file: %v", err))
		return result
	}

	if !json.Valid(raw) {
		result.Errors = append(result.Errors, "file is not valid JSON")
		return result
	}

	osvData, err := osv_schema.UnmarshalFromJson[any, any](raw)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("OSV parse error: %v", err))
		return result
	}

	if osvData.ID == "" {
		result.Errors = append(result.Errors, "missing required field: id")
	}
	if osvData.SchemaVersion == "" {
		result.Errors = append(result.Errors, "missing required field: schema_version")
	}

	if len(result.Errors) == 0 {
		result.Valid = true
		result.ID = osvData.ID
		result.Version = osvData.SchemaVersion
	}
	return result
}
