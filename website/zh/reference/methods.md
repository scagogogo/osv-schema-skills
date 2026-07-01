# 方法清单

SDK 最常用方法的速查表。全部已对照源码核实。

## 方法一览

按接收者类型分组——这就是你日常会用到的全部表面。

```mermaid
mindmap
  root((osv SDK))
    Aliases
      GetCVE
      Filter
    AffectedSlice
      HasEcosystem
      FilterByEcosystem
      Filter
    Package
      IsMaven
      GetGroupID
      GetArtifactID
    SeveritySlice
      GetCVSS3
      GetCVSS2
    Severity
      GetScore
      GetScoreAsFloat
      GetScoreAsPointer
    References
      FilterByType
    Event
      IsIntroduced
      IsFixed
      IsLastAffected
      IsLimit
    package-level
      UnmarshalFromJson
      UnmarshalFromJsonFile
```

## Aliases

| 方法 | 签名 | 说明 |
|------|------|------|
| `GetCVE` | `() string` | 第一个以 `CVE-` 开头的标识 |
| `Filter` | `(func(string) bool) Aliases` | 按谓词过滤别名 |

## AffectedSlice

| 方法 | 签名 | 说明 |
|------|------|------|
| `HasEcosystem` | `(Ecosystem) bool` | 是否有受影响条目匹配该生态 |
| `FilterByEcosystem` | `(Ecosystem) AffectedSlice` | 收窄到一个生态 |
| `Filter` | `(func(*Affected) bool) AffectedSlice` | 自定义谓词过滤 |

## Package

| 方法 | 签名 | 说明 |
|------|------|------|
| `IsMaven` | `() bool` | `Ecosystem == Maven` |
| `GetGroupID` | `() string` | Maven `groupId`（`:` 左侧） |
| `GetArtifactID` | `() string` | Maven `artifactId`（`:` 右侧） |

## SeveritySlice

| 方法 | 签名 | 说明 |
|------|------|------|
| `GetCVSS3` | `() *Severity` | CVSS v3 条目，或 `nil` |
| `GetCVSS2` | `() *Severity` | CVSS v2 条目，或 `nil` |

## Severity

| 方法 | 签名 | 说明 |
|------|------|------|
| `GetScore` | `() float64` | 把 CVSS 分数解析为 `float64` |
| `GetScoreAsFloat` | `() (float64, error)` | 解析分数，向量字符串畸形时返回 error |
| `GetScoreAsPointer` | `() *float64` | 分数指针（用于可空字段） |

## References

| 方法 | 签名 | 说明 |
|------|------|------|
| `FilterByType` | `(...ReferenceType) References` | 按 `ADVISORY`、`FIX` 等过滤（可传多个） |

## Event

| 方法 | 签名 | 说明 |
|------|------|------|
| `IsIntroduced` | `() bool` | 事件标记 introduced 版本 |
| `IsFixed` | `() bool` | 事件标记 fixed 版本 |
| `IsLastAffected` | `() bool` | 事件标记 last_affected 版本 |
| `IsLimit` | `() bool` | 事件标记 range limit |

## Parsing

| 函数 | 签名 | 说明 |
|------|------|------|
| `UnmarshalFromJson` | `([]byte) (*OsvSchema[Eco,DB], error)` | 从字节解析 |
| `UnmarshalFromJsonFile` | `(string) (*OsvSchema[Eco,DB], error)` | 从文件路径解析 |

```go
// 通用解析——两个泛型都用 any
v, err := osv.UnmarshalFromJsonFile[any, any]("vuln.json")

// 或附加生态/库专属数据
v, err := osv.UnmarshalFromJsonFile[MyEco, MyDB]("vuln.json")
```

## 方法调用关系图

```mermaid
graph TD
  OSV["OsvSchema"] --> AL["v.Aliases"]
  OSV --> SEV["v.Severity"]
  OSV --> AFF["v.Affected"]
  OSV --> REF["v.References"]

  AL -->|"GetCVE()"| CVE["string"]
  AL -->|"Filter(f)"| AL2["Aliases"]

  SEV -->|"GetCVSS3()"| S3["*Severity"]
  S3 -->|"GetScore()"| SCORE["float64"]

  AFF -->|"HasEcosystem(e)"| BOOL["bool"]
  AFF -->|"FilterByEcosystem(e)"| AFF2["AffectedSlice"]
  AFF --> PKG["a.Package"]
  PKG -->|"IsMaven()"| MBOOL["bool"]
  PKG -->|"GetGroupID()"| GID["string"]

  REF -->|"FilterByType(t)"| REF2["References"]
```

## 解析与校验的数据流

```mermaid
flowchart LR
  F["文件/字节"] --> U["UnmarshalFromJson[File]"]
  U --> V["*OsvSchema"]
  V --> CHK["检查 v.ID / v.SchemaVersion"]
  CHK --> OK{"非空?"}
  OK -->|"是"| VALID["✓ 有效"]
  OK -->|"否"| INVALID["✗ 无效"]
```

## Maven 坐标拆解

`GetGroupID` / `GetArtifactID` 按第一个 `:` 拆分 Maven 包名。仅当 `IsMaven()` 为真时才有意义。

```mermaid
flowchart LR
  N["package.Name<br/>'com.fasterxml.jackson.core:jackson-databind'"] --> CHK{"IsMaven()?"}
  CHK -->|否| SKIP["不是 Maven 包 → 跳过"]
  CHK -->|是| SPLIT["按第一个 ':' 拆分"]
  SPLIT --> G["GetGroupID()<br/>'com.fasterxml.jackson.core'"]
  SPLIT --> A["GetArtifactID()<br/>'jackson-databind'"]
```

## 一次真实查询，逐个方法

"某条 `GHSA-…` 是不是一个高危的 PyPI 漏洞，修复在哪？"——下面就是智能体（或你的代码）会走的精确方法链。

```mermaid
sequenceDiagram
  participant You as 你 / 智能体
  participant SDK as osv SDK
  You->>SDK: UnmarshalFromJsonFile[any,any](path)
  SDK-->>You: *OsvSchema
  You->>SDK: v.Affected.HasEcosystem(EcosystemPyPI)
  SDK-->>You: true
  You->>SDK: v.Severity.GetCVSS3()
  SDK-->>You: *Severity（向量字符串）
  You->>SDK: v.References.FilterByType(ReferenceTypeFix)
  SDK-->>You: References（修复链接）
  Note over You,SDK: 解析 → 过滤 → 取分 → 修复，全在同一个带类型的内核上
```

## 哪个方法返回什么

```mermaid
flowchart TD
  Q["你手上有什么？"] --> Q1["一个切片 → 想要单个"]
  Q1 --> M1["GetCVSS3 / GetCVSS2 / GetCVE<br/>→ 单个或 nil/空"]
  Q --> Q2["一个切片 → 想要子集"]
  Q2 --> M2["FilterByEcosystem / FilterByType / Filter<br/>→ 一个新切片"]
  Q --> Q3["一个切片 → 是否问题"]
  Q3 --> M3["HasEcosystem<br/>→ bool"]
  Q --> Q4["单个 → 一个派生值"]
  Q4 --> M4["GetScore / GetGroupID / IsMaven / IsFixed<br/>→ 标量"]
```

## 序列化辅助

大多数类型实现了 `sql.Scanner` 和 `driver.Valuer`，因此在 GORM 下能干净地作为 JSON 列存储。复杂嵌套类型（`AffectedSlice`、`SeveritySlice`、`Package`、`Credits`）会自动 marshal/unmarshal 为 JSON。

源码：根包 [`*.go`](https://github.com/scagogogo/osv-schema-skills)
