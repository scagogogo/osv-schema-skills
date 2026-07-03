# osv-affected

Analyze affected packages and version ranges.

> **Trigger:** mentions of affected packages, version ranges, impacted ecosystems, or determining which packages/versions are affected.
> **Skill source:** [`.claude/skills/osv-affected/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-affected/SKILL.md)

## CLI

```bash
osv parse -v vulnerability.json             # Full affected details + ranges
osv filter -e PyPI vulnerability.json       # Narrow to one ecosystem
osv query --ranges vulnerability.json       # Version ranges
osv query --events vulnerability.json       # Event timeline
```

## SDK

```go
// Presence
v.Affected.HasEcosystem(osv.EcosystemPyPI)

// Filter
pypi := v.Affected.FilterByEcosystem(osv.EcosystemPyPI)

// Iterate ranges & events
for _, a := range v.Affected {
    if a.Package == nil {
        continue // a missing package is rare but possible on untrusted data
    }
    fmt.Println(a.Package.Ecosystem, a.Package.Name)
    for _, r := range a.Ranges {
        fmt.Println("  range type:", r.Type)   // SEMVER / ECOSYSTEM / GIT
        for _, e := range r.Events {
            // e.IsIntroduced() / IsFixed() / IsLastAffected() / IsLimit()
        }
    }
}
```

## Structure

```mermaid
graph TD
  AFF["Affected[]"] --> PKG["package<br/>ecosystem · name · purl"]
  AFF --> VER["versions[]"]
  AFF --> RNG["ranges[]"]
  AFF --> ASEV["severity[] (per-affected)"]
  RNG --> TYPE["type: SEMVER/ECOSYSTEM/GIT"]
  RNG --> EVT["events[]"]
  EVT --> I["introduced"]
  EVT --> F["fixed"]
  EVT --> L["last_affected"]
  EVT --> LM["limit"]
```

## Affected data model

```mermaid
classDiagram
  class Affected {
    +Package *Package
    +Versions []string
    +Ranges []*Range
    +Severity []*Severity
    +EcosystemSpecific Eco
    +DatabaseSpecific DB
  }
  class Package {
    +Ecosystem Ecosystem
    +Name string
    +PackageUrl string
  }
  class Range {
    +Type RangeType
    +Repo string
    +Events []Event
    +DatabaseSpecific DB
  }
  class Event {
    +Introduced string
    +Fixed string
    +LastAffected string
    +Limit string
  }
  Affected --> Package
  Affected --> Range
  Range --> Event
```

## Decision tree

```mermaid
flowchart TD
  Q["What about affected?"] --> L["List all packages → parse -v"]
  Q --> E["Check ecosystem → HasEcosystem / filter -e"]
  Q --> R["Version ranges → query --ranges"]
  Q --> M["Maven GAV → query --maven"]
  Q --> EV["Event timeline → query --events"]
```

## Range type comparison

```mermaid
flowchart TD
  T["range.type"] --> SE["SEMVER<br/>semantic version range"]
  T --> EC["ECOSYSTEM<br/>most common, ecosystem versions"]
  T --> GI["GIT<br/>git commit range"]
```

- `RangeTypeEcosystem` (`ECOSYSTEM`) is the most common; `SEMVER` and `GIT` are less frequent.

## Is my version affected? — a worked example

The `versions[]` list is an explicit enumeration, but real records lean on `ranges[]`. To answer "is `X` affected" from a range, resolve its events. Example: `introduced: 1.2.0`, `fixed: 1.4.1`.

```mermaid
flowchart TD
  V["candidate version X"] --> C1{"X ≥ 1.2.0 ?"}
  C1 -->|no| SAFE1["not yet introduced → SAFE"]
  C1 -->|yes| C2{"X ≥ 1.4.1 ?"}
  C2 -->|yes| SAFE2["at/after fix → SAFE"]
  C2 -->|no| HIT["introduced, not fixed → AFFECTED"]
```

::: warning `versions[]` and `ranges[]` can disagree in shape
Some records list exact affected `versions[]`; others give only `ranges[]`; many give both. Prefer `ranges[]` for open-ended "everything since 1.2.0" cases, and treat `versions[]` as the authoritative enumeration when present. Never assume one implies the other.
:::

## Notes

- `RangeTypeEcosystem` (`ECOSYSTEM`) is the most common; `SEMVER` and `GIT` are less frequent
- Event fields are mutually exclusive per event object
- `affected[].severity` is optional per-affected severity, separate from top-level `severity`

## Cross-references

- [[osv-filter]] — narrow affected by ecosystem
- [[osv-query]] — extract ranges/events/maven
- [OSV Schema](/reference/osv-schema) — full type model
