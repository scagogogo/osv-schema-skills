# osv-parse

Parse an OSV JSON file and display structured vulnerability data.

> **Trigger:** mentions of OSV parsing, vulnerability JSON reading, CVE/GHSA data extraction, or when a user provides an OSV JSON file path.
> **Skill source:** [`.claude/skills/osv-parse/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-parse/SKILL.md)

## CLI

```bash
osv parse vulnerability.json           # Key fields (text)
osv parse -v vulnerability.json        # All fields (dates, details, ranges, credits)
osv parse -o json vulnerability.json   # JSON output
```

| Flag | Description |
|------|-------------|
| `-v, --verbose` | Show all fields |
| `-o, --output` | `text` (default) or `json` |

## SDK equivalent

```go
v, err := osv.UnmarshalFromJsonFile[any, any]("vulnerability.json")
fmt.Println(v.ID, v.Summary, v.Aliases.GetCVE())
```

## Decision tree

```mermaid
flowchart TD
  Q["Got an OSV JSON file?"] --> P["osv parse file.json"]
  P --> Need{"Need everything?"}
  Need -->|"no"| Done["Key fields shown"]
  Need -->|"yes"| V["osv parse -v file.json"]
  V --> Done2["All fields shown"]
```

## Output structure

```mermaid
flowchart TD
  OUT["parse output"] --> ID["ID"]
  OUT --> SV["schema version"]
  OUT --> SUM["summary / aliases / CVE"]
  OUT --> SEV["severity (CVSS)"]
  OUT --> AFF["affected packages"]
  OUT --> REF["references"]
  OUT --> VV["with -v also:<br/>published/modified/withdrawn/<br/>related/details/per-range events/credits"]
```

## What it prints

ID, schema version, summary, aliases/CVE, severity, affected packages, references. With `-v` it additionally shows published/modified dates, withdrawn, related, details, per-range events, and credits.

## What happens under the hood

`osv parse` is a thin shell over the SDK's `UnmarshalFromJsonFile` — the same call you'd make in Go. The text/JSON rendering is the only difference from the SDK path.

```mermaid
sequenceDiagram
  participant U as You / agent
  participant CLI as osv parse
  participant SDK as UnmarshalFromJsonFile[any,any]
  participant R as renderer
  U->>CLI: osv parse file.json [-v] [-o json]
  CLI->>SDK: read file → decode JSON
  SDK-->>CLI: *OsvSchema (typed core)
  CLI->>R: -o text → key fields (or all with -v)
  CLI->>R: -o json → re-marshal typed core
  R-->>U: printed output
```

## Text vs JSON: which to pick

```mermaid
flowchart TD
  WHO{"Who reads the output?"} -->|"a human at a terminal"| T["-o text (default)<br/>compact, key fields"]
  WHO -->|"a script / agent / pipe"| J["-o json<br/>full structure, machine-parseable"]
  T --> TV{"need every field?"}
  TV -->|yes| ADDV["add -v"]
  TV -->|no| OK["done"]
```

`-o json` re-marshals the typed `OsvSchema` core, so the output is a standard OSV record — field names match the [OSV Schema](/reference/osv-schema) exactly:

```bash
osv parse -o json vulnerability.json | jq '{id, summary, severity, affected}'
```

```json
{
  "id": "GHSA-vxv8-r8q2-63xw",
  "summary": "TensorFlow vulnerable to `CHECK` fail in `FractionalMaxPoolGrad`",
  "severity": [{ "type": "CVSS_V3", "score": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H" }],
  "affected": [{ "package": { "ecosystem": "PyPI", "name": "tensorflow" }, "ranges": [...] }]
}
```

::: tip Parsing never mutates the file
`parse` only reads. It decodes into the typed core and prints — it never writes back. To check that a file is *well-formed* before parsing, reach for [[osv-validate]]; a malformed file makes `parse` exit non-zero with the decode error.
:::

## Cross-references

- [[osv-validate]] — check the file is schema-valid first
- [[osv-filter]] / [[osv-query]] — narrow or extract from parsed data
