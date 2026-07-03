# CLI

The `osv` CLI is a thin shell front-end over the Go core — ideal for quick lookups, shell scripting, and CI pipelines.

## Install

See [Quick Start](/guide/quick-start) for install options (pre-built binary, `go install`, or build from source). Pre-built binaries cover:

| OS | Architectures |
|----|---------------|
| Linux | amd64, arm64, arm (v7) |
| macOS | amd64, arm64 |
| Windows | amd64, arm64 |

Download from [GitHub Releases](https://github.com/scagogogo/osv-schema-skills/releases).

## Commands

```mermaid
flowchart LR
  R["osv"] --> P["parse<br/>display fields"]
  R --> V["validate<br/>schema check"]
  R --> F["filter<br/>ecosystem/ref/alias"]
  R --> Q["query<br/>severity/maven/ranges/events"]
  R --> VER["version<br/>CLI + schema version"]
```

## How commands map to the core

```mermaid
flowchart TD
  P["osv parse"] --> U["UnmarshalFromJson"]
  V["osv validate"] --> U
  F["osv filter"] --> U
  Q["osv query"] --> U
  U --> SCH["OsvSchema struct"]
  SCH --> P2["parse: print fields"]
  SCH --> V2["validate: check id/version"]
  SCH --> F2["filter: call FilterBy*"]
  SCH --> Q2["query: call GetCVSS*/ranges"]
```

### `osv parse`

Parse an OSV JSON file and display its fields.

```bash
osv parse vulnerability.json           # Key fields (text)
osv parse -v vulnerability.json        # All fields (dates, details, credits, ranges)
osv parse -o json vulnerability.json   # JSON output
```

| Flag | Description |
|------|-------------|
| `-v, --verbose` | Show all fields: published/modified, withdrawn, related, details, credits, per-range events |
| `-o, --output` | Output format: `text` (default) or `json` |

Output includes ID, schema version, summary, aliases/CVE, severity, affected packages, and references.

### `osv validate`

Validate one or more OSV JSON files against the schema (parses, checks required `id` and `schema_version`).

```bash
osv validate vulnerability.json              # Single file
osv validate file1.json file2.json           # Batch
osv validate -o json vulnerability.json      # JSON output
```

Exits with code `1` if any file is invalid — friendly for CI gating.

| Flag | Description |
|------|-------------|
| `-o, --output` | Output format: `text` (default) or `json` |

### `osv filter`

Filter by affected package ecosystem, reference type, or alias pattern. At least one filter flag required; flags combine.

```bash
osv filter -e PyPI vulnerability.json        # Filter affected by ecosystem
osv filter -r FIX vulnerability.json         # Filter references by type
osv filter -a CVE vulnerability.json         # Filter aliases by pattern
osv filter -e PyPI -r FIX vulnerability.json # Combine
osv filter -o json -e PyPI vulnerability.json
```

| Flag | Description |
|------|-------------|
| `-e, --ecosystem` | Ecosystem name, case-sensitive per OSV spec (`PyPI`, `npm`, `Maven`) |
| `-r, --ref-type` | Reference type, auto-uppercased (`ADVISORY`, `FIX`, `WEB`) |
| `-a, --alias` | Alias prefix pattern, upper-cased before matching (`CVE`, `GHSA`, or `CVE-2024` match case-insensitively) |
| `-o, --output` | `text` (default) or `json` |

### `osv query`

Extract focused sub-information. At least one flag required; flags combine.

```bash
osv query --severity cvss3 vulnerability.json  # CVSS v3 entry + parsed score (0.0 on a vector string)
osv query --severity cvss2 vulnerability.json  # CVSS v2
osv query --maven vulnerability.json           # Maven groupId/artifactId decomposition
osv query --ranges vulnerability.json          # Version ranges per affected package
osv query --events vulnerability.json          # Event timeline (introduced/fixed/…)
osv query --ranges --events vulnerability.json # Combine
```

| Flag | Description |
|------|-------------|
| `--severity` | `cvss3` or `cvss2` |
| `--maven` | Decompose Maven `groupId:artifactId` |
| `--ranges` | Show version ranges |
| `--events` | Show event timeline |
| `-o, --output` | `text` (default) or `json` |

::: tip
`GetScore()` returns `0.0` when the OSV `score` field is a CVSS vector string rather than a number — see [Methods](/reference/methods#severity).
:::

### `osv version`

```bash
osv version
```

Prints the CLI version (injected at build time by goreleaser) and the supported OSV schema version:

```text
osv-cli version: dev
OSV schema version: 1.4.0
```

The `dev` placeholder is replaced with the release tag by goreleaser's ldflags. Unlike the other subcommands, `version` ignores `-o json` — it always prints these two text lines.

## Global flag

| Flag | Description |
|------|-------------|
| `-o, --output` | `text` (default) or `json` — applies to `parse`/`validate`/`filter`/`query`; `version` ignores it |

## Exit code conventions

```mermaid
flowchart LR
  RUN["Run subcommand"] --> RC{"Exit code"}
  RC -->|"0"| OK["Success / file valid"]
  RC -->|"1"| FAIL["Failure / some file invalid"]
  RC -->|"2+"| ERR["Argument or runtime error"]
```

## Typical pipeline

```mermaid
sequenceDiagram
  participant CI as CI pipeline
  participant CLI as osv CLI
  participant Core as Go core
  CI->>CLI: validate *.json
  CLI->>Core: UnmarshalFromJson
  Core-->>CLI: ok / errors
  CLI-->>CI: exit 0 or 1
  CI->>CLI: parse critical.json -o json
  Core-->>CLI: structured data
  CLI-->>CI: JSON to stdout
```

## Composing with `jq`

Because every subcommand speaks `-o json`, the CLI drops straight into a Unix pipeline. The `-o json` output is the same typed core re-marshalled, so field names match the [OSV Schema](/reference/osv-schema) exactly.

```mermaid
flowchart LR
  OSV["osv … -o json"] --> PIPE["| jq '…'"]
  PIPE --> SEL["select / map fields"]
  SEL --> OUT["scalar / filtered JSON"]
```

```bash
# Pull just the CVSS v3 vector string
osv query --severity cvss3 -o json vuln.json | jq -r '.severity.score'

# List every affected ecosystem across a directory, deduplicated
for f in advisories/*.json; do
  osv parse -o json "$f" | jq -r '.affected[].package.ecosystem'
done | sort -u

# Gate CI: fail if any file is invalid, then report the criticals
osv validate advisories/*.json || exit 1
for f in advisories/*.json; do
  osv parse -o json "$f" | jq 'select(.severity != null)'
done
```

::: tip Exit code + JSON compose cleanly
`validate` sets the exit code (`0`/`1`) *and* can emit JSON, so a single command both gates the pipeline and produces a machine-readable report — no second parse needed.
:::
