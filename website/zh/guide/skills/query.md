# osv-query

提取特定子信息：CVSS severity、Maven 拆分、版本范围、事件时间线。

> **触发条件：** 查询 CVSS 分数、Maven groupId/artifactId、版本范围，或从 OSV 数据中聚焦提取。
> **技能源码：** [`.claude/skills/osv-query/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-query/SKILL.md)

## CLI

```bash
osv query --severity cvss3 vulnerability.json  # CVSS v3 条目 + 解析分数
osv query --severity cvss2 vulnerability.json  # CVSS v2
osv query --maven vulnerability.json           # Maven groupId/artifactId
osv query --ranges vulnerability.json          # 版本范围
osv query --events vulnerability.json          # 事件时间线
osv query --ranges --events vulnerability.json # 组合
```

| 标志 | 说明 |
|------|------|
| `--severity` | `cvss3` 或 `cvss2` |
| `--maven` | 拆分 Maven `groupId:artifactId` |
| `--ranges` | 显示版本范围 |
| `--events` | 显示事件时间线 |
| `-o, --output` | `text`（默认）或 `json` |

至少需要一个标志。

## 四个提取维度

```mermaid
flowchart TD
  DATA["OSV 数据"] --> SEV["--severity<br/>Severity → GetCVSS3/2"]
  DATA --> MAV["--maven<br/>Package → GetGroupID/ArtifactID"]
  DATA --> RNG["--ranges<br/>Affected[].Ranges"]
  DATA --> EVT["--events<br/>Range.Events 时间线"]
  SEV --> OUT["提取结果"]
  MAV --> OUT
  RNG --> OUT
  EVT --> OUT
```

## SDK 等价

```go
// Severity
if s := v.Severity.GetCVSS3(); s != nil { fmt.Println(s.GetScore()) }

// Maven
for _, a := range v.Affected {
    if a.Package.IsMaven() {
        fmt.Println(a.Package.GetGroupID(), a.Package.GetArtifactID())
    }
}

// Ranges & events
for _, a := range v.Affected {
    for _, r := range a.Ranges {
        for _, e := range r.Events {
            // e.IsIntroduced() / IsFixed() / IsLastAffected() / IsLimit()
        }
    }
}
```

## 决策树

```mermaid
flowchart TD
  Q["提取什么？"] --> Sev["CVSS severity"]
  Q --> Mav["Maven GAV"]
  Q --> Rng["版本范围"]
  Q --> Evt["事件时间线"]
  Sev --> S["osv query --severity cvss3|cvss2"]
  Mav --> M["osv query --maven"]
  Rng --> R["osv query --ranges"]
  Evt --> E["osv query --events"]
  R --> Comb{"要组合吗？"}
  E --> Comb
  Comb -->|"是"| C["--ranges --events"]
```

## 版本范围与事件的关系

```mermaid
graph TD
  AFF["Affected"] --> RNG["ranges[]"]
  RNG --> TYPE["type: SEMVER / ECOSYSTEM / GIT"]
  RNG --> EVT["events[]"]
  EVT --> I["introduced<br/>受影响起始版本"]
  EVT --> F["fixed<br/>已修复版本"]
  EVT --> L["last_affected<br/>最后受影响版本"]
  EVT --> LM["limit<br/>范围上限"]
```

事件字段在每个 event 对象里互斥——一个 event 只会是 introduced/fixed/last_affected/limit 之一。

## 注意事项

- `GetCVSS3()` / `GetCVSS2()` 在 severity 类型缺失时返回 `nil`
- 当 OSV 的 `score` 是 CVSS 向量字符串而非数字时，`GetScore()` 返回 `0.0`——错误处理用 `GetScoreAsFloat()`
- Maven 拆分只适用于 `Maven` 生态的包
- 事件字段互斥：每个 event 对象只有 `introduced`/`fixed`/`last_affected`/`limit` 之一

## 交叉引用

- [[osv-parse]] — 先完整解析
- [[osv-severity]] — 更深的 severity 分析
- [[osv-affected]] — 更深的 affected/range 分析
