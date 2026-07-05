package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// resetValidateFlags 清空 validate 子命令的全局 flag 状态。
func resetValidateFlags() {
	outputFormat = "text"
}

func TestValidateCommandValidFile(t *testing.T) {
	resetValidateFlags()
	rootCmd.SetArgs([]string{"validate", fixturePath})
	err := runRoot()
	assert.Nil(t, err)
}

func TestValidateCommandJSON(t *testing.T) {
	resetValidateFlags()
	out, err := runCapture(t, []string{"validate", "-o", "json", fixturePath})
	assert.Nil(t, err)
	var got []map[string]any
	assert.Nil(t, json.Unmarshal([]byte(out), &got))
	assert.Equal(t, 1, len(got))
	assert.Equal(t, true, got[0]["valid"])
	assert.Equal(t, "GHSA-vxv8-r8q2-63xw", got[0]["id"])
}

func TestValidateCommandMultipleFilesMixed(t *testing.T) {
	resetValidateFlags()
	// 一个有效文件 + 一个无效文件
	tmpFile, err := os.CreateTemp("", "invalid-osv-*.json")
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("{invalid json")
	tmpFile.Close()

	// hasError=true 会调 osExit(1)，注入替换避免测试进程退出
	savedExit := osExit
	var gotCode int
	osExit = func(code int) { gotCode = code }
	defer func() { osExit = savedExit }()

	out, _ := runCapture(t, []string{"validate", fixturePath, tmpFile.Name()})
	assert.Equal(t, 1, gotCode, "should osExit(1) when any file invalid")
	// 文本输出含一个 ✓ 和一个 ✗
	assert.Contains(t, out, "✓")
	assert.Contains(t, out, "✗")
}

func TestValidateCommandInvalidFile(t *testing.T) {
	resetValidateFlags()
	tmpFile, err := os.CreateTemp("", "invalid-osv-*.json")
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("{invalid json")
	tmpFile.Close()

	savedExit := osExit
	var gotCode int
	osExit = func(code int) { gotCode = code }
	defer func() { osExit = savedExit }()

	_, _ = runCapture(t, []string{"validate", tmpFile.Name()})
	assert.Equal(t, 1, gotCode)
}

func TestValidateFileMissingID(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "missing-id-*.json")
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(`{"schema_version":"1.4.0"}`)
	tmpFile.Close()

	result := validateFile(tmpFile.Name())
	if assert.False(t, result.Valid) {
		assert.NotEmpty(t, result.Errors)
		// 至少包含 missing id 错误
		found := false
		for _, e := range result.Errors {
			if contains(e, "id") {
				found = true
				break
			}
		}
		assert.True(t, found, "expected an error about missing id")
	}
}

func TestValidateFileMissingSchemaVersion(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "missing-ver-*.json")
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(`{"id":"X"}`)
	tmpFile.Close()

	result := validateFile(tmpFile.Name())
	assert.False(t, result.Valid)
	found := false
	for _, e := range result.Errors {
		if contains(e, "schema_version") {
			found = true
			break
		}
	}
	assert.True(t, found, "expected an error about missing schema_version")
}

func TestValidateFileNonExistent(t *testing.T) {
	result := validateFile("nonexistent-file.json")
	assert.False(t, result.Valid)
	// 读文件失败错误
	assert.NotEmpty(t, result.Errors)
}

func TestValidateFileInvalidJSON(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "bad-json-*.json")
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(`{not valid json`)
	tmpFile.Close()

	result := validateFile(tmpFile.Name())
	assert.False(t, result.Valid)
	// json.Valid=false 分支
	found := false
	for _, e := range result.Errors {
		if contains(e, "valid JSON") {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestValidateFileParseError(t *testing.T) {
	// JSON 合法但结构不符合 OSV（比如字段类型错）会触发 UnmarshalFromJson 错误
	// 这里用一个合法 JSON 但 affected 是字符串而非数组，触发 decode 错误
	dir := t.TempDir()
	p := filepath.Join(dir, "decode-err.json")
	assert.Nil(t, os.WriteFile(p, []byte(`{"id":"X","schema_version":"1.4.0","affected":"not-an-array"}`), 0o644))

	result := validateFile(p)
	// affected 字段类型不对，UnmarshalFromJson 返回 error
	assert.False(t, result.Valid)
	found := false
	for _, e := range result.Errors {
		if contains(e, "OSV parse error") {
			found = true
			break
		}
	}
	assert.True(t, found, "expected OSV parse error, got: %v", result.Errors)
}

func TestValidateFileValid(t *testing.T) {
	result := validateFile(fixturePath)
	assert.True(t, result.Valid)
	assert.Equal(t, "GHSA-vxv8-r8q2-63xw", result.ID)
	assert.Equal(t, "1.4.0", result.Version)
}

// contains 是局部字符串包含 helper，避免引入 strings 包。
func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
