---
name: osv-query
description: Extract specific sub-information from OSV vulnerability data — severity details (CVSS v2/v3), Maven package decomposition, version ranges, or event timelines. Triggers on queries about CVSS scores, Maven groupId/artifactId, version ranges, or when user needs focused extraction from OSV data.
allowed-tools: "Bash(osv:*)"
argument-hint: <osv-json-file>
---

# OSV Query

> **Setup:** See `/osv-installation` for one-time CLI/SDK install.
> **Layers:** SDK (Go) → CLI (shell) — pick your entry point.

## When to Use

- Extract CVSS severity scores (v2 or v3) from a vulnerability
- Decompose Maven package names into groupId and artifactId
- View version ranges for affected packages
- Inspect event timelines (introduced/fixed/last_affected/limit)

## Decision Tree

```
OSV data → what sub-info to extract?
├─ CVSS severity?        → --severity cvss3 / GetCVSS3()
├─ Maven decomposition?  → --maven / IsMaven() + GetGroupID()
├─ Version ranges?       → --ranges / Range.Type + Repo
└─ Event timeline?       → --events / IsIntroduced/IsFixed/etc.
```

## Task Patterns

### Query CVSS severity scores

**Goal:** Get CVSS v3 score from a vulnerability.

| Layer | Approach |
|-------|----------|
| CLI | `osv query --severity cvss3 vulnerability.json` |
| SDK | `osvData.Severity.GetCVSS3().GetScore()` |

### Query Maven package decomposition

**Goal:** Split Maven package name into groupId:artifactId.

| Layer | Approach |
|-------|----------|
| CLI | `osv query --maven vulnerability.json` |
| SDK | `pkg.IsMaven()`, `pkg.GetGroupID()`, `pkg.GetArtifactID()` |

### Query version ranges

**Goal:** View version ranges for all affected packages.

| Layer | Approach |
|-------|----------|
| CLI | `osv query --ranges vulnerability.json` |
| SDK | `affected.Ranges` — iterate `Range[DatabaseSpecific]` |

### Query event timeline

**Goal:** Show introduced/fixed/last_affected/limit events.

| Layer | Approach |
|-------|----------|
| CLI | `osv query --events vulnerability.json` |
| SDK | `event.IsIntroduced()`, `event.IsFixed()`, etc. |

## API Reference

### SDK — SeveritySlice Methods

```go
// Get CVSS v3 severity entry (nil if not present)
func (s SeveritySlice) GetCVSS3() *Severity

// Get CVSS v2 severity entry (nil if not present)
func (s SeveritySlice) GetCVSS2() *Severity
```

### SDK — Severity Methods

```go
// Get score as float64 (0.0 if Score string is empty or unparseable)
func (s *Severity) GetScore() float64

// Get score as float64 with error
func (s *Severity) GetScoreAsFloat() (float64, error)

// Get score as *float64 pointer
func (s *Severity) GetScoreAsPointer() *float64
```

### SDK — Package Methods

```go
// Check if package belongs to Maven ecosystem
func (p *Package) IsMaven() bool

// Extract groupId from Maven package name (before last colon)
func (p *Package) GetGroupID() string

// Extract artifactId from Maven package name (after last colon)
func (p *Package) GetArtifactID() string
```

### SDK — Event Methods

```go
func (e *Event) IsIntroduced() bool    // version was first affected
func (e *Event) IsFixed() bool         // version contains the fix
func (e *Event) IsLastAffected() bool  // last affected version
func (e *Event) IsLimit() bool         // upper limit (not a known version)
```

### SDK — Range Type Constants

```go
RangeTypeSemver    RangeType = "SEMVER"
RangeTypeEcosystem RangeType = "ECOSYSTEM"
RangeTypeGit       RangeType = "GIT"
```

### CLI Commands

```bash
osv query --severity cvss3 <file>   # CVSS v3 score
osv query --severity cvss2 <file>   # CVSS v2 score
osv query --maven <file>            # Maven package decomposition
osv query --ranges <file>           # Version ranges
osv query --events <file>           # Event timeline (introduced/fixed/etc.)
osv query --ranges --events <file>  # Combined ranges + events
osv query -o json --severity cvss3 <file>  # JSON output
```

## Cross-References

- [[osv-parse]] — full parsing and display of OSV data
- [[osv-filter]] — filter by ecosystem before querying
- [[osv-severity]] — detailed severity analysis
- [[osv-affected]] — detailed affected package analysis

## Important Notes

- At least one query flag is required (`--severity`, `--maven`, `--ranges`, or `--events`)
- `GetCVSS3()` and `GetCVSS2()` return `nil` if no matching severity type exists
- `GetScore()` returns 0.0 for unparseable scores; use `GetScoreAsFloat()` for error handling
- Maven decomposition only works for packages in the `Maven` ecosystem — check `IsMaven()` first
- Event fields are mutually exclusive: only one of `Introduced`, `Fixed`, `LastAffected`, `Limit` is non-empty per event
- Range type `ECOSYSTEM` is the most common in OSV data; `SEMVER` and `GIT` are less frequent
