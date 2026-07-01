package main

import (
	"encoding/json"
	"fmt"
	"io"

	osv_schema "github.com/scagogogo/osv-schema-skills"
	"github.com/spf13/cobra"
)

var (
	querySeverity string
	queryMaven    bool
	queryRanges   bool
	queryEvents   bool
)

var queryCmd = &cobra.Command{
	Use:   "query <file>",
	Short: "Extract specific sub-information from OSV JSON data",
	Long: `Extract focused sub-information from an OSV JSON file:
  --severity cvss3|cvss2   CVSS severity entry and parsed score
  --maven                  Maven groupId:artifactId decomposition
  --ranges                 Version ranges per affected package
  --events                 Event timeline (introduced/fixed/last_affected/limit)

At least one flag is required. Flags can be combined.`,
	Args: cobra.ExactArgs(1),
	RunE: runQuery,
}

func init() {
	queryCmd.Flags().StringVar(&querySeverity, "severity", "", "query severity: cvss3 or cvss2")
	queryCmd.Flags().BoolVar(&queryMaven, "maven", false, "decompose Maven package names into groupId/artifactId")
	queryCmd.Flags().BoolVar(&queryRanges, "ranges", false, "show version ranges")
	queryCmd.Flags().BoolVar(&queryEvents, "events", false, "show event timeline")
	rootCmd.AddCommand(queryCmd)
}

func runQuery(cmd *cobra.Command, args []string) error {
	if querySeverity == "" && !queryMaven && !queryRanges && !queryEvents {
		return fmt.Errorf("at least one query flag is required (--severity, --maven, --ranges, or --events)")
	}

	osvData, err := parseOsvFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse OSV file %s: %w", args[0], err)
	}

	out := cmd.OutOrStdout()

	if outputFormat == "json" {
		return printQueryJSON(out, osvData)
	}
	return printQueryText(out, osvData)
}

func printQueryJSON(w io.Writer, o *osv_schema.OsvSchema[any, any]) error {
	result := map[string]any{
		"id": o.ID,
	}
	if querySeverity != "" {
		var s *osv_schema.Severity
		switch querySeverity {
		case "cvss3":
			s = o.Severity.GetCVSS3()
		case "cvss2":
			s = o.Severity.GetCVSS2()
		default:
			return fmt.Errorf("invalid --severity value %q (use cvss2 or cvss3)", querySeverity)
		}
		result["severity"] = toSeverityDTO(s)
	}
	if queryMaven {
		result["maven"] = collectMaven(o)
	}
	if queryRanges {
		result["ranges"] = collectRangesDTO(o)
	}
	if queryEvents {
		result["events"] = collectEventsDTO(o)
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// --- 干净的 JSON DTO（omitempty，对 AI Agent 友好）---

type severityDTO struct {
	Type         string  `json:"type,omitempty"`
	Score        string  `json:"score,omitempty"`
	NumericScore float64 `json:"numeric_score,omitempty"`
}

type mavenDTO struct {
	Name       string `json:"name,omitempty"`
	GroupID    string `json:"group_id,omitempty"`
	ArtifactID string `json:"artifact_id,omitempty"`
}

type rangeItemDTO struct {
	Package string     `json:"package,omitempty"`
	Type    string     `json:"type"`
	Repo    string     `json:"repo,omitempty"`
	Events  []eventDTO `json:"events,omitempty"`
}

type eventItemDTO struct {
	Package string `json:"package,omitempty"`
	eventDTO
}

func toSeverityDTO(s *osv_schema.Severity) *severityDTO {
	if s == nil {
		return nil
	}
	return &severityDTO{
		Type:         string(s.Type),
		Score:        s.Score,
		NumericScore: s.GetScore(),
	}
}

func collectMaven(o *osv_schema.OsvSchema[any, any]) []mavenDTO {
	out := make([]mavenDTO, 0)
	for _, a := range o.Affected {
		if a.Package == nil || !a.Package.IsMaven() {
			continue
		}
		out = append(out, mavenDTO{
			Name:       a.Package.Name,
			GroupID:    a.Package.GetGroupID(),
			ArtifactID: a.Package.GetArtifactID(),
		})
	}
	return out
}

func collectRangesDTO(o *osv_schema.OsvSchema[any, any]) []rangeItemDTO {
	out := make([]rangeItemDTO, 0)
	for _, a := range o.Affected {
		if a.Package == nil {
			continue
		}
		pkg := fmt.Sprintf("%s/%s", a.Package.Ecosystem, a.Package.Name)
		for _, r := range a.Ranges {
			dto := rangeItemDTO{Package: pkg, Type: string(r.Type), Repo: r.Repo}
			for _, e := range r.Events {
				dto.Events = append(dto.Events, toEventDTO(e))
			}
			out = append(out, dto)
		}
	}
	return out
}

func collectEventsDTO(o *osv_schema.OsvSchema[any, any]) []eventItemDTO {
	out := make([]eventItemDTO, 0)
	for _, a := range o.Affected {
		if a.Package == nil {
			continue
		}
		pkg := fmt.Sprintf("%s/%s", a.Package.Ecosystem, a.Package.Name)
		for _, r := range a.Ranges {
			for _, e := range r.Events {
				out = append(out, eventItemDTO{Package: pkg, eventDTO: toEventDTO(e)})
			}
		}
	}
	return out
}

func printQueryText(w io.Writer, o *osv_schema.OsvSchema[any, any]) error {
	fmt.Fprintf(w, "ID: %s\n\n", o.ID)

	if querySeverity != "" {
		var s *osv_schema.Severity
		switch querySeverity {
		case "cvss3":
			s = o.Severity.GetCVSS3()
		case "cvss2":
			s = o.Severity.GetCVSS2()
		default:
			return fmt.Errorf("invalid --severity value %q (use cvss2 or cvss3)", querySeverity)
		}
		fmt.Fprintf(w, "Severity (%s):\n", querySeverity)
		if s == nil {
			fmt.Fprintln(w, "  (none)")
		} else {
			fmt.Fprintf(w, "  Type:  %s\n", s.Type)
			fmt.Fprintf(w, "  Score: %s\n", s.Score)
			fmt.Fprintf(w, "  Numeric score: %.1f\n", s.GetScore())
		}
		fmt.Fprintln(w)
	}

	if queryMaven {
		fmt.Fprintln(w, "Maven decomposition:")
		mavens := collectMaven(o)
		if len(mavens) == 0 {
			fmt.Fprintln(w, "  (no Maven packages)")
		}
		for _, m := range mavens {
			fmt.Fprintf(w, "  %s → groupId=%s, artifactId=%s\n", m.Name, m.GroupID, m.ArtifactID)
		}
		fmt.Fprintln(w)
	}

	if queryRanges {
		fmt.Fprintln(w, "Version ranges:")
		for _, r := range collectRangesDTO(o) {
			fmt.Fprintf(w, "  %s", r.Package)
			if r.Repo != "" {
				fmt.Fprintf(w, " (repo: %s)", r.Repo)
			}
			fmt.Fprintf(w, " [%s]\n", r.Type)
			for _, ev := range r.Events {
				printEventDTO(w, "    ", ev)
			}
		}
		fmt.Fprintln(w)
	}

	if queryEvents {
		fmt.Fprintln(w, "Event timeline:")
		for _, ev := range collectEventsDTO(o) {
			fmt.Fprintf(w, "  %s: ", ev.Package)
			printEventDTOInline(w, ev.eventDTO)
			fmt.Fprintln(w)
		}
	}
	return nil
}

// printEventDTO 把单个事件按非空字段逐行输出
func printEventDTO(w io.Writer, prefix string, e eventDTO) {
	if e.Introduced != "" {
		fmt.Fprintf(w, "%sintroduced: %s\n", prefix, e.Introduced)
	}
	if e.Fixed != "" {
		fmt.Fprintf(w, "%sfixed: %s\n", prefix, e.Fixed)
	}
	if e.LastAffected != "" {
		fmt.Fprintf(w, "%slast_affected: %s\n", prefix, e.LastAffected)
	}
	if e.Limit != "" {
		fmt.Fprintf(w, "%slimit: %s\n", prefix, e.Limit)
	}
}

// printEventDTOInline 把单个事件按非空字段行内输出（key=value 形式）
func printEventDTOInline(w io.Writer, e eventDTO) {
	if e.Introduced != "" {
		fmt.Fprintf(w, "introduced=%s ", e.Introduced)
	}
	if e.Fixed != "" {
		fmt.Fprintf(w, "fixed=%s ", e.Fixed)
	}
	if e.LastAffected != "" {
		fmt.Fprintf(w, "last_affected=%s ", e.LastAffected)
	}
	if e.Limit != "" {
		fmt.Fprintf(w, "limit=%s ", e.Limit)
	}
}

