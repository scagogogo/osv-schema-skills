# Go SDK

The Go SDK is the type-safe foundation under both the CLI and the Skills. Use it when embedding OSV parsing/filtering/querying into a Go application.

## Install

```bash
go get -u github.com/scagogogo/osv-schema-skills
```

```go
import osv "github.com/scagogogo/osv-schema-skills"
```

## Quick start

```go
package main

import (
    "fmt"
    "log"

    osv "github.com/scagogogo/osv-schema-skills"
)

func main() {
    // Parse OSV data from a JSON file
    v, err := osv.UnmarshalFromJsonFile[any, any]("vulnerability.json")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s\n", v.ID)
    fmt.Printf("Summary: %s\n", v.Summary)

    // Get CVE from aliases
    if cve := v.Aliases.GetCVE(); cve != "" {
        fmt.Printf("CVE: %s\n", cve)
    }

    // Check if a specific ecosystem is affected
    if v.Affected.HasEcosystem("npm") {
        fmt.Println("Affects npm packages")
    }

    // Get CVSS v3 score
    if cvss3 := v.Severity.GetCVSS3(); cvss3 != nil {
        fmt.Printf("CVSS v3: %.1f\n", cvss3.GetScore())
    }
}
```

## Core type

```go
type OsvSchema[EcosystemSpecific, DatabaseSpecific any] struct {
    SchemaVersion    string
    ID               string
    Modified         time.Time
    Published        time.Time
    Withdrawn        string // string, not time.Time — check non-empty for withdrawn
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

Generic type parameters `EcosystemSpecific` and `DatabaseSpecific` let you attach custom data per ecosystem or vulnerability database. Use `any` for general-purpose parsing.

```mermaid
graph TD
  subgraph "Generic extension points"
    E["EcosystemSpecific<br/>per-ecosystem data"]
    D["DatabaseSpecific<br/>per-DB data"]
  end
  OSV["OsvSchema&lt;Eco, DB&gt;"] --> E
  OSV --> D
```

## Key methods

See the full table in [Reference → Methods](/reference/methods). Highlights:

| Type | Method | Description |
|------|--------|-------------|
| `OsvSchema` | `Affected.HasEcosystem(eco)` | Check if ecosystem is affected |
| `AffectedSlice` | `FilterByEcosystem(eco)` | Filter affected packages |
| `Aliases` | `GetCVE()` | Get first CVE identifier |
| `SeveritySlice` | `GetCVSS3()` / `GetCVSS2()` | Get CVSS severity entry |
| `Severity` | `GetScore()` | Parse score as float64 |
| `References` | `FilterByType(t)` | Filter by reference type |
| `Package` | `IsMaven()` / `GetGroupID()` / `GetArtifactID()` | Maven decomposition |

## Serialization

Every core type carries `json`, `yaml`, `mapstructure`, `db`, `bson`, `gorm` tags — JSON, YAML, mapstructure, GORM, and MongoDB (BSON) work out of the box.

## Design notes

- **Never nil from constructors** — `UnmarshalFromJsonFile` / `UnmarshalFromJson` return errors explicitly; the result is never a nil pointer on success.
- **Withdrawn is a string** — not `time.Time`. Check for a non-empty string to determine withdrawal status.
- **Database strategy** — simple fields are columns; complex nested structures (`AffectedSlice`, `SeveritySlice`) are stored as JSON strings via the GORM serializer.

## Requirements

- Go 1.18+
