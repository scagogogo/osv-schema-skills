package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fixturePath 指向根目录的测试数据文件。
const fixturePath = "../../test_data/GHSA-vxv8-r8q2-63xw.json"

// runCapture 执行一次 rootCmd（已设置 args），把 stdout/stderr 写入 buffer 返回。
// 调用方负责在调用前 reset 各子命令的全局 flag 状态。
func runCapture(t *testing.T, args []string) (stdout string, err error) {
	t.Helper()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err = runRoot()
	return buf.String(), err
}

func TestMainFunc_Success(t *testing.T) {
	// 成功路径：version 子命令，runRoot 返回 nil，main 不调用 osExit。
	rootCmd.SetArgs([]string{"version"})
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	savedExit := osExit
	called := false
	osExit = func(code int) { called = true; _ = code }
	defer func() { osExit = savedExit }()
	main()
	assert.False(t, called, "osExit should not be invoked on success")
	assert.Contains(t, buf.String(), "osv-cli version:")
}

func TestMainFunc_ErrorInvokesOsExit(t *testing.T) {
	// 失败路径：parse 不存在的文件，runRoot 返回 error，main 应调用 osExit(1)。
	// 但 rootCmd 的 Execute 在 RunE 返回 error 时，cobra 会把 error 打印到 stderr 并返回它。
	rootCmd.SetArgs([]string{"parse", "/no/such/file.json"})
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	savedExit := osExit
	gotCode := -99
	osExit = func(code int) { gotCode = code }
	defer func() { osExit = savedExit }()
	main()
	assert.Equal(t, 1, gotCode, "osExit should be called with code 1 on error")
}

func TestRunRoot_ReturnsErrorOnBadSubcommand(t *testing.T) {
	rootCmd.SetArgs([]string{"nonexistent-subcommand"})
	err := runRoot()
	assert.NotNil(t, err)
}

func TestParseOsvFile(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		v, err := parseOsvFile(fixturePath)
		assert.Nil(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, "GHSA-vxv8-r8q2-63xw", v.ID)
	})
	t.Run("missing file", func(t *testing.T) {
		_, err := parseOsvFile("/no/such.json")
		assert.NotNil(t, err)
	})
}

// guard against unused import
var _ = strings.TrimSpace
var _ = os.Args
