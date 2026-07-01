# osv-affected

分析受影响包与版本范围。

> **触发条件：** 提到受影响包、版本范围、受影响生态，或确定哪些包/版本受影响。
> **技能源码：** [`.claude/skills/osv-affected/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-affected/SKILL.md)

## CLI

```bash
osv parse -v vulnerability.json             # 完整 affected 详情 + 范围
osv filter -e PyPI vulnerability.json       # 收窄到一个生态
osv query --ranges vulnerability.json       # 版本范围
osv query --events vulnerability.json       # 事件时间线
```

## SDK

```go
// 是否存在
v.Affected.HasEcosystem(osv.EcosystemPyPI)

// 过滤
pypi := v.Affected.FilterByEcosystem(osv.EcosystemPyPI)

// 遍历范围与事件
for _, a := range v.Affected {
    fmt.Println(a.Package.Ecosystem, a.Package.Name)
    for _, r := range a.Ranges {
        fmt.Println("  range type:", r.Type)   // SEMVER / ECOSYSTEM / GIT
        for _, e := range r.Events {
            // e.IsIntroduced() / IsFixed() / IsLastAffected() / IsLimit()
        }
    }
}
```

## 结构

```mermaid
graph TD
  AFF["Affected[]"] --> PKG["package<br/>ecosystem · name · purl"]
  AFF --> VER["versions[]"]
  AFF --> RNG["ranges[]"]
  AFF --> ASEV["severity[]（每包）"]
  RNG --> TYPE["type: SEMVER/ECOSYSTEM/GIT"]
  RNG --> EVT["events[]"]
  EVT --> I["introduced"]
  EVT --> F["fixed"]
  EVT --> L["last_affected"]
  EVT --> LM["limit"]
```

## Affected 数据模型

```mermaid
classDiagram
  class Affected {
    +Package Package
    +Versions []string
    +Ranges []Range
    +Severity SeveritySlice
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
  Affected --> Package
  Affected --> Range
  Range --> Event
```

## 决策树

```mermaid
flowchart TD
  Q["关于 affected 要什么？"] --> L["列出所有包 → parse -v"]
  Q --> E["检查生态 → HasEcosystem / filter -e"]
  Q --> R["版本范围 → query --ranges"]
  Q --> M["Maven GAV → query --maven"]
  Q --> EV["事件时间线 → query --events"]
```

## 范围类型对比

```mermaid
flowchart TD
  T["range.type"] --> SE["SEMVER<br/>语义化版本范围"]
  T --> EC["ECOSYSTEM<br/>最常见，按生态版本"]
  T --> GI["GIT<br/>git 提交范围"]
```

- `RangeTypeEcosystem`（`ECOSYSTEM`）最常见；`SEMVER` 和 `GIT` 较少见。

## 注意事项

- `RangeTypeEcosystem`（`ECOSYSTEM`）最常见；`SEMVER` 和 `GIT` 较少
- 每个 event 对象的字段互斥
- `affected[].severity` 是可选的每包 severity，与顶层 `severity` 相互独立

## 交叉引用

- [[osv-filter]] — 按生态收窄 affected
- [[osv-query]] — 提取 ranges/events/maven
- [OSV Schema](/zh/reference/osv-schema) — 完整类型模型
