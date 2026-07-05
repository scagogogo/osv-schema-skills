package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunVersion(t *testing.T) {
	rootCmd.SetArgs([]string{"version"})
	out, err := runCapture(t, []string{"version"})
	assert.Nil(t, err)
	assert.Contains(t, out, "osv-cli version:")
	assert.Contains(t, out, "OSV schema version:")
	// 默认 cliVersion 是 "dev"（测试环境未注入 ldflags）
	assert.Contains(t, out, "osv-cli version: dev")
	assert.Contains(t, out, "1.4.0")
	// version 忽略 -o json
	_ = strings.Contains
}

func TestRunVersion_IgnoresJSONFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"version", "-o", "json"})
	out, err := runCapture(t, []string{"version", "-o", "json"})
	assert.Nil(t, err)
	// 仍然是纯文本两行，不输出 JSON
	assert.Contains(t, out, "osv-cli version:")
	assert.NotContains(t, out, "{")
}
