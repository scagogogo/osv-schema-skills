package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cliVersion    = "dev"
	schemaVersion = "1.4.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI and schema version information",
	Long:  "Display the CLI version and the supported OSV schema version.",
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "osv-cli version: %s\n", cliVersion)
	fmt.Fprintf(out, "OSV schema version: %s\n", schemaVersion)
	return nil
}
