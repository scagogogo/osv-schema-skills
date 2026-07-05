package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	osv_schema "github.com/scagogogo/osv-schema-skills"
	"github.com/stretchr/testify/assert"
)

// resetParseFlags 清空 parse 子命令的全局 flag 状态。
func resetParseFlags() {
	parseVerbose = false
	outputFormat = "text"
}

func TestParseCommand(t *testing.T) {
	resetParseFlags()
	out, err := runCapture(t, []string{"parse", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "GHSA-vxv8-r8q2-63xw")
	assert.Contains(t, out, "CVE-2022-35981")
	// 非 verbose 模式不输出 Published/Details
	assert.NotContains(t, out, "Published:")
	assert.NotContains(t, out, "Details:")
}

func TestParseCommandVerbose(t *testing.T) {
	resetParseFlags()
	out, err := runCapture(t, []string{"parse", "-v", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Published:")
	assert.Contains(t, out, "Modified:")
	assert.Contains(t, out, "Details:")
	assert.Contains(t, out, "range [ECOSYSTEM]")
	assert.Contains(t, out, "introduced=")
	assert.Contains(t, out, "fixed=")
}

func TestParseCommandJSON(t *testing.T) {
	resetParseFlags()
	out, err := runCapture(t, []string{"parse", "-o", "json", fixturePath})
	assert.Nil(t, err)
	var got map[string]any
	assert.Nil(t, json.Unmarshal([]byte(out), &got))
	assert.Equal(t, "GHSA-vxv8-r8q2-63xw", got["id"])
}

func TestParseCommandFileNotFound(t *testing.T) {
	resetParseFlags()
	rootCmd.SetArgs([]string{"parse", "nonexistent.json"})
	err := runRoot()
	assert.NotNil(t, err)
}

// --- printParseText 合成数据测试，覆盖所有 verbose 分支 ---

func TestPrintParseText_Minimal(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{ID: "X", SchemaVersion: "1.4.0", Summary: "s"}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, false))
	out := w.String()
	assert.Contains(t, out, "ID:             X")
	assert.Contains(t, out, "Schema Version: 1.4.0")
	assert.Contains(t, out, "Summary:        s")
	// 没有Aliases/Severity/Affected/References 区块
	assert.NotContains(t, out, "Aliases:")
	assert.NotContains(t, out, "Severity:")
	assert.NotContains(t, out, "Affected")
	assert.NotContains(t, out, "References:")
}

func TestPrintParseText_VerboseDates(t *testing.T) {
	resetParseFlags()
	pub := time.Date(2022, 9, 16, 22, 26, 57, 0, time.UTC)
	mod := time.Date(2022, 9, 19, 19, 33, 4, 0, time.UTC)
	o := &osv_schema.OsvSchema[any, any]{
		ID:        "X",
		Published: pub,
		Modified:  mod,
		Withdrawn: "2023-01-01T00:00:00Z",
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, true))
	out := w.String()
	assert.Contains(t, out, "Published:")
	assert.Contains(t, out, "Modified:")
	assert.Contains(t, out, "Withdrawn:")
}

func TestPrintParseText_VerboseZeroDatesOmitted(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{ID: "X"}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, true))
	out := w.String()
	assert.NotContains(t, out, "Published:")
	assert.NotContains(t, out, "Modified:")
	assert.NotContains(t, out, "Withdrawn:")
}

func TestPrintParseText_AliasesAndCVE(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID:      "X",
		Aliases: osv_schema.Aliases{"GHSA-1", "CVE-2024-1"},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, false))
	out := w.String()
	assert.Contains(t, out, "Aliases:        GHSA-1, CVE-2024-1")
	assert.Contains(t, out, "CVE:            CVE-2024-1")
}

func TestPrintParseText_AliasesNoCVE(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID:      "X",
		Aliases: osv_schema.Aliases{"GHSA-1"},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, false))
	out := w.String()
	assert.Contains(t, out, "Aliases:")
	assert.NotContains(t, out, "CVE:")
}

func TestPrintParseText_VerboseRelated(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID:      "X",
		Related: osv_schema.Related{"R1", "R2"},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, true))
	assert.Contains(t, w.String(), "Related:        R1, R2")
}

func TestPrintParseText_VerboseDetails(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{ID: "X", Details: "line1\nline2"}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, true))
	out := w.String()
	assert.Contains(t, out, "Details:")
	assert.Contains(t, out, "line1")
	assert.Contains(t, out, "line2")
}

func TestPrintParseText_VerboseCreditsWithNameTypeContact(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID:      "X",
		Credits: &osv_schema.Credits{Name: "Alice", Type: "FINDER", Contact: []string{"a@x"}},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, true))
	out := w.String()
	assert.Contains(t, out, "Credits:")
	assert.Contains(t, out, "Alice")
	assert.Contains(t, out, "(FINDER)")
	assert.Contains(t, out, "a@x")
}

func TestPrintParseText_VerboseCreditsNameOnly(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID:      "X",
		Credits: &osv_schema.Credits{Name: "Bob"},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, true))
	out := w.String()
	assert.Contains(t, out, "Bob")
	assert.NotContains(t, out, "(")
}

func TestPrintParseText_Severity(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID:       "X",
		Severity: osv_schema.SeveritySlice{{Type: osv_schema.SeverityTypeCVSS3, Score: "7.5"}},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, false))
	assert.Contains(t, w.String(), "Severity:")
	assert.Contains(t, w.String(), "CVSS_V3: 7.5")
}

func TestPrintParseText_AffectedWithVersions(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID: "X",
		Affected: osv_schema.AffectedSlice[any, any]{
			{
				Package:  &osv_schema.Package{Ecosystem: osv_schema.EcosystemPyPI, Name: "p"},
				Versions: []string{"1.0", "1.1"},
			},
		},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, false))
	out := w.String()
	assert.Contains(t, out, "PyPI/p")
	assert.Contains(t, out, "versions: 1.0, 1.1")
}

func TestPrintParseText_AffectedNilPackageSkipped(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID: "X",
		Affected: osv_schema.AffectedSlice[any, any]{
			{Package: nil},
		},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, false))
	out := w.String()
	// "Affected Packages:" 头仍会打印（len(Affected)>0），但 nil package 不输出包行
	assert.Contains(t, out, "Affected Packages:")
	// nil package 被跳过，不输出 "Ecosystem/Name" 形式的包行
	assert.NotContains(t, out, "PyPI/")
	assert.NotContains(t, out, "npm/")
}

func TestPrintParseText_References(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID: "X",
		References: osv_schema.References{
			{Type: osv_schema.ReferenceTypeFix, URL: "https://f"},
		},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, false))
	out := w.String()
	assert.Contains(t, out, "References:")
	assert.Contains(t, out, "[FIX] https://f")
}

// 覆盖 verbose 模式下 range 的 repo/last_affected/limit 分支
func TestPrintParseText_VerboseRangeAllEventTypes(t *testing.T) {
	resetParseFlags()
	o := &osv_schema.OsvSchema[any, any]{
		ID: "X",
		Affected: osv_schema.AffectedSlice[any, any]{
			{
				Package: &osv_schema.Package{Ecosystem: osv_schema.EcosystemPyPI, Name: "p"},
				Ranges: []*osv_schema.Range[any]{
					{
						Type: osv_schema.RangeTypeEcosystem,
						Repo: "https://repo",
						Events: osv_schema.Events{
							&osv_schema.Event{Introduced: "1.0", Fixed: "1.1", LastAffected: "2.0", Limit: "v1"},
						},
					},
				},
			},
		},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printParseText(w, o, true))
	out := w.String()
	assert.Contains(t, out, "range [ECOSYSTEM]")
	assert.Contains(t, out, "repo=https://repo")
	assert.Contains(t, out, "introduced=1.0")
	assert.Contains(t, out, "fixed=1.1")
	assert.Contains(t, out, "last_affected=2.0")
	assert.Contains(t, out, "limit=v1")
}

// guard against unused import
var _ = os.Stat
var _ = strings.TrimSpace
