# osv-query

Extract specific sub-information: CVSS severity, Maven decomposition, version ranges, event timelines.

> **Trigger:** queries about CVSS scores, Maven groupId/artifactId, version ranges, or focused extraction from OSV data.
> **Skill source:** [`.claude/skills/osv-query/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-query/SKILL.md)

## CLI

```bash
osv query --severity cvss3 vulnerability.json  # CVSS v3 entry + parsed score
osv query --severity cvss2 vulnerability.json  # CVSS v2
osv query --maven vulnerability.json           # Maven groupId/artifactId
osv query --ranges vulnerability.json          # Version ranges
osv query --events vulnerability.json          # Event timeline
osv query --ranges --events vulnerability.json # Combine
```

| Flag | Description |
|------|-------------|
| `--severity` | `cvss3` or `cvss2` |
| `--maven` | Decompose Maven `groupId:artifactId` |
| `--ranges` | Show version ranges |
| `--events` | Show event timeline |
| `-o, --output` | `text` (default) or `json` |

At least one flag is required.

`--severity -o json` returns the CVSS entry. The `score` is the raw vector string; the DTO also computes a `numeric_score` via `GetScore()`, but it carries `omitempty` — and since `GetScore()` returns `0.0` on a vector string, that field is dropped from the JSON:

```bash
osv query --severity cvss3 -o json test_data/GHSA-vxv8-r8q2-63xw.json
```

```json
{
  "id": "GHSA-vxv8-r8q2-63xw",
  "severity": {
    "type": "CVSS_V3",
    "score": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H"
  }
}
```

## The four extraction dimensions

```mermaid
flowchart TD
  DATA["OSV data"] --> SEV["--severity<br/>Severity → GetCVSS3/2"]
  DATA --> MAV["--maven<br/>Package → GetGroupID/ArtifactID"]
  DATA --> RNG["--ranges<br/>Affected[].Ranges"]
  DATA --> EVT["--events<br/>Range.Events timeline"]
  SEV --> OUT["extracted result"]
  MAV --> OUT
  RNG --> OUT
  EVT --> OUT
```

## SDK equivalent

```go
// Severity
if s := v.Severity.GetCVSS3(); s != nil { fmt.Println(s.GetScore()) }

// Maven
for _, a := range v.Affected {
    if a.Package.IsMaven() {
        fmt.Println(a.Package.GetGroupID(), a.Package.GetArtifactID())
    }
}

// Ranges & events
for _, a := range v.Affected {
    for _, r := range a.Ranges {
        for _, e := range r.Events {
            // e.IsIntroduced() / IsFixed() / IsLastAffected() / IsLimit()
        }
    }
}
```

## Decision tree

```mermaid
flowchart TD
  Q["What to extract?"] --> Sev["CVSS severity"]
  Q --> Mav["Maven GAV"]
  Q --> Rng["Version ranges"]
  Q --> Evt["Event timeline"]
  Sev --> S["osv query --severity cvss3|cvss2"]
  Mav --> M["osv query --maven"]
  Rng --> R["osv query --ranges"]
  Evt --> E["osv query --events"]
  R --> Comb{"Combine?"}
  E --> Comb
  Comb -->|"yes"| C["--ranges --events"]
```

## Version ranges vs events

```mermaid
graph TD
  AFF["Affected"] --> RNG["ranges[]"]
  RNG --> TYPE["type: SEMVER / ECOSYSTEM / GIT"]
  RNG --> EVT["events[]"]
  EVT --> I["introduced<br/>first affected version"]
  EVT --> F["fixed<br/>version with the fix"]
  EVT --> L["last_affected<br/>last affected version"]
  EVT --> LM["limit<br/>range upper bound"]
```

Event fields are mutually exclusive per event object — one of introduced/fixed/last_affected/limit each.

## A worked `--events` timeline

`--events` prints the raw ordered events. To turn them into a yes/no answer for a concrete version, walk them in order. Here is `introduced: 1.0.0` then `fixed: 1.5.0` resolved for three candidate versions:

```mermaid
flowchart LR
  subgraph events["range.events (in order)"]
    I["introduced 1.0.0"] --> F["fixed 1.5.0"]
  end
  events --> Q0["0.9.0 → before introduced → SAFE"]
  events --> Q1["1.2.0 → ≥1.0.0 and <1.5.0 → AFFECTED"]
  events --> Q2["1.5.0 → ≥ fixed → SAFE"]
```

::: tip The CLI gives you data, not a verdict
`osv query --events` deliberately stops at the raw timeline — it never decides "is version X affected", because that requires ecosystem-aware version comparison (see [RangeType](/reference/osv-schema#rangetype-—-how-versions-are-compared)). The walk above is what *you* implement on top of the per-event predicates.
:::

## Notes

- `GetCVSS3()` / `GetCVSS2()` return `nil` if the severity type is absent
- `GetScore()` returns `0.0` when the OSV `score` is a CVSS vector string rather than a number — use `GetScoreAsFloat()` for error handling
- Maven decomposition only applies to `Maven`-ecosystem packages
- Event fields are mutually exclusive: one of `introduced`/`fixed`/`last_affected`/`limit` per event

## Cross-references

- [[osv-parse]] — full parse first
- [[osv-severity]] — deeper severity analysis
- [[osv-affected]] — deeper affected/range analysis
