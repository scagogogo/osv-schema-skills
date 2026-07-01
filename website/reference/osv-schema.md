# OSV Schema

The core type models the [OSV Schema](https://ossf.github.io/osv-schema/) (currently `1.4.0`).

## Top-level structure

```mermaid
graph TD
  OSV["OsvSchema&lt;Eco, DB&gt;"] --> ID["id"]
  OSV --> SV["schema_version"]
  OSV --> Time["modified / published"]
  OSV --> W["withdrawn (string)"]
  OSV --> Ali["aliases / related"]
  OSV --> Sum["summary / details"]
  OSV --> SEV["severity: SeveritySlice"]
  OSV --> AFF["affected: AffectedSlice"]
  OSV --> REF["references: References"]
  OSV --> CRED["credits"]
  OSV --> DB["database_specific / ecosystem_specific"]
```

## Required vs optional

| Field | Required | Notes |
|-------|----------|-------|
| `schema_version` | ✅ | Currently `1.4.0` |
| `id` | ✅ | Unique record identifier |
| `modified` | ✅ | Last modification time |
| `published` | ❌ | First publication time |
| `withdrawn` | ❌ | **String**, not `time.Time` |
| `aliases` | ❌ | e.g. CVE-2024-XXXX |
| `affected` | ❌ | But usually present |
| `severity` | ❌ | CVSS v2 / v3 / v4 |

`osv validate` enforces `id` and `schema_version`.

## Affected → package → ranges → events

```mermaid
graph LR
  AFF["Affected"] --> PKG["package<br/>ecosystem · name · purl"]
  AFF --> VER["versions[]"]
  AFF --> RNG["ranges[]"]
  RNG --> EVT["events[]"]
  EVT --> I["introduced"]
  EVT --> F["fixed"]
  EVT --> L["last_affected"]
  EVT --> V["limit"]
```

## Source files

All types live in the root package `osv_schema`:

| File | Contents |
|------|----------|
| `osv_schema.go` | `OsvSchema` top-level type |
| `package.go` | `Package`, `Ecosystem` constants |
| `affected.go` | `Affected`, `AffectedSlice` |
| `severity.go` | `Severity`, `SeveritySlice` |
| `range.go` | `Range` |
| `event.go` | `Event` |
| `references.go` | `References` |
| `aliases.go` | `Aliases` |
| `related.go` | `Related` |
| `credits.go` | `Credits` |
| `unmarshal.go` | `UnmarshalFromJson` / `UnmarshalFromJsonFile` |
