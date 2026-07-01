# osv-validate

Validate OSV JSON files against the schema.

> **Trigger:** mentions of OSV validation, vulnerability format checking, schema compliance, or verifying a file is well-formed.
> **Skill source:** [`.claude/skills/osv-validate/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-validate/SKILL.md)

## CLI

```bash
osv validate vulnerability.json              # Single file
osv validate file1.json file2.json           # Batch
osv validate -o json vulnerability.json      # JSON output
```

Exits with code `1` if any file is invalid — CI-friendly.

| Flag | Description |
|------|-------------|
| `-o, --output` | `text` (default) or `json` |

## What it checks

- File is readable and valid JSON
- Parses as OSV (`UnmarshalFromJson`)
- Required fields present: `id` and `schema_version`

## Validation flow

```mermaid
flowchart TD
  F["Input file"] --> R{"Readable & valid JSON?"}
  R -->|"no"| E1["✗ error"]
  R -->|"yes"| P{"Parses as OSV?"}
  P -->|"no"| E2["✗ error"]
  P -->|"yes"| C{"id & schema_version<br/>both present?"}
  C -->|"no"| E3["✗ error"]
  C -->|"yes"| OK["✓ valid"]
```

## Decision tree

```mermaid
flowchart TD
  Q["Have OSV file(s)?"] --> V["osv validate *.json"]
  V --> R{"Exit code?"}
  R -->|"0"| OK["✓ all valid"]
  R -->|"1"| Fail["✗ some invalid — errors listed"]
  Fail --> Gate["CI gate fails"]
```

## Where it sits in CI

```mermaid
flowchart LR
  PUSH["Commit / PR"] --> CI["CI pipeline"]
  CI --> VAL["osv validate *.json"]
  VAL --> RC{"exit code"}
  RC -->|"0"| MERGE["Can merge ✓"]
  RC -->|"1"| BLK["Blocks merge ✗"]
```

## Batch semantics: one bad file fails the run

With multiple files the exit code is the logical AND of every result — a single invalid file makes the whole invocation exit `1`, but every file is still checked and reported. This is exactly the behaviour you want in a pre-merge gate over a directory of advisories.

```mermaid
flowchart TD
  START["osv validate a.json b.json c.json"] --> A["check a.json"]
  A --> B["check b.json"]
  B --> C["check c.json"]
  A -.result.-> AGG["aggregate"]
  B -.result.-> AGG
  C -.result.-> AGG
  AGG --> Q{"any invalid?"}
  Q -->|"no"| Z0["exit 0 — all ✓"]
  Q -->|"yes"| Z1["exit 1 — errors listed per file"]
```

## SDK equivalent

```go
raw, _ := os.ReadFile("vulnerability.json")
if !json.Valid(raw) { /* not JSON */ }
v, err := osv.UnmarshalFromJson[any, any](raw)
// then check v.ID != "" && v.SchemaVersion != ""
```

## Cross-references

- [[osv-parse]] — display a valid file's contents
- [[osv-installation]] — install the CLI first
