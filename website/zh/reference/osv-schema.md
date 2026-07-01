# OSV Schema

核心类型建模 [OSV Schema](https://ossf.github.io/osv-schema/)（当前 `1.4.0`）。

## 顶层结构

```mermaid
graph TD
  OSV["OsvSchema&lt;Eco, DB&gt;"] --> ID["id"]
  OSV --> SV["schema_version"]
  OSV --> Time["modified / published"]
  OSV --> W["withdrawn（字符串）"]
  OSV --> Ali["aliases / related"]
  OSV --> Sum["summary / details"]
  OSV --> SEV["severity: SeveritySlice"]
  OSV --> AFF["affected: AffectedSlice"]
  OSV --> REF["references: References"]
  OSV --> CRED["credits"]
  OSV --> DB["database_specific / ecosystem_specific"]
```

## 必需 vs 可选

| 字段 | 必需 | 说明 |
|------|------|------|
| `schema_version` | ✅ | 当前 `1.4.0` |
| `id` | ✅ | 唯一记录标识 |
| `modified` | ✅ | 最后修改时间 |
| `published` | ❌ | 首次发布时间 |
| `withdrawn` | ❌ | **字符串**，非 `time.Time` |
| `aliases` | ❌ | 如 CVE-2024-XXXX |
| `affected` | ❌ | 但通常存在 |
| `severity` | ❌ | CVSS v2 / v3 / v4 |

`osv validate` 强制 `id` 和 `schema_version`。

## 完整类型关系图

```mermaid
classDiagram
  class OsvSchema {
    +SchemaVersion string
    +ID string
    +Modified time.Time
    +Published time.Time
    +Withdrawn string
    +Aliases Aliases
    +Related Related
    +Summary string
    +Details string
    +Severity SeveritySlice
    +Affected AffectedSlice
    +References References
    +Credits *Credits
  }
  class Affected {
    +Package Package
    +Ranges []Range
    +Versions []string
  }
  class Package {
    +Ecosystem Ecosystem
    +Name string
    +Purl string
  }
  class Range {
    +Type RangeType
    +Events []Event
  }
  class Event {
    +Introduced string
    +Fixed string
    +LastAffected string
    +Limit string
  }
  class Severity {
    +Type SeverityType
    +Score string
  }
  class References {
    +[]Reference items
  }
  OsvSchema --> Affected
  OsvSchema --> Severity
  OsvSchema --> References
  Affected --> Package
  Affected --> Range
  Range --> Event
```

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

## 一条记录的生命周期

```mermaid
stateDiagram-v2
  [*] --> 发布: published
  发布 --> 修改: modified（每次更新）
  修改 --> 修改: 持续维护
  修改 --> 撤回: withdrawn 非空
  撤回 --> [*]
```

## 字段速查（按用途）

```mermaid
flowchart TD
  USE["你想知道..."] --> A1["是哪个漏洞？"] --> R1["id / aliases(CVE)"]
  USE --> A2["多严重？"] --> R2["severity (CVSS)"]
  USE --> A3["影响什么？"] --> R3["affected[].package"]
  USE --> A4["哪些版本？"] --> R4["affected[].ranges / events"]
  USE --> A5["怎么修？"] --> R5["references (FIX) + events.fixed"]
  USE --> A6["还撤回了吗？"] --> R6["withdrawn 非空?"]
```

## 源文件

所有类型在根包 `osv_schema` 中：

| 文件 | 内容 |
|------|------|
| `osv_schema.go` | `OsvSchema` 顶层类型 |
| `package.go` | `Package`、`Ecosystem` 常量 |
| `affected.go` | `Affected`、`AffectedSlice` |
| `severity.go` | `Severity`、`SeveritySlice` |
| `range.go` | `Range` |
| `event.go` | `Event` |
| `references.go` | `References` |
| `aliases.go` | `Aliases` |
| `related.go` | `Related` |
| `credits.go` | `Credits` |
| `unmarshal.go` | `UnmarshalFromJson` / `UnmarshalFromJsonFile` |
