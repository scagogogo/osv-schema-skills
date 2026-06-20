---
name: osv-filter
description: Filter OSV vulnerability data by ecosystem, reference type, or alias pattern. Triggers on mentions of filtering vulnerabilities by package ecosystem (npm, PyPI, Maven), reference type (ADVISORY, FIX), or alias pattern (CVE, GHSA).
allowed-tools: "Bash(osv:*)"
argument-hint: <osv-json-file>
---

# OSV Filter

> **Setup:** See `/osv-installation` for one-time CLI/SDK install.
> **Layers:** SDK (Go) → CLI (shell) — pick your entry point.

## When to Use

- Filter affected packages by ecosystem (e.g., show only npm or PyPI packages)
- Filter references by type (e.g., show only ADVISORY or FIX links)
- Filter aliases by pattern (e.g., show only CVE identifiers)
- Check if a vulnerability affects a specific ecosystem

## Decision Tree

```
OSV data → what to filter?
├─ By package ecosystem?        → --ecosystem PyPI / FilterByEcosystem()
├─ By reference type?           → --ref-type ADVISORY / FilterByType()
├─ By alias pattern?            → --alias CVE / Aliases.Filter()
└─ Check ecosystem presence?    → HasEcosystem()
```

## Task Patterns

### Filter by ecosystem

**Goal:** Show only PyPI-affected packages from a vulnerability.

| Layer | Approach |
|-------|----------|
| CLI | `osv filter -e PyPI vulnerability.json` |
| SDK | `osvData.Affected.FilterByEcosystem(osv.EcosystemPyPI)` |

### Check if ecosystem is affected

**Goal:** Does this vulnerability affect npm packages?

| Layer | Approach |
|-------|----------|
| CLI | `osv filter -e npm vulnerability.json` (check `Has Ecosystem` field) |
| SDK | `osvData.Affected.HasEcosystem(osv.EcosystemNPM)` |

### Filter references by type

**Goal:** Show only advisory and fix references.

| Layer | Approach |
|-------|----------|
| CLI | `osv filter -r ADVISORY vulnerability.json` |
| SDK | `osvData.References.FilterByType(osv.ReferenceTypeAdvisory)` |

### Filter aliases by pattern

**Goal:** Show only CVE identifiers from aliases.

| Layer | Approach |
|-------|----------|
| CLI | `osv filter -a CVE vulnerability.json` |
| SDK | `osvData.Aliases.Filter(func(alias string) bool { return strings.HasPrefix(alias, "CVE-") })` |

### Combine filters

**Goal:** Filter by both ecosystem and reference type.

```bash
osv filter -e PyPI -r FIX vulnerability.json
```

## API Reference

### SDK — AffectedSlice Methods

```go
// Check if any affected package belongs to the ecosystem
func (a AffectedSlice[E, D]) HasEcosystem(eco Ecosystem) bool

// Filter affected packages by ecosystem
func (a AffectedSlice[E, D]) FilterByEcosystem(eco Ecosystem) AffectedSlice[E, D]

// Custom filter with predicate function
func (a AffectedSlice[E, D]) Filter(fn func(*Affected[E, D]) bool) AffectedSlice[E, D]
```

### SDK — References Methods

```go
// Filter references by type
func (r References) FilterByType(refType ReferenceType) References
```

### SDK — Aliases Methods

```go
// Filter aliases with custom predicate
func (a Aliases) Filter(fn func(string) bool) Aliases

// Quick access to CVE identifier
func (a Aliases) GetCVE() string
```

### SDK — Ecosystem Constants

```go
EcosystemNPM       = "npm"
EcosystemPyPI      = "PyPI"
EcosystemMaven     = "Maven"
EcosystemNuGet     = "NuGet"
EcosystemRubyGems  = "RubyGems"
EcosystemGo        = "Go"
EcosystemCargo     = "Cargo"
EcosystemPub       = "Pub"
EcosystemHex       = "Hex"
EcosystemPackagist = "Packagist"
// ... and more
```

### SDK — ReferenceType Constants

```go
ReferenceTypeAdvisory   ReferenceType = "ADVISORY"
ReferenceTypeArticle    ReferenceType = "ARTICLE"
ReferenceTypeFix        ReferenceType = "FIX"
ReferenceTypePackage    ReferenceType = "PACKAGE"
ReferenceTypeReport     ReferenceType = "REPORT"
ReferenceTypeWeb        ReferenceType = "WEB"
// ... and more
```

### CLI Commands

```bash
osv filter -e <ecosystem> <file>       # Filter by ecosystem
osv filter -r <ref-type> <file>        # Filter by reference type
osv filter -a <alias-pattern> <file>   # Filter by alias pattern
osv filter -e PyPI -r FIX <file>       # Combine filters
osv filter -o json -e PyPI <file>      # JSON output
```

## Cross-References

- [[osv-parse]] — parse OSV JSON files first
- [[osv-query]] — extract specific sub-information
- [[osv-affected]] — detailed affected package analysis
- [[osv-severity]] — severity-based filtering

## Important Notes

- At least one filter flag is required (`--ecosystem`, `--ref-type`, or `--alias`)
- Ecosystem names are case-sensitive and must match OSV spec (e.g., `PyPI` not `pypi`, `Maven` not `maven`)
- Reference types are case-insensitive in CLI (auto-uppercased)
- `HasEcosystem` returns a boolean; `FilterByEcosystem` returns the filtered slice
- `GetCVE()` returns the first CVE alias found, or empty string if none
