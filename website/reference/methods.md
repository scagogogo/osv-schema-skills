# Methods

Quick reference for the SDK's most-used methods. All verified against source.

## Methods at a glance

Grouped by receiver type — this is the whole surface you will use day to day.

```mermaid
mindmap
  root((osv SDK))
    Aliases
      GetCVE
      Filter
    AffectedSlice
      HasEcosystem
      FilterByEcosystem
      Filter
    Package
      IsMaven
      GetGroupID
      GetArtifactID
    SeveritySlice
      GetCVSS3
      GetCVSS2
    Severity
      GetScore
      GetScoreAsFloat
      GetScoreAsPointer
    References
      FilterByType
    Event
      IsIntroduced
      IsFixed
      IsLastAffected
      IsLimit
    package-level
      UnmarshalFromJson
      UnmarshalFromJsonFile
```

## Aliases

| Method | Signature | Description |
|--------|-----------|-------------|
| `GetCVE` | `() string` | First identifier starting with `CVE-` |
| `Filter` | `(func(string) bool) Aliases` | Filter aliases by predicate |

## AffectedSlice

| Method | Signature | Description |
|--------|-----------|-------------|
| `HasEcosystem` | `(Ecosystem) bool` | Whether any affected entry matches ecosystem |
| `FilterByEcosystem` | `(Ecosystem) AffectedSlice` | Narrow to one ecosystem |
| `Filter` | `(func(*Affected) bool) AffectedSlice` | Custom predicate filter |

## Package

| Method | Signature | Description |
|--------|-----------|-------------|
| `IsMaven` | `() bool` | `Ecosystem == Maven` |
| `GetGroupID` | `() string` | Maven `groupId` (left of `:`) |
| `GetArtifactID` | `() string` | Maven `artifactId` (right of `:`) |

## SeveritySlice

| Method | Signature | Description |
|--------|-----------|-------------|
| `GetCVSS3` | `() *Severity` | CVSS v3 entry, or `nil` |
| `GetCVSS2` | `() *Severity` | CVSS v2 entry, or `nil` |

## Severity

| Method | Signature | Description |
|--------|-----------|-------------|
| `GetScore` | `() float64` | Parse the CVSS score as `float64` |
| `GetScoreAsFloat` | `() (float64, error)` | Parse score, returning an error if the vector string is malformed |
| `GetScoreAsPointer` | `() *float64` | Score as pointer (for nullable fields) |

All three share one parser (`GetScoreAsFloat`); the other two only differ in **how they report a parse failure** — and a `score` that holds a CVSS *vector string* (e.g. `CVSS:3.1/AV:N/…`) rather than a number *is* a parse failure. Pick the variant whose failure shape you can handle:

```mermaid
flowchart TD
  SCORE["Severity.score"] --> PARSE{"strconv.ParseFloat<br/>succeeds?"}
  PARSE -->|"yes (numeric string)"| NUM["the float value"]
  PARSE -->|"no (vector string / empty)"| FAIL["parse error"]
  NUM --> G["GetScore() → value"]
  NUM --> F["GetScoreAsFloat() → (value, nil)"]
  NUM --> P["GetScoreAsPointer() → &value"]
  FAIL --> G2["GetScore() → 0.0 ⚠️ silent"]
  FAIL --> F2["GetScoreAsFloat() → (0, err)"]
  FAIL --> P2["GetScoreAsPointer() → nil"]
```

::: warning `GetScore()` hides the vector-string case
Because `GetScore()` drops the error, a vector-string score is indistinguishable from a real `0.0`. When the distinction matters, use `GetScoreAsFloat()` (check `err`) or `GetScoreAsPointer()` (check `nil`) — and read the CVSS vector from `Severity.Score` directly.
:::

## References

| Method | Signature | Description |
|--------|-----------|-------------|
| `FilterByType` | `(...ReferenceType) References` | Filter by `ADVISORY`, `FIX`, etc. (accepts multiple) |

## Event

| Method | Signature | Description |
|--------|-----------|-------------|
| `IsIntroduced` | `() bool` | Event marks an introduced version |
| `IsFixed` | `() bool` | Event marks a fixed version |
| `IsLastAffected` | `() bool` | Event marks last affected version |
| `IsLimit` | `() bool` | Event marks a range limit |

## Parsing

| Function | Signature | Description |
|----------|-----------|-------------|
| `UnmarshalFromJson` | `([]byte) (*OsvSchema[Eco,DB], error)` | Parse from bytes |
| `UnmarshalFromJsonFile` | `(string) (*OsvSchema[Eco,DB], error)` | Parse from file path |

```go
// General-purpose parsing — use `any` for both generics
v, err := osv.UnmarshalFromJsonFile[any, any]("vuln.json")

// Or attach ecosystem/database-specific data
v, err := osv.UnmarshalFromJsonFile[MyEco, MyDB]("vuln.json")
```

## Method call graph

```mermaid
graph TD
  OSV["OsvSchema"] --> AL["v.Aliases"]
  OSV --> SEV["v.Severity"]
  OSV --> AFF["v.Affected"]
  OSV --> REF["v.References"]

  AL -->|"GetCVE()"| CVE["string"]
  AL -->|"Filter(f)"| AL2["Aliases"]

  SEV -->|"GetCVSS3()"| S3["*Severity"]
  S3 -->|"GetScore()"| SCORE["float64"]

  AFF -->|"HasEcosystem(e)"| BOOL["bool"]
  AFF -->|"FilterByEcosystem(e)"| AFF2["AffectedSlice"]
  AFF --> PKG["a.Package"]
  PKG -->|"IsMaven()"| MBOOL["bool"]
  PKG -->|"GetGroupID()"| GID["string"]

  REF -->|"FilterByType(t)"| REF2["References"]
```

## Parse & validate data flow

```mermaid
flowchart LR
  F["file/bytes"] --> U["UnmarshalFromJson[File]"]
  U --> V["*OsvSchema"]
  V --> CHK["check v.ID / v.SchemaVersion"]
  CHK --> OK{"non-empty?"}
  OK -->|"yes"| VALID["✓ valid"]
  OK -->|"no"| INVALID["✗ invalid"]
```

## Maven coordinate decomposition

`GetGroupID` / `GetArtifactID` split a Maven package name on the first `:`. They only make sense when `IsMaven()` is true.

```mermaid
flowchart LR
  N["package.Name<br/>'com.fasterxml.jackson.core:jackson-databind'"] --> CHK{"IsMaven()?"}
  CHK -->|no| SKIP["not a Maven package → skip"]
  CHK -->|yes| SPLIT["split on first ':'"]
  SPLIT --> G["GetGroupID()<br/>'com.fasterxml.jackson.core'"]
  SPLIT --> A["GetArtifactID()<br/>'jackson-databind'"]
```

## A real query, method by method

"Is `GHSA-…` a high-severity PyPI issue, and where's the fix?" — here's the exact method chain an agent (or your code) walks.

```mermaid
sequenceDiagram
  participant You as You / agent
  participant SDK as osv SDK
  You->>SDK: UnmarshalFromJsonFile[any,any](path)
  SDK-->>You: *OsvSchema
  You->>SDK: v.Affected.HasEcosystem(EcosystemPyPI)
  SDK-->>You: true
  You->>SDK: v.Severity.GetCVSS3()
  SDK-->>You: *Severity (vector string)
  You->>SDK: v.References.FilterByType(ReferenceTypeFix)
  SDK-->>You: References (fix links)
  Note over You,SDK: parse → filter → score → fix, all on one typed core
```

## Which method returns what

```mermaid
flowchart TD
  Q["What do you have?"] --> Q1["a slice → want one item"]
  Q1 --> M1["GetCVSS3 / GetCVSS2 / GetCVE<br/>→ single item or nil/empty"]
  Q --> Q2["a slice → want a subset"]
  Q2 --> M2["FilterByEcosystem / FilterByType / Filter<br/>→ a new slice"]
  Q --> Q3["a slice → yes/no question"]
  Q3 --> M3["HasEcosystem<br/>→ bool"]
  Q --> Q4["one item → a derived value"]
  Q4 --> M4["GetScore / GetGroupID / IsMaven / IsFixed<br/>→ scalar"]
```

## Serialization helpers

Most types implement `sql.Scanner` and `driver.Valuer`, so they store cleanly as JSON columns under GORM. The complex nested types (`AffectedSlice`, `SeveritySlice`, `Package`, `Credits`) marshal themselves to/from JSON automatically.

Source: root package [`*.go`](https://github.com/scagogogo/osv-schema-skills)
