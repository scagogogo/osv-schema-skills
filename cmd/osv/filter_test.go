package main

import (
	"encoding/json"
	"testing"

	osv_schema "github.com/scagogogo/osv-schema-skills"
	"github.com/stretchr/testify/assert"
)

// resetFilterFlags 清空 filter 子命令的全局 flag 状态，避免测试间残留。
func resetFilterFlags() {
	filterEcosystem = ""
	filterRefType = ""
	filterAlias = ""
	outputFormat = "text"
}

func TestRunFilter_NoFlagsReturnsError(t *testing.T) {
	resetFilterFlags()
	rootCmd.SetArgs([]string{"filter", fixturePath})
	err := runRoot()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "at least one filter flag")
}

func TestRunFilter_BadFileReturnsError(t *testing.T) {
	resetFilterFlags()
	rootCmd.SetArgs([]string{"filter", "-e", "PyPI", "/no/such.json"})
	err := runRoot()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to parse OSV file")
}

func TestRunFilter_EcosystemText(t *testing.T) {
	resetFilterFlags()
	out, err := runCapture(t, []string{"filter", "-e", "PyPI", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Ecosystem filter: PyPI")
	assert.Contains(t, out, "Has ecosystem: true")
	assert.Contains(t, out, "PyPI/tensorflow")
}

func TestRunFilter_EcosystemNoMatch(t *testing.T) {
	resetFilterFlags()
	out, err := runCapture(t, []string{"filter", "-e", "Maven", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Has ecosystem: false")
	assert.Contains(t, out, "Matching packages (0):")
}

func TestRunFilter_EcosystemJSON(t *testing.T) {
	resetFilterFlags()
	out, err := runCapture(t, []string{"filter", "-e", "PyPI", "-o", "json", fixturePath})
	assert.Nil(t, err)
	var got map[string]any
	assert.Nil(t, json.Unmarshal([]byte(out), &got))
	assert.Equal(t, "GHSA-vxv8-r8q2-63xw", got["id"])
	assert.Equal(t, true, got["has_ecosystem"])
	affected, ok := got["affected"].([]any)
	assert.True(t, ok)
	assert.True(t, len(affected) >= 1)
}

func TestRunFilter_RefTypeText(t *testing.T) {
	resetFilterFlags()
	out, err := runCapture(t, []string{"filter", "-r", "FIX", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Reference filter: FIX")
	assert.Contains(t, out, "(0 matches)")
}

func TestRunFilter_RefTypeJSON(t *testing.T) {
	resetFilterFlags()
	out, err := runCapture(t, []string{"filter", "-r", "advisory", "-o", "json", fixturePath})
	assert.Nil(t, err)
	var got map[string]any
	assert.Nil(t, json.Unmarshal([]byte(out), &got))
	refs, ok := got["references"].([]any)
	assert.True(t, ok)
	assert.Equal(t, 1, len(refs))
}

func TestRunFilter_AliasText(t *testing.T) {
	resetFilterFlags()
	out, err := runCapture(t, []string{"filter", "-a", "CVE", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Alias filter: CVE-")
	assert.Contains(t, out, "CVE-2022-35981")
}

func TestRunFilter_AliasNoMatch(t *testing.T) {
	resetFilterFlags()
	out, err := runCapture(t, []string{"filter", "-a", "GHSA", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "(0 matches)")
}

func TestRunFilter_AliasJSON(t *testing.T) {
	resetFilterFlags()
	out, err := runCapture(t, []string{"filter", "-a", "cve", "-o", "json", fixturePath})
	assert.Nil(t, err)
	var got map[string]any
	assert.Nil(t, json.Unmarshal([]byte(out), &got))
	aliases, ok := got["aliases"].([]any)
	assert.True(t, ok)
	assert.Equal(t, 1, len(aliases))
	assert.Equal(t, "CVE-2022-35981", aliases[0])
}

func TestRunFilter_CombineAllText(t *testing.T) {
	resetFilterFlags()
	out, err := runCapture(t, []string{"filter", "-e", "PyPI", "-r", "WEB", "-a", "CVE", fixturePath})
	assert.Nil(t, err)
	assert.Contains(t, out, "Ecosystem filter: PyPI")
	assert.Contains(t, out, "Reference filter: WEB")
	assert.Contains(t, out, "Alias filter: CVE-")
}

// --- DTO 单元测试 ---

func TestToEventDTO_AllFields(t *testing.T) {
	e := &osv_schema.Event{Introduced: "1.0.0", Fixed: "1.0.1", LastAffected: "2.0.0", Limit: "v1"}
	got := toEventDTO(e)
	assert.Equal(t, "1.0.0", got.Introduced)
	assert.Equal(t, "1.0.1", got.Fixed)
	assert.Equal(t, "2.0.0", got.LastAffected)
	assert.Equal(t, "v1", got.Limit)
}

func TestToEventDTO_Empty(t *testing.T) {
	e := &osv_schema.Event{}
	got := toEventDTO(e)
	assert.Equal(t, eventDTO{}, got)
}

func TestToAffectedDTOs_WithPackageRanges(t *testing.T) {
	in := osv_schema.AffectedSlice[any, any]{
		{
			Package:  &osv_schema.Package{Ecosystem: osv_schema.EcosystemPyPI, Name: "req"},
			Versions: []string{"1.0"},
			Ranges: []*osv_schema.Range[any]{
				{Type: osv_schema.RangeTypeEcosystem, Repo: "r", Events: osv_schema.Events{
					&osv_schema.Event{Introduced: "0", Fixed: "1.0"},
				}},
			},
		},
	}
	got := toAffectedDTOs(in)
	assert.Equal(t, 1, len(got))
	assert.Equal(t, "PyPI", got[0].Package.Ecosystem)
	assert.Equal(t, "req", got[0].Package.Name)
	assert.Equal(t, []string{"1.0"}, got[0].Versions)
	assert.Equal(t, 1, len(got[0].Ranges))
	assert.Equal(t, "r", got[0].Ranges[0].Repo)
	assert.Equal(t, 1, len(got[0].Ranges[0].Events))
	assert.Equal(t, "0", got[0].Ranges[0].Events[0].Introduced)
}

func TestToAffectedDTOs_NilPackageSkipped(t *testing.T) {
	in := osv_schema.AffectedSlice[any, any]{
		{Package: nil},
	}
	got := toAffectedDTOs(in)
	assert.Equal(t, 1, len(got))
	assert.Equal(t, "", got[0].Package.Ecosystem)
}

func TestToReferenceDTOs(t *testing.T) {
	in := osv_schema.References{
		{Type: osv_schema.ReferenceTypeFix, URL: "https://f"},
		{Type: osv_schema.ReferenceTypeWeb, URL: "https://w"},
	}
	got := toReferenceDTOs(in)
	assert.Equal(t, 2, len(got))
	assert.Equal(t, "FIX", got[0].Type)
	assert.Equal(t, "https://f", got[0].URL)
}

// --- aliasPrefix 单元测试 ---

func TestAliasPrefix(t *testing.T) {
	cases := []struct{ in, want string }{
		{"CVE", "CVE-"},
		{"cve", "CVE-"},
		{"CVE-", "CVE-"},
		{"CVE-2024", "CVE-2024"},
		{"GHSA", "GHSA-"},
		{"ghsa-abc", "GHSA-ABC"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, aliasPrefix(c.in), "input %q", c.in)
	}
}
