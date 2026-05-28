package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	osv_schema "github.com/scagogogo/osv-schema"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse <file>",
	Short: "Parse an OSV JSON file and display key fields",
	Long:  "Parse an OSV JSON file and display its vulnerability ID, summary, severity, affected packages and other key information.",
	Args:  cobra.ExactArgs(1),
	RunE:  runParse,
}

func init() {
	rootCmd.AddCommand(parseCmd)
}

func runParse(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	osvData, err := parseOsvFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse OSV file %s: %w", filePath, err)
	}

	out := cmd.OutOrStdout()

	if outputFormat == "json" {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(osvData)
	}
	return printParseText(out, osvData)
}

func printParseText(w io.Writer, o *osv_schema.OsvSchema[any, any]) error {
	fmt.Fprintf(w, "ID:             %s\n", o.ID)
	fmt.Fprintf(w, "Schema Version: %s\n", o.SchemaVersion)
	fmt.Fprintf(w, "Summary:        %s\n", o.Summary)

	if len(o.Aliases) > 0 {
		fmt.Fprintf(w, "Aliases:        %s\n", strings.Join(o.Aliases, ", "))
		if cve := o.Aliases.GetCVE(); cve != "" {
			fmt.Fprintf(w, "CVE:            %s\n", cve)
		}
	}

	if len(o.Severity) > 0 {
		fmt.Fprintln(w, "\nSeverity:")
		for _, s := range o.Severity {
			fmt.Fprintf(w, "  %s: %s (score: %.1f)\n", s.Type, s.Score, s.GetScore())
		}
	}

	if len(o.Affected) > 0 {
		fmt.Fprintln(w, "\nAffected Packages:")
		for _, a := range o.Affected {
			if a.Package != nil {
				fmt.Fprintf(w, "  %s/%s", a.Package.Ecosystem, a.Package.Name)
				if len(a.Versions) > 0 {
					fmt.Fprintf(w, " (versions: %s)", strings.Join(a.Versions, ", "))
				}
				fmt.Fprintln(w)
			}
		}
	}

	if len(o.References) > 0 {
		fmt.Fprintln(w, "\nReferences:")
		for _, r := range o.References {
			fmt.Fprintf(w, "  [%s] %s\n", r.Type, r.URL)
		}
	}

	return nil
}
