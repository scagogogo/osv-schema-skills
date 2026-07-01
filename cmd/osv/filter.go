package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	osv_schema "github.com/scagogogo/osv-schema-skills"
	"github.com/spf13/cobra"
)

var (
	filterEcosystem string
	filterRefType   string
	filterAlias     string
)

var filterCmd = &cobra.Command{
	Use:   "filter <file>",
	Short: "Filter OSV JSON data by ecosystem, reference type, or alias pattern",
	Long: `Filter an OSV JSON file by:
  - affected package ecosystem   (-e, --ecosystem)
  - reference type               (-r, --ref-type)
  - alias pattern                (-a, --alias)

At least one filter flag is required. Filters can be combined.`,
	Args: cobra.ExactArgs(1),
	RunE: runFilter,
}

func init() {
	filterCmd.Flags().StringVarP(&filterEcosystem, "ecosystem", "e", "", "filter affected packages by ecosystem (e.g. PyPI, npm, Maven)")
	filterCmd.Flags().StringVarP(&filterRefType, "ref-type", "r", "", "filter references by type (e.g. ADVISORY, FIX)")
	filterCmd.Flags().StringVarP(&filterAlias, "alias", "a", "", "filter aliases by prefix pattern (e.g. CVE, GHSA)")
	rootCmd.AddCommand(filterCmd)
}

func runFilter(cmd *cobra.Command, args []string) error {
	if filterEcosystem == "" && filterRefType == "" && filterAlias == "" {
		return fmt.Errorf("at least one filter flag is required (--ecosystem, --ref-type, or --alias)")
	}

	osvData, err := parseOsvFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse OSV file %s: %w", args[0], err)
	}

	out := cmd.OutOrStdout()

	if outputFormat == "json" {
		return printFilterJSON(out, osvData)
	}
	return printFilterText(out, osvData)
}

func printFilterJSON(w io.Writer, o *osv_schema.OsvSchema[any, any]) error {
	result := map[string]any{
		"id": o.ID,
	}
	if filterEcosystem != "" {
		eco := osv_schema.Ecosystem(filterEcosystem)
		result["has_ecosystem"] = o.Affected.HasEcosystem(eco)
		result["affected"] = toAffectedDTOs(o.Affected.FilterByEcosystem(eco))
	}
	if filterRefType != "" {
		rt := osv_schema.ReferenceType(strings.ToUpper(filterRefType))
		result["references"] = toReferenceDTOs(o.References.FilterByType(rt))
	}
	if filterAlias != "" {
		prefix := aliasPrefix(filterAlias)
		result["aliases"] = o.Aliases.Filter(func(alias string) bool {
			return strings.HasPrefix(strings.ToUpper(alias), prefix)
		})
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// --- 干净的 JSON DTO（omitempty，避免空字段污染输出，对 AI Agent 更友好）---

type eventDTO struct {
	Introduced   string `json:"introduced,omitempty"`
	Fixed        string `json:"fixed,omitempty"`
	LastAffected string `json:"last_affected,omitempty"`
	Limit        string `json:"limit,omitempty"`
}

type rangeDTO struct {
	Type   string     `json:"type"`
	Repo   string     `json:"repo,omitempty"`
	Events []eventDTO `json:"events,omitempty"`
}

type packageDTO struct {
	Ecosystem string `json:"ecosystem,omitempty"`
	Name      string `json:"name,omitempty"`
	Purl      string `json:"purl,omitempty"`
}

type affectedDTO struct {
	Package  packageDTO  `json:"package,omitempty"`
	Ranges   []rangeDTO  `json:"ranges,omitempty"`
	Versions []string    `json:"versions,omitempty"`
}

type referenceDTO struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}

func toEventDTO(e *osv_schema.Event) eventDTO {
	return eventDTO{
		Introduced:   e.Introduced,
		Fixed:        e.Fixed,
		LastAffected: e.LastAffected,
		Limit:        e.Limit,
	}
}

func toAffectedDTOs(in osv_schema.AffectedSlice[any, any]) []affectedDTO {
	out := make([]affectedDTO, 0, len(in))
	for _, a := range in {
		dto := affectedDTO{}
		if a.Package != nil {
			dto.Package = packageDTO{
				Ecosystem: string(a.Package.Ecosystem),
				Name:      a.Package.Name,
				Purl:      a.Package.PackageUrl,
			}
		}
		dto.Versions = a.Versions
		for _, r := range a.Ranges {
			rdto := rangeDTO{Type: string(r.Type), Repo: r.Repo}
			for _, e := range r.Events {
				rdto.Events = append(rdto.Events, toEventDTO(e))
			}
			dto.Ranges = append(dto.Ranges, rdto)
		}
		out = append(out, dto)
	}
	return out
}

func toReferenceDTOs(in osv_schema.References) []referenceDTO {
	out := make([]referenceDTO, 0, len(in))
	for _, r := range in {
		out = append(out, referenceDTO{Type: string(r.Type), URL: r.URL})
	}
	return out
}

func printFilterText(w io.Writer, o *osv_schema.OsvSchema[any, any]) error {
	fmt.Fprintf(w, "ID: %s\n\n", o.ID)

	if filterEcosystem != "" {
		eco := osv_schema.Ecosystem(filterEcosystem)
		fmt.Fprintf(w, "Ecosystem filter: %s\n", filterEcosystem)
		fmt.Fprintf(w, "  Has ecosystem: %v\n", o.Affected.HasEcosystem(eco))
		filtered := o.Affected.FilterByEcosystem(eco)
		fmt.Fprintf(w, "  Matching packages (%d):\n", len(filtered))
		for _, a := range filtered {
			if a.Package != nil {
				fmt.Fprintf(w, "    - %s/%s\n", a.Package.Ecosystem, a.Package.Name)
			}
		}
		fmt.Fprintln(w)
	}

	if filterRefType != "" {
		rt := osv_schema.ReferenceType(strings.ToUpper(filterRefType))
		filtered := o.References.FilterByType(rt)
		fmt.Fprintf(w, "Reference filter: %s (%d matches)\n", rt, len(filtered))
		for _, r := range filtered {
			fmt.Fprintf(w, "  [%s] %s\n", r.Type, r.URL)
		}
		fmt.Fprintln(w)
	}

	if filterAlias != "" {
		prefix := aliasPrefix(filterAlias)
		filtered := o.Aliases.Filter(func(alias string) bool {
			return strings.HasPrefix(strings.ToUpper(alias), prefix)
		})
		fmt.Fprintf(w, "Alias filter: %s* (%d matches)\n", prefix, len(filtered))
		for _, a := range filtered {
			fmt.Fprintf(w, "  - %s\n", a)
		}
	}
	return nil
}

// aliasPrefix 把 "CVE" 归一化为 "CVE-"，让用户可传 CVE 或 CVE-2024 前缀
func aliasPrefix(pattern string) string {
	p := strings.ToUpper(pattern)
	if !strings.HasSuffix(p, "-") && !strings.Contains(p, "-") {
		return p + "-"
	}
	return p
}
