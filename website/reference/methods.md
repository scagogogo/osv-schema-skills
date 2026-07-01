# Methods

Quick reference for the SDK's most-used methods. All verified against source.

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
| `GetScoreAsFloat` | `() float64` | Alias of `GetScore` |
| `GetScoreAsPointer` | `() *float64` | Score as pointer (for nullable fields) |

## References

| Method | Signature | Description |
|--------|-----------|-------------|
| `FilterByType` | `(ReferenceType) References` | Filter by `ADVISORY`, `FIX`, etc. |

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

## Serialization helpers

Most types implement `sql.Scanner` and `driver.Valuer`, so they store cleanly as JSON columns under GORM. The complex nested types (`AffectedSlice`, `SeveritySlice`, `Package`, `Credits`) marshal themselves to/from JSON automatically.

Source: root package [`*.go`](https://github.com/scagogogo/osv-schema-skills)
