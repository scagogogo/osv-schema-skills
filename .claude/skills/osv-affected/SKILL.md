---
name: osv-affected
description: Analyze affected packages and version ranges in OSV vulnerability records. Triggers on mentions of affected packages, version ranges, impacted ecosystems, package vulnerability status, or when user needs to determine which packages/versions are affected by a vulnerability.
allowed-tools: "Bash(osv:*)"
argument-hint: <osv-json-file>
---

# OSV Affected Package Analysis

> **Setup:** See `/osv-installation` for one-time CLI/SDK install.
> **Layers:** SDK (Go) → CLI (shell) — pick your entry point.

## When to Use

- Determine which packages and versions are affected by a vulnerability
- Check if a specific ecosystem is affected (npm, PyPI, Maven, etc.)
- Analyze version ranges and the introduced→fixed event timeline
- Extract package URLs (purl) and Maven decomposition
- Get per-package severity ratings

## Decision Tree

```
Affected data → what do you need?
├─ List all affected packages?  → Affected slice iteration / osv parse -v
├─ Check specific ecosystem?    → HasEcosystem() / osv filter -e
├─ Filter by ecosystem?         → FilterByEcosystem() / osv filter -e
├─ Version ranges?              → Range.Type + Events / osv query --ranges
├─ Maven decomposition?         → IsMaven() / GetGroupID() / osv query --maven
└─ Per-package severity?        → Affected[i].Severity
```

## Task Patterns

### List all affected packages

**Goal:** Show all packages affected by a vulnerability.

| Layer | Approach |
|-------|----------|
| CLI | `osv parse vulnerability.json` |
| SDK | Iterate `osvData.Affected` slice |

```go
for _, affected := range osvData.Affected {
    if affected.Package != nil {
        fmt.Printf("%s/%s\n", affected.Package.Ecosystem, affected.Package.Name)
    }
}
```

### Check if a specific version is affected

**Goal:** Is version `1.2.3` of package `foo` affected?

```go
for _, affected := range osvData.Affected {
    if affected.Package != nil && affected.Package.Name == "foo" {
        for _, v := range affected.Versions {
            if v == "1.2.3" {
                fmt.Println("Version 1.2.3 IS affected")
            }
        }
    }
}
```

### Analyze version ranges and events

**Goal:** Understand the version range (introduced → fixed) for an affected package.

| Layer | Approach |
|-------|----------|
| CLI | `osv query --events vulnerability.json` |
| SDK | Iterate `affected.Ranges[i].Events` |

```go
for _, affected := range osvData.Affected {
    for _, r := range affected.Ranges {
        fmt.Printf("Range type: %s\n", r.Type)
        for _, event := range r.Events {
            switch {
            case event.IsIntroduced():
                fmt.Printf("  Introduced: %s\n", event.Introduced)
            case event.IsFixed():
                fmt.Printf("  Fixed: %s\n", event.Fixed)
            case event.IsLastAffected():
                fmt.Printf("  Last affected: %s\n", event.LastAffected)
            case event.IsLimit():
                fmt.Printf("  Limit: %s\n", event.Limit)
            }
        }
    }
}
```

### Maven package decomposition

**Goal:** Split `org.apache.commons:commons-lang3` into groupId and artifactId.

| Layer | Approach |
|-------|----------|
| CLI | `osv query --maven vulnerability.json` |
| SDK | `pkg.IsMaven()`, `pkg.GetGroupID()`, `pkg.GetArtifactID()` |

```go
for _, affected := range osvData.Affected {
    if affected.Package != nil && affected.Package.IsMaven() {
        fmt.Printf("GroupID: %s\n", affected.Package.GetGroupID())
        fmt.Printf("ArtifactID: %s\n", affected.Package.GetArtifactID())
    }
}
```

## API Reference

### SDK — AffectedSlice Methods

```go
func (a AffectedSlice[E, D]) HasEcosystem(eco Ecosystem) bool
func (a AffectedSlice[E, D]) FilterByEcosystem(eco Ecosystem) AffectedSlice[E, D]
func (a AffectedSlice[E, D]) Filter(fn func(*Affected[E, D]) bool) AffectedSlice[E, D]
```

### SDK — Affected Struct

```go
type Affected[EcosystemSpecific, DatabaseSpecific any] struct {
    Package          *Package                           `json:"package"`
    Severity         SeveritySlice                      `json:"severity"`
    Versions         []string                           `json:"versions"`
    Ranges           []Range[DatabaseSpecific]          `json:"ranges"`
    EcosystemSpecific EcosystemSpecific                 `json:"ecosystem_specific"`
    DatabaseSpecific  DatabaseSpecific                  `json:"database_specific"`
}
```

### SDK — Package Struct

```go
type Package struct {
    Ecosystem  Ecosystem `json:"ecosystem"`
    Name       string    `json:"name"`
    PackageUrl string    `json:"purl"`
}
```

### SDK — Range Struct

```go
type Range[DatabaseSpecific any] struct {
    Type             RangeType             `json:"type"`
    Repo             string                `json:"repo"`
    Events           Events                `json:"events"`
    DatabaseSpecific DatabaseSpecific      `json:"database_specific"`
}
```

### CLI Commands

```bash
osv parse -v <file>                    # Full affected package details
osv filter -e <ecosystem> <file>       # Filter by ecosystem
osv query --ranges <file>              # Version ranges
osv query --events <file>              # Event timeline
osv query --maven <file>               # Maven decomposition
```

## Cross-References

- [[osv-parse]] — full parsing of OSV data
- [[osv-filter]] — filter affected packages by ecosystem
- [[osv-query]] — extract ranges, events, maven info
- [[osv-severity]] — severity analysis including per-package

## Important Notes

- `Affected.Package` can be `nil` — always check before accessing fields
- `Versions` is a flat list of explicit version strings; `Ranges` provide the introduced→fixed timeline
- A package can have multiple `Range` entries (e.g., one SEMVER range and one ECOSYSTEM range)
- Maven packages use `groupId:artifactId` as the `Name` field — use `GetGroupID()`/`GetArtifactID()` to decompose
- `EcosystemSpecific` and `DatabaseSpecific` are generic type parameters that vary by source database
- `purl` (Package URL) follows the [purl spec](https://github.com/package-url/purl-spec)
