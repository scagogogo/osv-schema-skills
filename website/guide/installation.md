# Installation

Install the `osv` CLI, the Go SDK, and enable the Claude Code Skills.

## Requirements

- **Go 1.18+** (for SDK and building from source)
- Internet access only needed for `go get` / `go install` / downloading binaries

## Install options at a glance

```mermaid
flowchart TD
  START["What to install?"] --> Q1{"Pick one"}
  Q1 -->|"Pre-built binary<br/>(fastest, no Go)"| BIN["Download tar.gz from Release"]
  Q1 -->|"go install<br/>(needs Go)"| GI["go install ...@latest"]
  Q1 -->|"Build from source"| SRC["git clone + go build"]
  BIN --> VER["osv version verify"]
  GI --> VER
  SRC --> VER
  VER --> OK["Ready ✓"]
```

## Three ways in, one core

Whichever you install, all three access layers resolve to the same Go core — so a fact you learn via the CLI holds in the SDK and the skills.

```mermaid
flowchart TD
  BIN["osv binary"] --> CORE
  SDK["Go import"] --> CORE
  SK["Claude Code skills"] --> SDK
  SK -.shells out.-> BIN
  CORE["osv_schema Go core<br/>parse · validate · filter · query"]
```

## CLI

::: tabs
== Pre-built binary

Pre-built binaries ship for every tag via goreleaser. If the latest release has no pre-built assets yet (e.g. before the first goreleaser-tagged release), fall back to `go install` below.

| OS | Architectures |
|----|---------------|
| Linux | amd64, arm64, arm (v7) |
| macOS | amd64, arm64 |
| Windows | amd64, arm64 |

The archive name is composed from the version, OS, and arch — build yours by filling the same template:

```mermaid
flowchart LR
  T["osv_&lt;version&gt;_&lt;os&gt;_&lt;arch&gt;.&lt;ext&gt;"] --> V["version → v0.1.0"]
  T --> O["os → linux / darwin / windows"]
  T --> A["arch → amd64 / arm64 / arm"]
  T --> E["ext → tar.gz (unix) · zip (windows)"]
```

```bash
# Linux amd64 example — swap version/platform for your case.
# Replace v0.1.0 with the newest tag from the Releases page.
VERSION=v0.1.0
curl -fsSL -o osv.tar.gz \
  https://github.com/scagogogo/osv-schema-skills/releases/download/${VERSION}/osv_${VERSION}_linux_amd64.tar.gz
tar -xzf osv.tar.gz osv
chmod +x osv && sudo mv osv /usr/local/bin/
osv version
```

Verify integrity with the bundled `checksums.txt`:

```bash
sha256sum -c checksums.txt --ignore-missing
```

Releases: <https://github.com/scagogogo/osv-schema-skills/releases>

== Go install

```bash
go install github.com/scagogogo/osv-schema-skills/cmd/osv@latest
osv version
```

`go install` drops the binary in `$(go env GOPATH)/bin`. If `osv version` then says *command not found*, that directory isn't on your `PATH`:

```mermaid
flowchart TD
  GI["go install …@latest"] --> LOC["binary → $(go env GOPATH)/bin"]
  LOC --> Q{"osv version works?"}
  Q -->|yes| OK["ready ✓"]
  Q -->|"command not found"| FIX["add to PATH:<br/>export PATH=\$PATH:\$(go env GOPATH)/bin"]
  FIX --> OK
```

== Build from source

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills
go build -o osv ./cmd/osv/
./osv version
```
:::

## Go SDK

```bash
go get -u github.com/scagogogo/osv-schema-skills
```

```go
import osv "github.com/scagogogo/osv-schema-skills"
```

See the [Go SDK guide](/guide/sdk) for usage.

## Claude Code Skills

The 7 skills activate automatically when Claude Code opens this repo — no install step:

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills
claude   # skills are live
```

Or install as a plugin — the manifest is already in `.claude-plugin/`, so once the marketplace listing is live you can add it directly:

```bash
claude plugin add scagogogo/osv-schema-skills
```

See [Skills Overview](/guide/skills).

## Verify

```bash
osv version                                   # CLI + schema version
osv parse test_data/GHSA-vxv8-r8q2-63xw.json  # parse a sample record
```
