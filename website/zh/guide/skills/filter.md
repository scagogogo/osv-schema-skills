# osv-filter

按生态、引用类型或别名模式过滤 OSV 数据。

> **触发条件：** 提到按包生态过滤（npm、PyPI、Maven）、按引用类型过滤（ADVISORY、FIX），或按别名模式过滤（CVE、GHSA）。
> **技能源码：** [`.claude/skills/osv-filter/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-filter/SKILL.md)

## CLI

```bash
osv filter -e PyPI vulnerability.json        # 按生态
osv filter -r FIX vulnerability.json         # 按引用类型
osv filter -a CVE vulnerability.json         # 按别名模式
osv filter -e PyPI -r FIX vulnerability.json # 组合
osv filter -o json -e PyPI vulnerability.json
```

| 标志 | 说明 |
|------|------|
| `-e, --ecosystem` | 生态，按 OSV 规范区分大小写（`PyPI`、`npm`、`Maven`） |
| `-r, --ref-type` | 引用类型，自动转大写（`ADVISORY`、`FIX`、`WEB`） |
| `-a, --alias` | 别名前缀（`CVE`、`GHSA`、`CVE-2024`） |
| `-o, --output` | `text`（默认）或 `json` |

至少需要一个过滤标志。

## 三个过滤维度

```mermaid
flowchart TD
  DATA["OSV 数据"] --> E["-e 生态<br/>Affected → FilterByEcosystem"]
  DATA --> R["-r 引用类型<br/>References → FilterByType"]
  DATA --> A["-a 别名模式<br/>Aliases → Filter(前缀)"]
  E --> OUT["过滤结果"]
  R --> OUT
  A --> OUT
```

## SDK 等价

```go
// 生态
pypi := v.Affected.FilterByEcosystem(osv.EcosystemPyPI)
hasNpm := v.Affected.HasEcosystem(osv.EcosystemNpm)

// 引用
fixes := v.References.FilterByType(osv.ReferenceTypeFix)

// 别名
cves := v.Aliases.Filter(func(a string) bool {
    return strings.HasPrefix(strings.ToUpper(a), "CVE-")
})
```

## 决策树

```mermaid
flowchart TD
  Q["过滤什么？"] --> Eco["按生态"]
  Q --> Ref["按引用类型"]
  Q --> Ali["按别名模式"]
  Eco --> FE["osv filter -e &lt;eco&gt;"]
  Ref --> FR["osv filter -r &lt;type&gt;"]
  Ali --> FA["osv filter -a &lt;pattern&gt;"]
  FE --> Comb{"要组合吗？"}
  FR --> Comb
  FA --> Comb
  Comb -->|"是"| C["链式标志: -e ... -r ..."]
  Comb -->|"否"| Done["结果"]
```

## 组合过滤的执行顺序

```mermaid
flowchart LR
  IN["原始数据"] --> E["-e 生态过滤"]
  E --> R["-r 引用过滤"]
  R --> A["-a 别名过滤"]
  A --> OUT["最终结果"]
```

各标志独立作用于原数据的不同切片，组合即取交集。

## 注意事项

- 生态名区分大小写（`PyPI`，不是 `pypi`）
- 引用类型在 CLI 中自动转大写
- `HasEcosystem` 返回 bool；`FilterByEcosystem` 返回过滤后的切片

## 交叉引用

- [[osv-parse]] — 先解析
- [[osv-query]] — 过滤后提取字段
- 完整常量列表见 [生态系统](/zh/reference/ecosystems)
