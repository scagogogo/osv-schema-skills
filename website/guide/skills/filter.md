# osv-filter

Filter OSV data by ecosystem, reference type, or alias pattern.

> **Trigger:** mentions of filtering by package ecosystem (npm, PyPI, Maven), reference type (ADVISORY, FIX), or alias pattern (CVE, GHSA).
> **Skill source:** [`.claude/skills/osv-filter/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-filter/SKILL.md)

## CLI

```bash
osv filter -e PyPI vulnerability.json        # By ecosystem
osv filter -r FIX vulnerability.json         # By reference type
osv filter -a CVE vulnerability.json         # By alias pattern
osv filter -e PyPI -r FIX vulnerability.json # Combine
osv filter -o json -e PyPI vulnerability.json
```

| Flag | Description |
|------|-------------|
| `-e, --ecosystem` | Ecosystem, case-sensitive per OSV spec (`PyPI`, `npm`, `Maven`) |
| `-r, --ref-type` | Reference type, auto-uppercased (`ADVISORY`, `FIX`, `WEB`) |
| `-a, --alias` | Alias prefix, upper-cased before matching (`CVE`, `GHSA`, `CVE-2024`) |
| `-o, --output` | `text` (default) or `json` |

At least one filter flag is required.

The default text output reports the filter, a yes/no `Has ecosystem` check, and matching counts. With no `-o json`, each dimension prints its own block:

```bash
osv filter -e PyPI test_data/GHSA-vxv8-r8q2-63xw.json
```

```text
ID: GHSA-vxv8-r8q2-63xw

Ecosystem filter: PyPI
  Has ecosystem: true
  Matching packages (9):
    - PyPI/tensorflow
    ...
```

`-o json` returns the filtered `affected` entries — note each event object carries only its one non-empty field (`omitempty` at work):

```bash
osv filter -e PyPI -o json test_data/GHSA-vxv8-r8q2-63xw.json
```

```json
{
  "affected": [
    {
      "package": { "ecosystem": "PyPI", "name": "tensorflow" },
      "ranges": [{ "type": "ECOSYSTEM", "events": [ { "introduced": "0" }, { "fixed": "2.7.2" } ] }]
    }
  ]
}
```

## The three filter dimensions

```mermaid
flowchart TD
  DATA["OSV data"] --> E["-e ecosystem<br/>Affected → FilterByEcosystem"]
  DATA --> R["-r reference type<br/>References → FilterByType"]
  DATA --> A["-a alias pattern<br/>Aliases → Filter(prefix)"]
  E --> OUT["filtered result"]
  R --> OUT
  A --> OUT
```

## SDK equivalent

```go
// Ecosystem
pypi := v.Affected.FilterByEcosystem(osv.EcosystemPyPI)
hasNpm := v.Affected.HasEcosystem(osv.EcosystemNpm)

// References
fixes := v.References.FilterByType(osv.ReferenceTypeFix)

// Aliases
cves := v.Aliases.Filter(func(a string) bool {
    return strings.HasPrefix(strings.ToUpper(a), "CVE-")
})
```

## Decision tree

```mermaid
flowchart TD
  Q["What to filter?"] --> Eco["By ecosystem"]
  Q --> Ref["By reference type"]
  Q --> Ali["By alias pattern"]
  Eco --> FE["osv filter -e &lt;eco&gt;"]
  Ref --> FR["osv filter -r &lt;type&gt;"]
  Ali --> FA["osv filter -a &lt;pattern&gt;"]
  FE --> Comb{"Combine?"}
  FR --> Comb
  FA --> Comb
  Comb -->|"yes"| C["chain flags: -e ... -r ..."]
  Comb -->|"no"| Done["results"]
```

## Execution order of combined filters

```mermaid
flowchart LR
  IN["raw data"] --> E["-e ecosystem filter"]
  E --> R["-r reference filter"]
  R --> A["-a alias filter"]
  A --> OUT["final result"]
```

Each flag independently acts on a different slice of the original data; combining them takes the intersection.

## Matching semantics per flag

The three flags do **not** match the same way — this is the most common source of "why did my filter return nothing?".

```mermaid
flowchart TD
  IN["your flag value"] --> E{"which flag?"}
  E -->|"-e ecosystem"| EX["exact, case-sensitive<br/>'PyPI' ✓  ·  'pypi' ✗"]
  E -->|"-r ref-type"| UP["auto-uppercased, then exact<br/>'fix' → 'FIX' ✓"]
  E -->|"-a alias"| PF["prefix match (uppercased)<br/>'CVE' matches 'CVE-2024-1234'"]
```

::: warning `-e` is the strict one
Ecosystem is compared verbatim against the OSV spec's exact casing, so `-e pypi` silently returns nothing. Reference types are forgiving (auto-uppercased) and aliases are prefix-based. When a filter comes back empty, check `-e` casing first against the [Ecosystems](/reference/ecosystems) list.
:::

## Notes

- Ecosystem names are case-sensitive (`PyPI`, not `pypi`)
- Reference types are auto-uppercased in the CLI
- Alias prefixes are upper-cased before matching, so `-a cve` works like `-a CVE`
- `HasEcosystem` returns a bool; `FilterByEcosystem` returns the filtered slice

## Cross-references

- [[osv-parse]] — parse first
- [[osv-query]] — extract fields after filtering
- See [Ecosystems](/reference/ecosystems) for the full constant list
