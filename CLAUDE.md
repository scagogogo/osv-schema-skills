# CLAUDE.md

This file provides guidance to Claude Code when working in this repository.

## Project Overview

**osv-schema-skills** is a Go library + CLI + Claude Code Skills bundle for the OSV (Open Source Vulnerability) Schema. It provides type-safe parsing, validation, filtering, and querying of vulnerability data in the OSV format — accessible via 6 specialized Claude Code skills and a Go SDK.

## Repository Structure

```
.
├── .claude/skills/       # 6 Claude Code Skills (SKILL.md files)
│   ├── osv-parse/        # Parse and display OSV JSON data
│   ├── osv-validate/     # Validate OSV JSON files
│   ├── osv-filter/       # Filter by ecosystem, reference type, alias
│   ├── osv-query/        # Extract severity, maven, ranges, events
│   ├── osv-severity/     # CVSS severity analysis
│   ├── osv-affected/     # Affected package and version analysis
│   └── osv-installation/ # Installation and setup guide
├── .claude-plugin/       # Claude Code plugin & marketplace manifests
├── cmd/osv/              # CLI binary entrypoint (cobra)
├── docs/superpowers/     # Historical development plans
├── test_data/            # OSV JSON test fixtures
├── *.go                  # Core library (root package osv_schema)
└── *_test.go             # Tests
```

## Build & Test

```bash
# Build everything
go build ./...

# Build CLI binary
go build -o osv ./cmd/osv/

# Run all tests
go test ./...

# Run specific tests
go test -v -run TestUnmarshal ./...
```

## Code Style

- Standard Go conventions: `gofmt`, `go vet`
- Tests use `github.com/stretchr/testify` for assertions
- Test files are named `*_test.go` alongside the source files
- Comments are mixed Chinese/English — maintain consistency with existing style
- All core types use Go generics for `EcosystemSpecific` and `DatabaseSpecific` extensibility

## Key Design Decisions

- **Generic type parameters**: `OsvSchema[EcosystemSpecific, DatabaseSpecific]` allows flexible customization per ecosystem/database
- **Multiple serialization tags**: Each field has `json`, `yaml`, `mapstructure`, `db`, `bson`, `gorm` tags for broad compatibility
- **Database strategy**: Simple fields stored as columns; complex nested structures (AffectedSlice, SeveritySlice) stored as JSON strings via GORM serializer
- **Never nil from NewVersion-style constructors**: Unmarshal functions return errors explicitly
- **Withdrawn field is string**: Not `time.Time` — check for non-empty string to determine withdrawal status

## Skills

The 6 skills in `.claude/skills/*/SKILL.md` follow this structure:
- YAML frontmatter with `name`, `description`, `allowed-tools`, `argument-hint`
- Two access paths per skill: SDK (Go) and CLI
- Each skill has: When to Use, Decision Tree, Task Patterns, API Reference, Cross-References, Important Notes

When adding or modifying skills:
- Match the existing format and depth
- Cover both access paths (SDK and CLI)
- Include code examples for each
- Add an "Important Notes" section with gotchas
- Use `[[skill-name]]` for cross-references between skills
