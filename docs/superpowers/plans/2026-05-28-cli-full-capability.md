# CLI Full Capability Exposure Plan

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 将 osv_schema 库的所有 26 个公开域方法通过 CLI 暴露，补全缺失的输出字段，新增 filter 和 query 子命令

**Architecture:** 用户通过 CLI 操作 OSV JSON → parse 展示完整字段（含 verbose 标志控制）→ filter 按生态/引用类型/别名过滤 → query 提取特定子信息（severity/maven/ranges）→ 所有输出支持 text/json 格式。复用现有 `parseOsvFile` 共享函数和 `--output` 全局标志。

**Tech Stack:** Go 1.18, cobra v1.10.2, 现有 osv_schema 包

**Scope:** Medium
**Risk:** Medium

**Risks:**
- T1 修改 parse 输出格式，可能影响已有脚本解析 → 缓解：默认输出保持向后兼容，新增内容通过 `--verbose` 标志启用
- T2 filter 的 `AffectedSlice.Filter` 接收函数参数，CLI 通过字符串匹配模拟 → 缓解：用 ecosystem 名称匹配作为最常见场景

**Autonomy Level:** Full

---

### Task 1: 增强 parse 子命令 — 补全所有缺失字段，添加 --verbose 标志

**Depends on:** None
**Files:**
- Modify: `cmd/osv/parse.go:1-86`
- Modify: `cmd/osv/parse_test.go`

- [ ] **Step 1: 修改 parse.go — 添加 --verbose 标志和完整字段输出**

文件: `cmd/osv/parse.go`（替换整个文件）

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	osv_schema "github.com/scagogogo/osv-schema"
	"github.com/spf13/cobra"
)

var verbose bool

var parseCmd = &cobra.Command{
	Use:   "parse <file>",
	Short: "Parse an OSV JSON file and display key fields",
	Long:  "Parse an OSV JSON file and display its vulnerability ID, summary, severity, affected packages and other key information. Use --verbose to show all fields including dates, details, ranges, credits, and related IDs.",
	Args:  cobra.ExactArgs(1),
	RunE:  runParse,
}

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show all fields including dates, details, ranges, credits, and related IDs")
}

func runParse(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	osvData, err := parseOsvFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse OSV file %s: %w", filePath, err)
	}

	if outputFormat == "json" {
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(osvData)
	}
	return printParseText(cmd.OutOrStdout(), osvData)
}

func printParseText(w io.Writer, o *osv_schema.OsvSchema[any, any]) error {
	fmt.Fprintf(w, "ID:             %s\n", o.ID)
	fmt.Fprintf(w, "Schema Version: %s\n", o.SchemaVersion)
	fmt.Fprintf(w, "Summary:        %s\n", o.Summary)

	if verbose {
		fmt.Fprintf(w, "Published:      %s\n", o.Published.Format("2006-01-02"))
		fmt.Fprintf(w, "Modified:       %s\n", o.Modified.Format("2006-01-02"))
		if o.Withdrawn != "" {
			fmt.Fprintf(w, "Withdrawn:      %s\n", o.Withdrawn)
		}
	}

	if len(o.Aliases) > 0 {
		fmt.Fprintf(w, "Aliases:        %s\n", strings.Join(o.Aliases, ", "))
		if cve := o.Aliases.GetCVE(); cve != "" {
			fmt.Fprintf(w, "CVE:            %s\n", cve)
		}
	}

	if verbose && len(o.Related) > 0 {
		fmt.Fprintf(w, "Related:        %s\n", strings.Join(o.Related, ", "))
	}

	if verbose && o.Details != "" {
		fmt.Fprintln(w, "\nDetails:")
		fmt.Fprintln(w, o.Details)
	}

	if len(o.Severity) > 0 {
		fmt.Fprintln(w, "\nSeverity:")
		for _, s := range o.Severity {
			score, err := s.GetScoreAsFloat()
			if err != nil && verbose {
				fmt.Fprintf(w, "  %s: %s (score parse error: %v)\n", s.Type, s.Score, err)
			} else {
				fmt.Fprintf(w, "  %s: %s (score: %.1f)\n", s.Type, s.Score, score)
			}
		}
		if cvss3 := o.Severity.GetCVSS3(); cvss3 != nil {
			fmt.Fprintf(w, "  CVSS v3: %s\n", cvss3.Score)
		}
		if cvss2 := o.Severity.GetCVSS2(); cvss2 != nil {
			fmt.Fprintf(w, "  CVSS v2: %s\n", cvss2.Score)
		}
	}

	if len(o.Affected) > 0 {
		fmt.Fprintln(w, "\nAffected Packages:")
		for _, a := range o.Affected {
			if a.Package != nil {
				fmt.Fprintf(w, "  %s/%s", a.Package.Ecosystem, a.Package.Name)
				if a.Package.PackageUrl != "" {
					fmt.Fprintf(w, " (purl: %s)", a.Package.PackageUrl)
				}
				if a.Package.IsMaven() {
					fmt.Fprintf(w, " [Maven: %s:%s]", a.Package.GetGroupID(), a.Package.GetArtifactID())
				}
				if len(a.Versions) > 0 {
					fmt.Fprintf(w, " (versions: %s)", strings.Join(a.Versions, ", "))
				}
				fmt.Fprintln(w)

				if verbose && len(a.Ranges) > 0 {
					for _, r := range a.Ranges {
						fmt.Fprintf(w, "    Range [%s]", r.Type)
						if r.Repo != "" {
							fmt.Fprintf(w, " repo=%s", r.Repo)
						}
						fmt.Fprintln(w)
						for _, e := range r.Events {
							switch {
							case e.IsIntroduced():
								fmt.Fprintf(w, "      introduced: %s\n", e.Introduced)
							case e.IsFixed():
								fmt.Fprintf(w, "      fixed: %s\n", e.Fixed)
							case e.IsLastAffected():
								fmt.Fprintf(w, "      last_affected: %s\n", e.LastAffected)
							case e.IsLimit():
								fmt.Fprintf(w, "      limit: %s\n", e.Limit)
							}
						}
					}
				}

				if verbose && len(a.Severity) > 0 {
					for _, s := range a.Severity {
						fmt.Fprintf(w, "    Severity: %s: %s\n", s.Type, s.Score)
					}
				}
			}
		}
	}

	if len(o.References) > 0 {
		fmt.Fprintln(w, "\nReferences:")
		for _, r := range o.References {
			fmt.Fprintf(w, "  [%s] %s\n", r.Type, r.URL)
		}
	}

	if verbose && o.Credits != nil {
		fmt.Fprintln(w, "\nCredits:")
		fmt.Fprintf(w, "  Name: %s", o.Credits.Name)
		if o.Credits.Type != "" {
			fmt.Fprintf(w, " (%s)", o.Credits.Type)
		}
		fmt.Fprintln(w)
		if len(o.Credits.Contact) > 0 {
			fmt.Fprintf(w, "  Contact: %s\n", strings.Join(o.Credits.Contact, ", "))
		}
	}

	return nil
}
```

- [ ] **Step 2: 更新 parse_test.go — 添加 verbose 模式测试**

文件: `cmd/osv/parse_test.go`（替换整个文件）

```go
package main

import (
	"bytes"
	"os"
	"testing"
)

func TestParseCommand(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"parse", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("parse command failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("GHSA-vxv8-r8q2-63xw")) {
		t.Errorf("expected output to contain vulnerability ID, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("CVE-2022-35981")) {
		t.Errorf("expected output to contain CVE alias, got: %s", output)
	}
}

func TestParseCommandVerbose(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"parse", "-v", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("parse --verbose command failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Published:")) {
		t.Errorf("expected verbose output to contain Published date, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Range [ECOSYSTEM]")) {
		t.Errorf("expected verbose output to contain Range info, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("introduced:")) {
		t.Errorf("expected verbose output to contain range events, got: %s", output)
	}
}

func TestParseCommandJSON(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"parse", "-o", "json", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("parse command with json output failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte(`"id"`)) {
		t.Errorf("expected JSON output to contain id field, got: %s", output)
	}
}

func TestParseCommandFileNotFound(t *testing.T) {
	rootCmd.SetArgs([]string{"parse", "nonexistent.json"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}
```

- [ ] **Step 3: 验证增强后的 parse 命令**

```bash
cd /home/cc11001100/github/scagogogo/osv-schema && /usr/local/go/bin/go build -o osv ./cmd/osv/ && ./osv parse test_data/GHSA-vxv8-r8q2-63xw.json && echo "=== VERBOSE ===" && ./osv parse -v test_data/GHSA-vxv8-r8q2-63xw.json
```

Expected:
  - Exit code: 0
  - 默认输出包含 ID/Summary/Aliases/Severity/Affected/References
  - verbose 输出额外包含 Published/Modified/Range/introduced/fixed/Credits

- [ ] **Step 4: 质量门禁检查**

```bash
cd /home/cc11001100/github/scagogogo/osv-schema && /usr/local/go/bin/go vet ./... && /usr/local/go/bin/go test ./...
```

Expected:
  - Exit code: 0

- [ ] **Step 5: 提交**

```bash
cd /home/cc11001100/github/scagogogo/osv-schema && rm -f osv && git add cmd/osv/parse.go cmd/osv/parse_test.go && git -c user.name="CC11001100" -c user.email="CC11001100@qq.com" commit -m "$(cat <<'EOF'
feat(cli): enhance parse with --verbose flag, full field output, and library method coverage

- Add --verbose/-v flag to show dates, details, ranges, credits, related IDs
- Display Package.Purl and Maven GroupID/ArtifactID via Package methods
- Show SeveritySlice.GetCVSS3/GetCVSS2 shortcut methods
- Display Range events using Event.Is* methods
- Use Severity.GetScoreAsFloat for proper error reporting in verbose mode
- Backward compatible: default output unchanged, new content only with -v

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

---

### Task 2: 实现 filter 子命令 — 按生态/引用类型过滤

**Depends on:** Task 1
**Files:**
- Create: `cmd/osv/filter.go`
- Create: `cmd/osv/filter_test.go`

- [ ] **Step 1: 创建 filter 子命令 — 按生态过滤受影响包、按类型过滤引用**

```go
// cmd/osv/filter.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	osv_schema "github.com/scagogogo/osv-schema"
	"github.com/spf13/cobra"
)

var (
	filterEcosystem string
	filterRefType   string
	filterAlias    string
)

var filterCmd = &cobra.Command{
	Use:   "filter <file>",
	Short: "Filter OSV data by ecosystem, reference type, or alias pattern",
	Long:  "Filter an OSV JSON file to show only matching affected packages (by ecosystem), references (by type), or aliases (by pattern). Uses AffectedSlice.FilterByEcosystem, References.FilterByType, and Aliases.Filter methods.",
	Args:  cobra.ExactArgs(1),
	RunE:  runFilter,
}

func init() {
	rootCmd.AddCommand(filterCmd)
	filterCmd.Flags().StringVarP(&filterEcosystem, "ecosystem", "e", "", "Filter affected packages by ecosystem (e.g. PyPI, Maven, npm)")
	filterCmd.Flags().StringVarP(&filterRefType, "ref-type", "r", "", "Filter references by type (ADVISORY, FIX, ARTICLE, REPORT, WEB, PACKAGE, etc.)")
	filterCmd.Flags().StringVarP(&filterAlias, "alias", "a", "", "Filter aliases by substring match (e.g. CVE, GHSA)")
}

func runFilter(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	osvData, err := parseOsvFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse OSV file %s: %w", filePath, err)
	}

	out := cmd.OutOrStdout()
	result := &filterResult{ID: osvData.ID}

	hasFilter := false

	if filterEcosystem != "" {
		hasFilter = true
		ecosystem := osv_schema.Ecosystem(filterEcosystem)
		result.HasEcosystem = osvData.Affected.HasEcosystem(ecosystem)
		result.FilteredAffected = osvData.Affected.FilterByEcosystem(ecosystem)
	}

	if filterRefType != "" {
		hasFilter = true
		refType := osv_schema.ReferenceType(strings.ToUpper(filterRefType))
		result.FilteredReferences = osvData.References.FilterByType(refType)
	}

	if filterAlias != "" {
		hasFilter = true
		result.FilteredAliases = osvData.Aliases.Filter(func(alias string) bool {
			return strings.Contains(strings.ToUpper(alias), strings.ToUpper(filterAlias))
		})
	}

	if !hasFilter {
		return fmt.Errorf("at least one filter flag is required (--ecosystem, --ref-type, or --alias)")
	}

	if outputFormat == "json" {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}
	return printFilterText(out, result)
}

type filterResult struct {
	ID                 string                               `json:"id"`
	HasEcosystem       bool                                 `json:"has_ecosystem,omitempty"`
	FilteredAffected   osv_schema.AffectedSlice[any, any]   `json:"filtered_affected,omitempty"`
	FilteredReferences osv_schema.References                 `json:"filtered_references,omitempty"`
	FilteredAliases    osv_schema.Aliases                    `json:"filtered_aliases,omitempty"`
}

func printFilterText(w io.Writer, r *filterResult) error {
	fmt.Fprintf(w, "ID: %s\n", r.ID)

	if r.HasEcosystem || len(r.FilteredAffected) > 0 {
		fmt.Fprintf(w, "Has Ecosystem: %v\n", r.HasEcosystem)
		if len(r.FilteredAffected) > 0 {
			fmt.Fprintln(w, "Filtered Affected Packages:")
			for _, a := range r.FilteredAffected {
				if a.Package != nil {
					fmt.Fprintf(w, "  %s/%s", a.Package.Ecosystem, a.Package.Name)
					if a.Package.IsMaven() {
						fmt.Fprintf(w, " [%s:%s]", a.Package.GetGroupID(), a.Package.GetArtifactID())
					}
					fmt.Fprintln(w)
				}
			}
		} else {
			fmt.Fprintln(w, "  (no matching packages)")
		}
	}

	if len(r.FilteredReferences) > 0 {
		fmt.Fprintln(w, "Filtered References:")
		for _, ref := range r.FilteredReferences {
			fmt.Fprintf(w, "  [%s] %s\n", ref.Type, ref.URL)
		}
	} else if filterRefType != "" {
		fmt.Fprintln(w, "  (no matching references)")
	}

	if len(r.FilteredAliases) > 0 {
		fmt.Fprintf(w, "Filtered Aliases: %s\n", strings.Join(r.FilteredAliases, ", "))
	} else if filterAlias != "" {
		fmt.Fprintln(w, "  (no matching aliases)")
	}

	return nil
}
```

- [ ] **Step 2: 创建 filter 测试**

```go
// cmd/osv/filter_test.go
package main

import (
	"bytes"
	"os"
	"testing"
)

func TestFilterByEcosystem(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"filter", "-e", "PyPI", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("filter --ecosystem command failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Has Ecosystem: true")) {
		t.Errorf("expected Has Ecosystem true for PyPI, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("PyPI/tensorflow")) {
		t.Errorf("expected PyPI/tensorflow in filtered output, got: %s", output)
	}
}

func TestFilterByEcosystemNotFound(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"filter", "-e", "Maven", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("filter --ecosystem Maven command failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Has Ecosystem: false")) {
		t.Errorf("expected Has Ecosystem false for Maven, got: %s", output)
	}
}

func TestFilterByRefType(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"filter", "-r", "ADVISORY", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("filter --ref-type command failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("nvd.nist.gov")) {
		t.Errorf("expected ADVISORY reference in filtered output, got: %s", output)
	}
}

func TestFilterByAlias(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"filter", "-a", "CVE", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("filter --alias command failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("CVE-2022-35981")) {
		t.Errorf("expected CVE alias in filtered output, got: %s", output)
	}
}

func TestFilterNoFlags(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	rootCmd.SetArgs([]string{"filter", testDataPath})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no filter flags provided, got nil")
	}
}
```

- [ ] **Step 3: 验证 filter 命令**

```bash
cd /home/cc11001100/github/scagogogo/osv-schema && /usr/local/go/bin/go build -o osv ./cmd/osv/ && ./osv filter -e PyPI test_data/GHSA-vxv8-r8q2-63xw.json && echo "---" && ./osv filter -r FIX test_data/GHSA-vxv8-r8q2-63xw.json && echo "---" && ./osv filter -a CVE test_data/GHSA-vxv8-r8q2-63xw.json
```

Expected:
  - Exit code: 0
  - ecosystem 过滤输出包含 PyPI 包
  - ref-type 过滤输出包含匹配的引用
  - alias 过滤输出包含 CVE

- [ ] **Step 4: 质量门禁检查**

```bash
cd /home/cc11001100/github/scagogogo/osv-schema && /usr/local/go/bin/go vet ./... && /usr/local/go/bin/go test ./...
```

Expected:
  - Exit code: 0

- [ ] **Step 5: 提交**

```bash
cd /home/cc11001100/github/scagogogo/osv-schema && rm -f osv && git add cmd/osv/filter.go cmd/osv/filter_test.go && git -c user.name="CC11001100" -c user.email="CC11001100@qq.com" commit -m "$(cat <<'EOF'
feat(cli): add filter subcommand with ecosystem/reference/alias filtering

Uses AffectedSlice.FilterByEcosystem/HasEcosystem, References.FilterByType,
and Aliases.Filter to expose library filtering capabilities.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

---

### Task 3: 实现 query 子命令 — 提取特定子信息

**Depends on:** Task 1
**Files:**
- Create: `cmd/osv/query.go`
- Create: `cmd/osv/query_test.go`

- [ ] **Step 1: 创建 query 子命令 — 提取 severity、maven、ranges 等子信息**

```go
// cmd/osv/query.go
package main

import (
	"encoding/json"
	"fmt"
	"io"

	osv_schema "github.com/scagogogo/osv-schema"
	"github.com/spf13/cobra"
)

var (
	querySeverity string // cvss2 or cvss3
	queryMaven   bool
	queryRanges  bool
	queryEvents  bool
)

var queryCmd = &cobra.Command{
	Use:   "query <file>",
	Short: "Query specific sub-information from an OSV JSON file",
	Long:  "Extract specific sub-information from an OSV JSON file: severity details (CVSS v2/v3), Maven package decomposition, or range/event data. Uses SeveritySlice.GetCVSS3/GetCVSS2, Package.IsMaven/GetGroupID/GetArtifactID, and Event.Is* methods.",
	Args:  cobra.ExactArgs(1),
	RunE:  runQuery,
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVar(&querySeverity, "severity", "", "Extract severity info: cvss2 or cvss3")
	queryCmd.Flags().BoolVar(&queryMaven, "maven", false, "Show Maven package decomposition (GroupID/ArtifactID)")
	queryCmd.Flags().BoolVar(&queryRanges, "ranges", false, "Show version ranges for all affected packages")
	queryCmd.Flags().BoolVar(&queryEvents, "events", false, "Show event details (introduced/fixed/last_affected/limit) for all ranges")
}

func runQuery(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	osvData, err := parseOsvFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse OSV file %s: %w", filePath, err)
	}

	out := cmd.OutOrStdout()

	hasQuery := false

	if querySeverity != "" {
		hasQuery = true
		return printSeverityQuery(out, osvData)
	}
	if queryMaven {
		hasQuery = true
		printMavenQuery(out, osvData)
	}
	if queryRanges {
		hasQuery = true
		printRangesQuery(out, osvData, false)
	}
	if queryEvents {
		hasQuery = true
		printRangesQuery(out, osvData, true)
	}

	if !hasQuery {
		return fmt.Errorf("at least one query flag is required (--severity, --maven, --ranges, or --events)")
	}
	return nil
}

type severityQueryResult struct {
	ID       string                `json:"id"`
	CVSS2    *osv_schema.Severity  `json:"cvss2,omitempty"`
	CVSS3    *osv_schema.Severity  `json:"cvss3,omitempty"`
	ScoreV2  float64               `json:"score_v2,omitempty"`
	ScoreV3  float64               `json:"score_v3,omitempty"`
}

func printSeverityQuery(w io.Writer, o *osv_schema.OsvSchema[any, any]) error {
	result := &severityQueryResult{ID: o.ID}

	if cvss3 := o.Severity.GetCVSS3(); cvss3 != nil {
		result.CVSS3 = cvss3
		result.ScoreV3 = cvss3.GetScore()
	}
	if cvss2 := o.Severity.GetCVSS2(); cvss2 != nil {
		result.CVSS2 = cvss2
		result.ScoreV2 = cvss2.GetScore()
	}

	if outputFormat == "json" {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	fmt.Fprintf(w, "ID: %s\n", result.ID)
	if result.CVSS3 != nil {
		fmt.Fprintf(w, "CVSS v3: %s (score: %.1f)\n", result.CVSS3.Score, result.ScoreV3)
	} else {
		fmt.Fprintln(w, "CVSS v3: (none)")
	}
	if result.CVSS2 != nil {
		fmt.Fprintf(w, "CVSS v2: %s (score: %.1f)\n", result.CVSS2.Score, result.ScoreV2)
	} else {
		fmt.Fprintln(w, "CVSS v2: (none)")
	}
	return nil
}

func printMavenQuery(w io.Writer, o *osv_schema.OsvSchema[any, any]) {
	fmt.Fprintf(w, "ID: %s\n", o.ID)
	found := false
	for _, a := range o.Affected {
		if a.Package != nil && a.Package.IsMaven() {
			found = true
			fmt.Fprintf(w, "  Maven Package: %s\n", a.Package.Name)
			fmt.Fprintf(w, "    GroupID:    %s\n", a.Package.GetGroupID())
			fmt.Fprintf(w, "    ArtifactID: %s\n", a.Package.GetArtifactID())
		}
	}
	if !found {
		fmt.Fprintln(w, "  (no Maven packages found)")
	}
}

func printRangesQuery(w io.Writer, o *osv_schema.OsvSchema[any, any], showEvents bool) {
	fmt.Fprintf(w, "ID: %s\n", o.ID)
	for _, a := range o.Affected {
		if a.Package == nil {
			continue
		}
		if len(a.Ranges) == 0 {
			continue
		}
		fmt.Fprintf(w, "  %s/%s:\n", a.Package.Ecosystem, a.Package.Name)
		for _, r := range a.Ranges {
			fmt.Fprintf(w, "    Range Type: %s", r.Type)
			if r.Repo != "" {
				fmt.Fprintf(w, " (repo: %s)", r.Repo)
			}
			fmt.Fprintln(w)

			if showEvents {
				for _, e := range r.Events {
					switch {
					case e.IsIntroduced():
						fmt.Fprintf(w, "      introduced: %s\n", e.Introduced)
					case e.IsFixed():
						fmt.Fprintf(w, "      fixed: %s\n", e.Fixed)
					case e.IsLastAffected():
						fmt.Fprintf(w, "      last_affected: %s\n", e.LastAffected)
					case e.IsLimit():
						fmt.Fprintf(w, "      limit: %s\n", e.Limit)
					}
				}
			} else {
				fmt.Fprintf(w, "      %d events\n", len(r.Events))
			}
		}
	}
}
```

- [ ] **Step 2: 创建 query 测试**

```go
// cmd/osv/query_test.go
package main

import (
	"bytes"
	"os"
	"testing"
)

func TestQuerySeverityCVSS3(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"query", "--severity", "cvss3", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("query --severity cvss3 failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("CVSS v3:")) {
		t.Errorf("expected CVSS v3 in output, got: %s", output)
	}
}

func TestQueryMaven(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"query", "--maven", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("query --maven failed: %v", err)
	}
	output := buf.String()
	// 测试数据没有 Maven 包，应显示 "no Maven packages found"
	if !bytes.Contains([]byte(output), []byte("no Maven packages found")) {
		t.Errorf("expected no Maven packages message, got: %s", output)
	}
}

func TestQueryRanges(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"query", "--ranges", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("query --ranges failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Range Type: ECOSYSTEM")) {
		t.Errorf("expected ECOSYSTEM range type, got: %s", output)
	}
}

func TestQueryEvents(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"query", "--events", testDataPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("query --events failed: %v", err)
	}
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("introduced:")) {
		t.Errorf("expected introduced event, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("fixed:")) {
		t.Errorf("expected fixed event, got: %s", output)
	}
}

func TestQueryNoFlags(t *testing.T) {
	testDataPath := "../../test_data/GHSA-vxv8-r8q2-63xw.json"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skip("Test data file not found")
	}

	rootCmd.SetArgs([]string{"query", testDataPath})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no query flags provided, got nil")
	}
}
```

- [ ] **Step 3: 验证 query 命令**

```bash
cd /home/cc11001100/github/scagogogo/osv-schema && /usr/local/go/bin/go build -o osv ./cmd/osv/ && ./osv query --severity cvss3 test_data/GHSA-vxv8-r8q2-63xw.json && echo "---" && ./osv query --ranges test_data/GHSA-vxv8-r8q2-63xw.json && echo "---" && ./osv query --events test_data/GHSA-vxv8-r8q2-63xw.json
```

Expected:
  - Exit code: 0
  - severity 输出包含 CVSS v3 信息
  - ranges 输出包含 Range Type
  - events 输出包含 introduced/fixed 事件

- [ ] **Step 4: 质量门禁检查**

```bash
cd /home/cc11001100/github/scagogogo/osv-schema && /usr/local/go/bin/go vet ./... && /usr/local/go/bin/go test ./...
```

Expected:
  - Exit code: 0

- [ ] **Step 5: 提交**

```bash
cd /home/cc11001100/github/scagogogo/osv-schema && rm -f osv && git add cmd/osv/query.go cmd/osv/query_test.go && git -c user.name="CC11001100" -c user.email="CC11001100@qq.com" commit -m "$(cat <<'EOF'
feat(cli): add query subcommand for severity, maven, ranges, and events

Uses SeveritySlice.GetCVSS3/GetCVSS2, Package.IsMaven/GetGroupID/GetArtifactID,
and Event.IsIntroduced/IsFixed/IsLastAffected/IsLimit methods.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```
