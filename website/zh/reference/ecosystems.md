# 生态系统

SDK 把全部 19 个 OSV 生态定义为类型化常量——杜绝字符串拼错的隐患。

## 全部列表

| 常量 | 值 | 说明 |
|------|----|------|
| `EcosystemGo` | `Go` | Go 模块路径 |
| `EcosystemNpm` | `npm` | NPM 包名 |
| `EcosystemPyPI` | `PyPI` | 规范化的 PyPI 名 |
| `EcosystemRubyGems` | `RubyGems` | Gem 名 |
| `EcosystemCratesIo` | `crates.io` | Rust crate |
| `EcosystemPackagist` | `Packagist` | PHP |
| `EcosystemMaven` | `Maven` | Java——name 为 `groupId:artifactId` |
| `EcosystemNuGet` | `NuGet` | .NET |
| `EcosystemHex` | `Hex` | Erlang/Elixir |
| `EcosystemPub` | `Pub` | Dart |
| `EcosystemLinux` | `Linux` | 仅内核 |
| `EcosystemDebian` | `Debian` | 可能带 `:<RELEASE>` 后缀 |
| `EcosystemAlpine` | `Alpine` | 需带 `:v<RELEASE>` 后缀 |
| `EcosystemRocky` | `Rocky` | 可能带 `:<RELEASE>` 后缀 |
| `EcosystemAlmaLinux` | `AlmaLinux` | 可能带 `:<RELEASE>` 后缀 |
| `EcosystemAndroid` | `Android` | 组件名 |
| `EcosystemOSSFuzz` | `OSS-Fuzz` | Fuzz 目标 |
| `EcosystemConanCenter` | `ConanCenter` | C/C++ |
| `EcosystemGitHubActions` | `GitHub Actions` | `{owner}/{repo}` |

## 按类别分组

```mermaid
flowchart TD
  ALL["19 个生态常量"] --> LANG["语言生态"]
  ALL --> OS["操作系统发行版"]
  ALL --> PLAT["平台/其他"]

  LANG --> L1["npm (JS)"]
  LANG --> L2["PyPI (Python)"]
  LANG --> L3["Maven (Java)"]
  LANG --> L4["Go / RubyGems / crates.io"]
  LANG --> L5["NuGet / Hex / Pub / Packagist / ConanCenter"]

  OS --> O1["Linux"]
  OS --> O2["Debian / Alpine"]
  OS --> O3["Rocky / AlmaLinux"]

  PLAT --> P1["Android"]
  PLAT --> P2["OSS-Fuzz"]
  PLAT --> P3["GitHub Actions"]
```

## 生态 → 版本方案

一个生态不只是给注册表起名——它还隐含了*版本如何排序*，而这正是 `range.type` 需要的（见 [RangeType](/zh/reference/osv-schema#rangetype-——-版本如何比较)）。正因如此，你不能用普通字符串 `<` 来比较版本。

```mermaid
flowchart LR
  subgraph SEMVER["类 SEMVER 优先级"]
    npm & CratesIo["crates.io"] & Go
  end
  subgraph ECO["生态原生排序"]
    PyPI["PyPI (PEP 440)"] & Maven["Maven（点分+限定符）"] & Debian["Debian (dpkg)"] & RubyGems
  end
  subgraph GIT["GIT 提交图"]
    Linux & Android
  end
```

## 命名约定与后缀

`package.name` 字符串并非自由格式——每个生态都有自己的形状，有些还带强制或可选的 `:<release>` 后缀。弄错这一点是 `HasEcosystem` 匹配静默失败的最常见原因。

```mermaid
flowchart TD
  E["生态"] --> C{"name 形状"}
  C -->|"Maven"| M["groupId:artifactId<br/>→ GetGroupID / GetArtifactID"]
  C -->|"GitHub Actions"| GA["owner/repo"]
  C -->|"Go"| GO["完整模块路径"]
  C -->|"Alpine"| AL["Alpine:v3.18（后缀必需）"]
  C -->|"Debian / Rocky / AlmaLinux"| DEB["Base 或 Base:release（后缀可选）"]
  C -->|"多数语言包"| PL["普通注册表名"]
```

::: tip 发行版后缀属于生态字符串，而非 name
对 Alpine 而言 release 后缀（`Alpine:v3.18`）是**必需**的；对 Debian/Rocky/AlmaLinux 则是可选的。后缀挂在*生态*字符串上，所以精确匹配的 `HasEcosystem(EcosystemAlpine)` 不会匹配 `Alpine:v3.18`。只有在你把后缀归一化掉之后，才应与基础常量比较。
:::

## 用法

```go
// 检查单个生态
if v.Affected.HasEcosystem(osv.EcosystemPyPI) {
    // ...
}

// 过滤受影响条目
pypiAffected := v.Affected.FilterByEcosystem(osv.EcosystemPyPI)

// Maven 拆分
for _, a := range v.Affected {
    if a.Package != nil && a.Package.IsMaven() {
        fmt.Println(a.Package.GetGroupID())    // groupId
        fmt.Println(a.Package.GetArtifactID()) // artifactId
    }
}
```

消费 `Ecosystem` 常量的入口有两个——一个问是/否，另一个返回缩小后的切片。两者都遍历同一个 `AffectedSlice`：

```mermaid
flowchart TD
  A["v.Affected: AffectedSlice"] --> H["HasEcosystem(eco)"]
  A --> F["FilterByEcosystem(eco)"]
  H --> W["遍历条目，<br/>匹配 entry.Package.Ecosystem == eco"]
  F --> W
  W -->|"有任一匹配"| BOOL["HasEcosystem → true"]
  W -->|"无匹配"| BOOL2["HasEcosystem → false"]
  W --> KEEP["保留匹配条目"]
  KEEP --> SUB["FilterByEcosystem → AffectedSlice（子集）"]
```

::: tip `Ecosystem` 是带类型字符串，不是自由字符串
传常量如 `osv.EcosystemPyPI`，而非字面量 `"PyPI"`。常量带有 OSV 规范要求的精确大小写，于是比较区分大小写且杜绝拼写错误。
:::

::: warning `FilterByEcosystem` 假定 `Package` 非 nil
`HasEcosystem` 会跳过 `Package` 为 `nil` 的条目（它先检查 `item.Package != nil`）。`FilterByEcosystem` 则**不**检查——它直接解引用 `affected.Package.Ecosystem`，故某条 `affected` 的 `package` 为 `null`/缺失时会 panic。实践中每条合规的 OSV `affected` 条目都带 `package`，但若你解析的是不受信任的数据，请先用 [[osv-validate]] 校验，或自行对切片做防护。
:::

## Maven name 的拆分

```mermaid
flowchart LR
  NAME["package.name = 'org.apache:commons'"] --> SEP["按第一个 ':' 拆分"]
  SEP --> G["GetGroupID → 'org.apache'"]
  SEP --> A["GetArtifactID → 'commons'"]
```

`GetGroupID` / `GetArtifactID` 用 `strings.SplitN(name, ":", 2)`，故 `org.apache.commons:collections4` 这样的 name 会拆成 `org.apache.commons`（group）和 `collections4`（artifact）——只有**第一个** `:` 是分隔符。两者都对 nil 安全：`nil` 的 `Package` 或不含 `:` 的 name 返回 `""`，而非 panic。它们对任意 `Package` 都可调用，但只有 `IsMaven()` 为真时才有意义。

源码：[`package.go`](https://github.com/scagogogo/osv-schema-skills/blob/main/package.go)
