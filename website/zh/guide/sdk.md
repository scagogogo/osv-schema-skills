# Go SDK

Go SDK 是 CLI 和技能之下的类型安全基石。当你要把 OSV 解析/过滤/查询嵌入 Go 应用时用它。

## 安装

```bash
go get -u github.com/scagogogo/osv-schema-skills
```

```go
import osv "github.com/scagogogo/osv-schema-skills"
```

## 快速开始

```go
package main

import (
    "fmt"
    "log"

    osv "github.com/scagogogo/osv-schema-skills"
)

func main() {
    // 从 JSON 文件解析 OSV 数据
    v, err := osv.UnmarshalFromJsonFile[any, any]("vulnerability.json")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s\n", v.ID)
    fmt.Printf("Summary: %s\n", v.Summary)

    // 从 aliases 取 CVE
    if cve := v.Aliases.GetCVE(); cve != "" {
        fmt.Printf("CVE: %s\n", cve)
    }

    // 检查是否影响某生态
    if v.Affected.HasEcosystem("npm") {
        fmt.Println("影响 npm 包")
    }

    // 取 CVSS v3 分数
    if cvss3 := v.Severity.GetCVSS3(); cvss3 != nil {
        fmt.Printf("CVSS v3: %.1f\n", cvss3.GetScore())
    }
}
```

## 从 JSON 到代码：对象生命周期

```mermaid
flowchart LR
  F["vulnerability.json"] --> U["UnmarshalFromJsonFile"]
  U --> V["*OsvSchema 结构体"]
  V --> Q["查询字段:<br/>v.ID / v.Summary / v.Aliases"]
  V --> F2["过滤:<br/>v.Affected.FilterByEcosystem"]
  V --> S["取 severity:<br/>v.Severity.GetCVSS3"]
```

## 核心类型

```go
type OsvSchema[EcosystemSpecific, DatabaseSpecific any] struct {
    SchemaVersion    string
    ID               string
    Modified         time.Time
    Published        time.Time
    Withdrawn        string // string，不是 time.Time——非空即表示已撤回
    Aliases          Aliases
    Related          Related
    Summary          string
    Details          string
    Severity         SeveritySlice
    Affected         AffectedSlice[EcosystemSpecific, DatabaseSpecific]
    References       References
    DatabaseSpecific DatabaseSpecific
    Credits          *Credits
}
```

泛型参数 `EcosystemSpecific` 和 `DatabaseSpecific` 让你按生态或漏洞库附加自定义数据。通用解析用 `any`。

```mermaid
graph TD
  subgraph "泛型扩展点"
    E["EcosystemSpecific<br/>每生态数据"]
    D["DatabaseSpecific<br/>每库数据"]
  end
  OSV["OsvSchema&lt;Eco, DB&gt;"] --> E
  OSV --> D
```

## 类型关系一图

```mermaid
classDiagram
  class OsvSchema {
    +ID string
    +SchemaVersion string
    +Aliases Aliases
    +Severity SeveritySlice
    +Affected AffectedSlice
    +References References
  }
  class Aliases {
    +GetCVE() string
    +Filter(func) Aliases
  }
  class SeveritySlice {
    +GetCVSS3() *Severity
    +GetCVSS2() *Severity
  }
  class AffectedSlice {
    +HasEcosystem(Ecosystem) bool
    +FilterByEcosystem(Ecosystem) AffectedSlice
  }
  class Package {
    +IsMaven() bool
    +GetGroupID() string
    +GetArtifactID() string
  }
  class References {
    +FilterByType(...ReferenceType) References
  }
  OsvSchema --> Aliases
  OsvSchema --> SeveritySlice
  OsvSchema --> AffectedSlice
  OsvSchema --> References
  AffectedSlice --> Package
```

## 关键方法

完整表见 [参考 → 方法清单](/zh/reference/methods)。要点：

| 类型 | 方法 | 说明 |
|------|------|------|
| `OsvSchema` | `Affected.HasEcosystem(eco)` | 检查是否影响某生态 |
| `AffectedSlice` | `FilterByEcosystem(eco)` | 过滤受影响包 |
| `Aliases` | `GetCVE()` | 取第一个 CVE 标识 |
| `SeveritySlice` | `GetCVSS3()` / `GetCVSS2()` | 取 CVSS severity 条目 |
| `Severity` | `GetScore()` | 解析分数为 float64 |
| `References` | `FilterByType(t)` | 按引用类型过滤 |
| `Package` | `IsMaven()` / `GetGroupID()` / `GetArtifactID()` | Maven 拆分 |

## 序列化

每个核心类型都带 `json`、`yaml`、`mapstructure`、`db`、`bson`、`gorm` 标签——JSON、YAML、mapstructure、GORM 和 MongoDB（BSON）开箱即用。

```mermaid
flowchart LR
  T["带标签的核心类型"] --> JSON["json ✓"]
  T --> YAML["yaml ✓"]
  T --> MS["mapstructure ✓"]
  T --> GORM["gorm / db ✓"]
  T --> BSON["bson ✓"]
```

## 带类型的厂商字段——实战示例

`[any, any]` 适合大多数解析，但当你反复读取一个已知形态的 `database_specific`（比如 GitHub 的公告块）时，给它一个具体类型，编译器就会替你检查字段访问。

```go
// 定义你关心的厂商块形态
type GHSA struct {
    Severity   string   `json:"severity"`
    CWEIDs     []string `json:"cwe_ids"`
    GitHubURL  string   `json:"github_reviewed_at"`
}

// 用具体类型作为 DatabaseSpecific 解析
v, err := osv.UnmarshalFromJsonFile[any, GHSA]("ghsa.json")
if err != nil {
    log.Fatal(err)
}
// v.DatabaseSpecific 现在是带类型的 GHSA——无需 map[string]any 强转
fmt.Println(v.DatabaseSpecific.Severity, v.DatabaseSpecific.CWEIDs)
```

```mermaid
flowchart LR
  RAW["database_specific: { … }"] --> ANY["[any, any]<br/>→ map[string]any<br/>（每次访问都要强转）"]
  RAW --> TYPED["[any, GHSA]<br/>→ 带类型结构体<br/>（编译期检查）"]
  ANY --> COST["运行时类型断言"]
  TYPED --> SAFE["字段访问安全"]
```

::: tip 只为需要的类型付费
两个参数相互独立。只给 `DatabaseSpecific` 定类型、`EcosystemSpecific` 留 `any`（反之亦然）——不必两个块都建模才能给其中一个上类型。
:::

## 设计要点

- **成功时永不 nil，出错时必为 nil**——`UnmarshalFromJsonFile` / `UnmarshalFromJson` 失败时返回 `(nil, err)`，成功时返回非 nil 的 `*OsvSchema`。碰指针前先检查 `err`。
- **Withdrawn 是字符串**——不是 `time.Time`。用非空字符串判断撤回状态。
- **数据库策略**——简单字段做列；复杂嵌套结构（`AffectedSlice`、`SeveritySlice`）经 GORM serializer 存为 JSON 字符串。

## 环境要求

- Go 1.18+
