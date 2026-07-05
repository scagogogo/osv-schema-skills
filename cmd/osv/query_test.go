package main

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	osv_schema "github.com/scagogogo/osv-schema-skills"
	"github.com/stretchr/testify/assert"
)

// resetQueryFlags 清空 query 子命令的全局 flag 状态。
func resetQueryFlags() {
	querySeverity = ""
	queryMaven = false
	queryRanges = false
	queryEvents = false
	outputFormat = "text"
}

func TestRunQuery_NoFlagsReturnsError(t *testing.T) {
	resetQueryFlags()
	rootCmd.SetArgs([]string{"query", fixturePath})
	err := runRoot()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "at least one query flag")
}

func TestRunQuery_BadFileReturnsError(t *testing.T) {
	resetQueryFlags()
	rootCmd.SetArgs([]string{"query", "--severity", "cvss3", "/no/such.json"})
	err := runRoot()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to parse OSV file")
}

func TestRunQuery_SeverityCvss3Text(t *testing.T) {
	resetQueryFlags()
	out, err := runCapture(t, []string{"query", "--severity", "cvss3", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Severity (cvss3):")
	assert.Contains(t, out, "CVSS_V3")
	assert.Contains(t, out, "CVSS:3.1/")
}

func TestRunQuery_SeverityCvss2Text_None(t *testing.T) {
	resetQueryFlags()
	out, err := runCapture(t, []string{"query", "--severity", "cvss2", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Severity (cvss2):")
	assert.Contains(t, out, "(none)")
}

func TestRunQuery_SeverityInvalid(t *testing.T) {
	resetQueryFlags()
	rootCmd.SetArgs([]string{"query", "--severity", "cvss4", fixturePath})
	err := runRoot()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid --severity value")
}

func TestRunQuery_SeverityJSON(t *testing.T) {
	resetQueryFlags()
	out, err := runCapture(t, []string{"query", "--severity", "cvss3", "-o", "json", fixturePath})
	assert.Nil(t, err)
	var got map[string]any
	assert.Nil(t, json.Unmarshal([]byte(out), &got))
	assert.Equal(t, "GHSA-vxv8-r8q2-63xw", got["id"])
	sev, ok := got["severity"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "CVSS_V3", sev["type"])
}

func TestRunQuery_SeverityJSON_Cvss2Nil(t *testing.T) {
	resetQueryFlags()
	out, err := runCapture(t, []string{"query", "--severity", "cvss2", "-o", "json", fixturePath})
	assert.Nil(t, err)
	var got map[string]any
	assert.Nil(t, json.Unmarshal([]byte(out), &got))
	assert.Nil(t, got["severity"])
}

func TestRunQuery_RangesText(t *testing.T) {
	resetQueryFlags()
	out, err := runCapture(t, []string{"query", "--ranges", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Version ranges:")
	assert.Contains(t, out, "[ECOSYSTEM]")
	assert.Contains(t, out, "introduced:")
	assert.Contains(t, out, "fixed:")
}

func TestRunQuery_RangesJSON(t *testing.T) {
	resetQueryFlags()
	out, err := runCapture(t, []string{"query", "--ranges", "-o", "json", fixturePath})
	assert.Nil(t, err)
	var got map[string]any
	assert.Nil(t, json.Unmarshal([]byte(out), &got))
	ranges, ok := got["ranges"].([]any)
	assert.True(t, ok)
	assert.True(t, len(ranges) >= 1)
}

func TestRunQuery_EventsText(t *testing.T) {
	resetQueryFlags()
	out, err := runCapture(t, []string{"query", "--events", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Event timeline:")
	assert.Contains(t, out, "introduced=")
	assert.Contains(t, out, "fixed=")
}

func TestRunQuery_EventsJSON(t *testing.T) {
	resetQueryFlags()
	out, err := runCapture(t, []string{"query", "--events", "-o", "json", fixturePath})
	assert.Nil(t, err)
	var got map[string]any
	assert.Nil(t, json.Unmarshal([]byte(out), &got))
	events, ok := got["events"].([]any)
	assert.True(t, ok)
	assert.True(t, len(events) >= 1)
}

// --- DTO / collector 单元测试（用合成数据覆盖 Maven 等分支）---

func TestToSeverityDTO_Nil(t *testing.T) {
	assert.Nil(t, toSeverityDTO(nil))
}

func TestToSeverityDTO_NonNil(t *testing.T) {
	s := &osv_schema.Severity{Type: osv_schema.SeverityTypeCVSS3, Score: "7.5"}
	got := toSeverityDTO(s)
	assert.NotNil(t, got)
	assert.Equal(t, "CVSS_V3", got.Type)
	assert.Equal(t, "7.5", got.Score)
	assert.InDelta(t, 7.5, got.NumericScore, 1e-9)
}

func TestCollectMaven_WithMavenPackage(t *testing.T) {
	o := &osv_schema.OsvSchema[any, any]{
		Affected: osv_schema.AffectedSlice[any, any]{
			{Package: &osv_schema.Package{Ecosystem: osv_schema.EcosystemMaven, Name: "org.apache:lib"}},
			{Package: &osv_schema.Package{Ecosystem: osv_schema.EcosystemPyPI, Name: "requests"}},
			{Package: nil},
		},
	}
	got := collectMaven(o)
	assert.Equal(t, 1, len(got))
	assert.Equal(t, "org.apache", got[0].GroupID)
	assert.Equal(t, "lib", got[0].ArtifactID)
}

func TestCollectMaven_NoMaven(t *testing.T) {
	o := &osv_schema.OsvSchema[any, any]{
		Affected: osv_schema.AffectedSlice[any, any]{
			{Package: &osv_schema.Package{Ecosystem: osv_schema.EcosystemPyPI, Name: "requests"}},
		},
	}
	got := collectMaven(o)
	assert.Equal(t, 0, len(got))
}

func TestCollectRangesDTO(t *testing.T) {
	o := &osv_schema.OsvSchema[any, any]{
		Affected: osv_schema.AffectedSlice[any, any]{
			{
				Package: &osv_schema.Package{Ecosystem: osv_schema.EcosystemPyPI, Name: "req"},
				Ranges: []*osv_schema.Range[any]{
					{Type: osv_schema.RangeTypeEcosystem, Repo: "r", Events: osv_schema.Events{
						&osv_schema.Event{Introduced: "0", Fixed: "1.0"},
					}},
				},
			},
			{Package: nil}, // skipped
		},
	}
	got := collectRangesDTO(o)
	assert.Equal(t, 1, len(got))
	assert.Equal(t, "PyPI/req", got[0].Package)
	assert.Equal(t, "r", got[0].Repo)
	assert.Equal(t, 1, len(got[0].Events))
}

func TestCollectEventsDTO(t *testing.T) {
	o := &osv_schema.OsvSchema[any, any]{
		Affected: osv_schema.AffectedSlice[any, any]{
			{
				Package: &osv_schema.Package{Ecosystem: osv_schema.EcosystemNpm, Name: "lodash"},
				Ranges: []*osv_schema.Range[any]{
					{Events: osv_schema.Events{
						&osv_schema.Event{Introduced: "1.0"},
						&osv_schema.Event{Fixed: "2.0"},
					}},
				},
			},
			{Package: nil}, // skipped
		},
	}
	got := collectEventsDTO(o)
	assert.Equal(t, 2, len(got))
	assert.Equal(t, "npm/lodash", got[0].Package)
}

// --- printEventDTO / printEventDTOInline 单元测试 ---

func TestPrintEventDTO_AllFields(t *testing.T) {
	w := new(bytes.Buffer)
	e := eventDTO{Introduced: "1.0", Fixed: "1.1", LastAffected: "2.0", Limit: "v1"}
	printEventDTO(w, "  ", e)
	out := w.String()
	assert.Contains(t, out, "introduced: 1.0")
	assert.Contains(t, out, "fixed: 1.1")
	assert.Contains(t, out, "last_affected: 2.0")
	assert.Contains(t, out, "limit: v1")
}

func TestPrintEventDTO_Empty(t *testing.T) {
	w := new(bytes.Buffer)
	printEventDTO(w, "  ", eventDTO{})
	assert.Equal(t, "", w.String())
}

func TestPrintEventDTOInline_AllFields(t *testing.T) {
	w := new(bytes.Buffer)
	e := eventDTO{Introduced: "1.0", Fixed: "1.1", LastAffected: "2.0", Limit: "v1"}
	printEventDTOInline(w, e)
	out := w.String()
	assert.Contains(t, out, "introduced=1.0")
	assert.Contains(t, out, "fixed=1.1")
	assert.Contains(t, out, "last_affected=2.0")
	assert.Contains(t, out, "limit=v1")
}

func TestPrintEventDTOInline_Empty(t *testing.T) {
	w := new(bytes.Buffer)
	printEventDTOInline(w, eventDTO{})
	assert.Equal(t, "", w.String())
}

// --- printQueryText / printQueryJSON 直接调用覆盖 Maven/Ranges/Events 文本分支 ---

func TestPrintQueryText_MavenNone(t *testing.T) {
	resetQueryFlags()
	queryMaven = true
	o := &osv_schema.OsvSchema[any, any]{ID: "X"}
	w := new(bytes.Buffer)
	assert.Nil(t, printQueryText(w, o))
	assert.Contains(t, w.String(), "(no Maven packages)")
}

func TestPrintQueryText_MavenWithDecomposition(t *testing.T) {
	resetQueryFlags()
	queryMaven = true
	o := &osv_schema.OsvSchema[any, any]{
		ID: "X",
		Affected: osv_schema.AffectedSlice[any, any]{
			{Package: &osv_schema.Package{Ecosystem: osv_schema.EcosystemMaven, Name: "g:a"}},
		},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printQueryText(w, o))
	out := w.String()
	assert.Contains(t, out, "g:a → groupId=g, artifactId=a")
}

func TestPrintQueryText_RangesWithRepo(t *testing.T) {
	resetQueryFlags()
	queryRanges = true
	o := &osv_schema.OsvSchema[any, any]{
		ID: "X",
		Affected: osv_schema.AffectedSlice[any, any]{
			{
				Package: &osv_schema.Package{Ecosystem: osv_schema.EcosystemPyPI, Name: "p"},
				Ranges: []*osv_schema.Range[any]{
					{Type: osv_schema.RangeTypeEcosystem, Repo: "myrepo", Events: osv_schema.Events{
						&osv_schema.Event{Introduced: "0", LastAffected: "1.0", Limit: "v1"},
					}},
				},
			},
		},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printQueryText(w, o))
	out := w.String()
	assert.Contains(t, out, "(repo: myrepo)")
	assert.Contains(t, out, "introduced: 0")
	assert.Contains(t, out, "last_affected: 1.0")
	assert.Contains(t, out, "limit: v1")
}

func TestPrintQueryJSON_InvalidSeverity(t *testing.T) {
	resetQueryFlags()
	querySeverity = "cvss5"
	o := &osv_schema.OsvSchema[any, any]{ID: "X"}
	w := new(bytes.Buffer)
	err := printQueryJSON(w, o)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid --severity value")
}

func TestPrintQueryJSON_Maven(t *testing.T) {
	resetQueryFlags()
	queryMaven = true
	o := &osv_schema.OsvSchema[any, any]{
		ID: "X",
		Affected: osv_schema.AffectedSlice[any, any]{
			{Package: &osv_schema.Package{Ecosystem: osv_schema.EcosystemMaven, Name: "g:a"}},
		},
	}
	w := new(bytes.Buffer)
	assert.Nil(t, printQueryJSON(w, o))
	var got map[string]any
	assert.Nil(t, json.Unmarshal(w.Bytes(), &got))
	mavens, ok := got["maven"].([]any)
	assert.True(t, ok)
	assert.Equal(t, 1, len(mavens))
}

// guard against unused import
var _ io.Writer = (*bytes.Buffer)(nil)
