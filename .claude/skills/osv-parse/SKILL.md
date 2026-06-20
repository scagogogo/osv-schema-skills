---
name: osv-parse
description: Parse an OSV (Open Source Vulnerability) JSON file and display structured vulnerability data. Triggers on mentions of OSV parsing, vulnerability JSON reading, CVE/GHSA data extraction, or when user provides an OSV JSON file path.
allowed-tools: "Bash(osv:*)"
argument-hint: <osv-json-file>
---

# OSV Parse

> **Setup:** See `/osv-installation` for one-time CLI/SDK install.
> **Layers:** SDK (Go) → CLI (shell) — pick your entry point.

## When to Use

- Parse an OSV JSON file to extract vulnerability ID, summary, severity, affected packages
- Read and display vulnerability data from OSV-format JSON files
- Extract CVE/GHSA identifiers from vulnerability records
- Inspect the full structure of an OSV vulnerability entry

## Decision Tree

```
OSV JSON file → what do you need?
├─ Quick overview?              → osv parse <file>
├─ All fields including dates?  → osv parse -v <file>
├─ Machine-readable output?     → osv parse -o json <file>
└─ Programmatic access?         → SDK: UnmarshalFromJsonFile()
```

## Task Patterns

### Parse and display vulnerability overview

**Goal:** Read `vulnerability.json` and show ID, summary, severity, affected packages.

| Layer | Approach |
|-------|----------|
| SDK | `osv.UnmarshalFromJsonFile[any, any]("vulnerability.json")` |
| CLI | `osv parse vulnerability.json` |

### Parse with full details (verbose)

**Goal:** Show all fields including dates, details, ranges, credits, related IDs.

| Layer | Approach |
|-------|----------|
| CLI | `osv parse -v vulnerability.json` |

### Parse with JSON output

**Goal:** Get machine-readable JSON output for piping to other tools.

| Layer | Approach |
|-------|----------|
| CLI | `osv parse -o json vulnerability.json` |

## API Reference

### SDK — Unmarshal Functions

```go
// Parse from JSON byte slice
osv.UnmarshalFromJson[EcosystemSpecific, DatabaseSpecific](data []byte) (*OsvSchema[E, D], error)

// Parse from JSON file path
osv.UnmarshalFromJsonFile[EcosystemSpecific, DatabaseSpecific](filePath string) (*OsvSchema[E, D], error)
```

### SDK — OsvSchema Struct

```go
type OsvSchema[EcosystemSpecific, DatabaseSpecific any] struct {
    SchemaVersion    string
    ID               string
    Modified         time.Time
    Published        time.Time
    Withdrawn        string
    Aliases          Aliases
    Related          Related
    Summary          string
    Details          string
    Severity         SeveritySlice
    Affected         AffectedSlice[EcosystemSpecific, DatabaseSpecific]
    References       References
    DatabaseSpecific DatabaseSpecific
    Credits          *Credits
}
```

### CLI Commands

```bash
osv parse <file>              # Display key fields
osv parse -v <file>           # Verbose: all fields including dates, details, ranges
osv parse -o json <file>      # JSON output format
osv parse -o json -v <file>   # Full JSON with all fields
```

## Cross-References

- [[osv-validate]] — validate OSV JSON files against the schema
- [[osv-filter]] — filter parsed data by ecosystem, reference type, or alias
- [[osv-query]] — extract specific sub-information (severity, maven, ranges)
- [[osv-severity]] — detailed CVSS severity analysis
- [[osv-affected]] — affected package analysis

## Important Notes

- The generic type parameters `EcosystemSpecific` and `DatabaseSpecific` default to `any` for general use
- `Withdrawn` is a string field (not `time.Time`) — check for non-empty string to determine if a vulnerability is withdrawn
- `Aliases` is a custom slice type with helper methods like `GetCVE()`
- Use `-o json` for machine-readable output when piping to other tools
