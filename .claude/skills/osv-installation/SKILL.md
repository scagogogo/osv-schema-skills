---
name: osv-installation
description: Install and set up the OSV Schema Skills CLI and Go SDK. Triggers on first use, setup requests, installation instructions, or when user needs to install the osv CLI tool or Go library.
allowed-tools: "Bash(go:*)"
---

# OSV Schema Skills — Installation

> Install the Go SDK, CLI tool, and Claude Code Skills integration.

## Installation Options

### Option 1: Go Install (Recommended)

```bash
# Install CLI tool
go install github.com/scagogogo/osv-schema-skills/cmd/osv@latest

# Verify installation
osv version
```

### Option 2: Clone & Build from Source

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills

# Build CLI
go build -o osv ./cmd/osv/

# Verify
./osv version
```

### Option 3: Use as Go Library

```bash
go get -u github.com/scagogogo/osv-schema-skills
```

```go
import osv "github.com/scagogogo/osv-schema-skills"
```

## Go SDK Quick Start

```go
package main

import (
    "fmt"
    "log"

    osv "github.com/scagogogo/osv-schema-skills"
)

func main() {
    // Parse OSV data from JSON file
    vulnerability, err := osv.UnmarshalFromJsonFile[any, any]("vulnerability.json")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Vulnerability ID: %s\n", vulnerability.ID)
    fmt.Printf("Summary: %s\n", vulnerability.Summary)

    // Get CVE from aliases
    if cve := vulnerability.Aliases.GetCVE(); cve != "" {
        fmt.Printf("CVE: %s\n", cve)
    }

    // Check if specific ecosystem is affected
    if vulnerability.Affected.HasEcosystem("npm") {
        fmt.Println("This vulnerability affects npm packages")
    }
}
```

## CLI Quick Start

```bash
# Parse an OSV JSON file
osv parse vulnerability.json

# Parse with full details
osv parse -v vulnerability.json

# Validate an OSV JSON file
osv validate vulnerability.json

# Filter by ecosystem
osv filter -e PyPI vulnerability.json

# Query severity
osv query --severity cvss3 vulnerability.json
```

## Claude Code Skills Integration

This repository includes 7 Claude Code Skills in `.claude/skills/`:

| Skill | Purpose |
|-------|---------|
| `osv-parse` | Parse and display OSV JSON data |
| `osv-validate` | Validate OSV JSON files |
| `osv-filter` | Filter by ecosystem, reference type, alias |
| `osv-query` | Extract severity, maven, ranges, events |
| `osv-severity` | CVSS severity analysis |
| `osv-affected` | Affected package and version analysis |
| `osv-installation` | Setup & installation guide (this skill) |

Skills are automatically available when Claude Code opens this repository.

## Requirements

- **Go 1.18+** (for SDK and CLI)
- **Internet access** only needed for `go get`/`go install`

## Important Notes

- The CLI binary is named `osv` — ensure it's on your `$PATH`
- Use `any` for generic type parameters when you don't need ecosystem/database-specific data
- The library supports JSON, YAML, and database serialization (GORM, MongoDB)
