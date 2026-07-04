# osv-severity

Analyze CVSS severity data from OSV records.

> **Trigger:** mentions of CVSS scores, vulnerability severity assessment, risk rating, or evaluating impact.
> **Skill source:** [`.claude/skills/osv-severity/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-severity/SKILL.md)

## CLI

Severity is queried via `osv query`:

```bash
osv query --severity cvss3 vulnerability.json  # CVSS v3 entry + parsed score (0.0 on a vector string)
osv query --severity cvss2 vulnerability.json  # CVSS v2
```

Or see all severities at once with `osv parse -v`.

## SDK

```go
// CVSS v3 entry (nil if absent — check before use)
s := v.Severity.GetCVSS3()
if s == nil {
    return // no CVSS v3 entry
}

// Parsed numeric score
fmt.Println(s.GetScore())        // float64, 0.0 if unparseable
score, err := s.GetScoreAsFloat() // with error
ptr := s.GetScoreAsPointer()     // *float64, nil on error
```

## CVSS score table

| Score range | Severity |
|-------------|----------|
| 0.1–3.9 | Low |
| 4.0–6.9 | Medium |
| 7.0–8.9 | High |
| 9.0–10.0 | Critical |

```mermaid
flowchart LR
  N0["0.0 = no score<br/>(vector or missing)"] --> N1["0.1"]
  N1 -->|"Low"| N4["4.0"]
  N4 -->|"Medium"| N7["7.0"]
  N7 -->|"High"| N9["9.0"]
  N9 -->|"Critical"| N10["10.0"]
  N0 -.-> X["not rankable"]
```

The bands are contiguous — `0.1` is the first rankable score, so a real `0.0` never lands in *Low*; it means "no numeric score" (a vector string that `GetScore()` couldn't parse, or a missing field). That is why the table starts at `0.1` while the SDK getter returns `0.0` for the unrankable case.

## Decision tree

```mermaid
flowchart TD
  Q["Assessing risk?"] --> Ask{"Which CVSS?"}
  Ask -->|"v3"| V3["GetCVSS3() / query --severity cvss3"]
  Ask -->|"v2"| V2["GetCVSS2() / query --severity cvss2"]
  V3 --> S["GetScore()"]
  V2 --> S
  S --> Band["Map to Low/Medium/High/Critical"]
```

## Parsing path: vector vs number

```mermaid
flowchart TD
  SRC["OSV score field"] --> T{"Number or vector?"}
  T -->|"number e.g. 7.5"| NUM["GetScore() returns 7.5"]
  T -->|"vector e.g. CVSS:3.1/AV:N/..."| VEC["GetScore() returns 0.0<br/>parse vector yourself"]
  VEC --> GSF["GetScoreAsFloat()<br/>returns error hint"]
  NUM --> BAND["→ High"]
```

## Top-level vs per-affected severity

```mermaid
flowchart TD
  TOP["Top-level severity<br/>v.Severity (SeveritySlice)"] --> G3["GetCVSS3() record-level CVSS"]
  AFF["affected[].severity<br/>([]*Severity, optional, per-affected)"] --> P3["per-affected-range CVSS"]
```

`affected[].severity` is an optional severity slice scoped to a single affected entry (type `[]*Severity` — note this is a bare slice, *not* `SeveritySlice`, so it has no `GetCVSS3()` helper; iterate it directly). It is separate from the top-level `severity`.

## Anatomy of a CVSS vector

When `score` is a vector string rather than a number, this is what those slash-separated tokens mean — the reason `GetScore()` can't just `ParseFloat` it.

```mermaid
flowchart LR
  V["CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H"] --> H["CVSS:3.1<br/>version prefix"]
  V --> AV["AV:N<br/>Attack Vector = Network"]
  V --> AC["AC:L<br/>Attack Complexity = Low"]
  V --> PR["PR:N<br/>Privileges Required = None"]
  V --> UI["UI:N<br/>User Interaction = None"]
  V --> CIA["C:H / I:H / A:H<br/>Confidentiality / Integrity / Availability"]
```

::: tip Vector → number needs a CVSS calculator
The numeric 0–10 score is *derived* from these metrics by the CVSS formula, not stored in the string. That is why the SDK hands you the vector verbatim and leaves scoring to a dedicated CVSS library — the OSV record itself only guarantees the vector.
:::

## Notes

- OSV `score` may be a CVSS vector string (`CVSS:3.1/AV:N/...`) rather than a number — in that case `GetScore()` returns `0.0`. Parse the vector yourself if you need the numeric score from a vector.
- `SeverityTypeCVSS2 = "CVSS_V2"`, `SeverityTypeCVSS3 = "CVSS_V3"`

## Cross-references

- [[osv-query]] — the `--severity` flag lives here
- [[osv-affected]] — per-affected severity (`affected[].severity`)
- [Methods](/reference/methods#severity) — full severity API
