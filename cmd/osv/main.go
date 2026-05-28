package main

import (
	"fmt"
	"os"

	osv_schema "github.com/scagogogo/osv-schema"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:   "osv",
	Short: "OSV Schema CLI - parse, validate and format OSV vulnerability data",
	Long:  "A command-line tool for working with OSV (Open Source Vulnerability) schema data. Parse, validate and inspect vulnerability JSON files.",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// parseOsvFile 是子命令共用的文件解析函数
func parseOsvFile(filePath string) (*osv_schema.OsvSchema[any, any], error) {
	return osv_schema.UnmarshalFromJsonFile[any, any](filePath)
}
