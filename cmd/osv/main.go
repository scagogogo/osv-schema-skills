package main

import (
	"fmt"
	"os"

	osv_schema "github.com/scagogogo/osv-schema-skills"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	// osExit 可在测试中替换，避免 main() 的 error 路径调用真正的 os.Exit 终止测试进程。
	osExit = os.Exit
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
	if err := runRoot(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		osExit(1)
	}
}

// runRoot 执行根命令并返回错误，便于测试。
func runRoot() error {
	return rootCmd.Execute()
}

// parseOsvFile 是子命令共用的文件解析函数
func parseOsvFile(filePath string) (*osv_schema.OsvSchema[any, any], error) {
	return osv_schema.UnmarshalFromJsonFile[any, any](filePath)
}
