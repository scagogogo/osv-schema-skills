---
name: osv-severity
description: Analyze CVSS severity data from OSV vulnerability records. Triggers on mentions of CVSS scores, vulnerability severity assessment, risk rating, or when user needs to evaluate the severity/impact of a vulnerability.
allowed-tools: "Bash(osv:*)"
argument-hint: <osv-json-file>
---

# OSV Severity Analysis

> **Setup:** See `/osv-installation` for one-time CLI/SDK install.
> **Layers:** SDK (Go) → CLI (shell) — pick your entry point.

## When to Use

- Assess the severity of a vulnerability using CVSS scores
- Compare CVSS v2 and v3 scores for the same vulnerability
- Determine risk level (critical/high/medium/low) from CVSS scores
- Extract CVSS vector strings for detailed impact analysis

## Decision Tree

```
Severity data → what do you need?
├─ Quick CVSS score?           → GetCVSS3().GetScore() / osv query --severity
├─ Both v2 and v3?             → GetCVSS2() + GetCVSS3()
├─ Raw score string?           → Severity.Score field
├─ Severity type?              → Severity.Type (CVSS_V2, CVSS_V3)
└─ Per-package severity?       → Affected[i].Severity
```

## Task Patterns

### Get CVSS v3 score

**Goal:** What is the CVSS v3 score of this vulnerability?

| Layer | Approach |
|-------|----------|
| CLI | `osv query --severity cvss3 vulnerability.json` |
| SDK | `osvData.Severity.GetCVSS3().GetScore()` |

### Compare CVSS v2 and v3

**Goal:** Compare CVSS v2 vs v3 scores.

```go
cvss3 := osvData.Severity.GetCVSS3()
cvss2 := osvData.Severity.GetCVSS2()
if cvss3 != nil {
    fmt.Printf("CVSS v3: %.1f\n", cvss3.GetScore())
}
if cvss2 != nil {
    fmt.Printf("CVSS v2: %.1f\n", cvss2.GetScore())
}
```

### Risk level classification

**Goal:** Classify vulnerability risk level from CVSS score.

```go
if cvss3 := osvData.Severity.GetCVSS3(); cvss3 != nil {
    score := cvss3.GetScore()
    switch {
    case score >= 9.0:
        fmt.Println("CRITICAL")
    case score >= 7.0:
        fmt.Println("HIGH")
    case score >= 4.0:
        fmt.Println("MEDIUM")
    default:
        fmt.Println("LOW")
    }
}
```

### Per-package severity

**Goal:** Some affected packages may have their own severity entries.

```go
for _, affected := range osvData.Affected {
    for _, sev := range affected.Severity {
        fmt.Printf("  %s: %s (score: %.1f)\n", sev.Type, sev.Score, sev.GetScore())
    }
}
```

## API Reference

### SDK — SeverityType Constants

```go
SeverityTypeCVSSv2 SeverityType = "CVSS_V2"
SeverityTypeCVSSv3 SeverityType = "CVSS_V3"
```

### SDK — Severity Struct

```go
type Severity struct {
    Type  SeverityType `json:"type"`  // CVSS_V2 or CVSS_V3
    Score string       `json:"score"` // CVSS vector string or numeric score
}
```

### SDK — SeveritySlice Methods

```go
func (s SeveritySlice) GetCVSS3() *Severity  // first CVSS_V3 entry
func (s SeveritySlice) GetCVSS2() *Severity  // first CVSS_V2 entry
```

### SDK — Severity Methods

```go
func (s *Severity) GetScore() float64           // score as float64 (0.0 on error)
func (s *Severity) GetScoreAsFloat() (float64, error)  // score with error
func (s *Severity) GetScoreAsPointer() *float64        // score as pointer
```

### CLI Commands

```bash
osv query --severity cvss3 <file>   # CVSS v3 score and details
osv query --severity cvss2 <file>   # CVSS v2 score and details
osv parse -v <file>                 # Verbose parse includes severity
```

## Cross-References

- [[osv-query]] — general querying including severity
- [[osv-parse]] — full parsing with severity display
- [[osv-affected]] — per-package severity analysis
- [[osv-filter]] — filter by ecosystem then check severity

## Important Notes

- `GetScore()` parses the `Score` string field — it may contain a CVSS vector string (e.g., `CVSS:3.1/AV:N/AC:L/...`) or a numeric value
- For CVSS vector strings, `GetScore()` extracts the base score; the full vector provides more detail (attack vector, complexity, etc.)
- A vulnerability may have both CVSS v2 and v3 scores — they can differ significantly
- `SeveritySlice` may be empty — always check with `len()` or `GetCVSS3() != nil`
- Per-package severity (`Affected[i].Severity`) is separate from top-level severity
