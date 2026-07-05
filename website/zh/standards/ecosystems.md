# 生态系统命名标准

OSV 用**生态限定的包名**标识受影响组件。本页解释命名约定、为什么大小写敏感、本工具箱支持哪些生态。

---

## 为什么要有生态？

光名字 `requests` 是歧义的——可能是 PyPI 包、RubyGem 或 npm 模块。OSV 要求每个包携带 `ecosystem` 字段，使 `ecosystem` + `name` 组合在全球唯一无歧义：

```json
{
  "package": {
    "ecosystem": "PyPI",
    "name": "requests"
  }
}
```

这是 `osv filter -e PyPI` 的基础——没有生态，过滤"PyPI 包"就得靠命名约定猜。

---

## 大小写敏感

按 OSV 规范，生态名**区分大小写**。本工具箱精确匹配：

| 生态 | 正确 | 错误 |
|------|------|------|
| Python | `PyPI` | ~~`pypi`~~、~~`PyPi`~~ |
| Node.js | `npm` | ~~`NPM`~~、~~`Npm`~~ |
| Go | `Go` | ~~`go`~~、~~`GO`~~ |
| Java | `Maven` | ~~`maven`~~ |
| .NET | `NuGet` | ~~`nuget`~~ |
| Rust | `crates.io` | ~~`Crates.io`~~ |

`osv filter -e` 标志**原样**传值——不做大小写归一化。这是有意的：规范说大小写敏感，工具箱就忠实如此。

::: tip 引用类型和别名则大小写不敏感
只有 `ecosystem` 区分大小写。`-r`（引用类型）和 `-a`（别名模式）标志在匹配前自动转大写，所以 `fix`、`Fix`、`FIX` 都行；`cve`、`Cve`、`CVE` 都行。
:::

---

## 支持的生态

本工具箱为 OSV 标准化的生态定义了常量：

| 常量 | 值 | 说明 |
|------|-----|------|
| `EcosystemGo` | `Go` | 名字是 Go module path |
| `EcosystemNpm` | `npm` | 名字是 npm 包名 |
| `EcosystemOSSFuzz` | `OSS-Fuzz` | OSS-Fuzz 项目报告 |
| `EcosystemPyPI` | `PyPI` | 名字是规范化的 PyPI 包名 |
| `EcosystemRubyGems` | `RubyGems` | 名字是 gem 名 |
| `EcosystemCratesIo` | `crates.io` | 名字是 Rust crate 名 |
| `EcosystemPackagist` | `Packagist` | PHP 包管理器 |
| `EcosystemMaven` | `Maven` | 名字是 `groupId:artifactId` |
| `EcosystemNuGet` | `NuGet` | .NET 包名 |
| `EcosystemLinux` | `Linux` | 仅 `name: Kernel` |
| `EcosystemDebian` | `Debian` | 源包名；可选 `:<RELEASE>` 后缀（如 `Debian:7`） |
| `EcosystemAlpine` | `Alpine` | 源包名；**必需** `:v<RELEASE>` 后缀（如 `Alpine:v3.16`） |
| `EcosystemHex` | `Hex` | Erlang/Elixir 包 |
| `EcosystemAndroid` | `Android` | Android 组件名（Framework、Media Framework、Kernel……） |
| `EcosystemGitHubActions` | `GitHub Actions` | 名字是 action 的 `{owner}/{repo}` |
| `EcosystemPub` | `Pub` | Dart 包名 |
| `EcosystemConanCenter` | `ConanCenter` | C/C++ Conan 包名 |
| `EcosystemRocky` | `Rocky` | 源包名；可选 `:<RELEASE>` 后缀 |
| `EcosystemAlmaLinux` | `AlmaLinux` | 源包名；可选 `:<RELEASE>` 后缀 |

完整清单及详情见 [生态系统参考](/zh/reference/ecosystems)。

---

## 发行版范围限定

某些 Linux 发行版生态在生态字符串本身（不是包名）携带 `:<RELEASE>` 后缀，把记录限定到特定发行版：

| 生态 | 后缀形式 | 示例 | 是否必需？ |
|------|---------|------|-----------|
| `Debian` | `:<RELEASE>`（数字） | `Debian:7` | 可选 |
| `Alpine` | `:v<RELEASE>`（数字，`v` 前缀） | `Alpine:v3.16` | **必需** |
| `Rocky` | `:<RELEASE>`（数字） | `Rocky:9` | 可选 |
| `AlmaLinux` | `:<RELEASE>`（数字） | `AlmaLinux:9` | 可选 |

这是 OSV 规范的一部分，不是工具箱的约定——后缀放在 `ecosystem` 字段里，所以 `osv filter -e "Debian:7"` 只匹配 Debian 7 的记录。注意 Alpine 是最严格的：其发行版后缀**必须**有，且要带 `v` 前缀。

---

## Maven 的特殊命名

Maven 包名是 `groupId:artifactId`：

```json
{
  "package": {
    "ecosystem": "Maven",
    "name": "org.apache.logging.log4j:log4j-core"
  }
}
```

CLI 的 `--maven` 标志将其拆成独立的 `groupId` 和 `artifactId` 字段：

```bash
osv query --maven -o json vuln.json
```

```json
{
  "maven": {
    "groupId": "org.apache.logging.log4j",
    "artifactId": "log4j-core"
  }
}
```

---

## 每个生态的版本比较

每个生态有自己的版本排序规则。OSV 的 `range.type` 字段说明用哪套规则：

| `range.type` | 比较方式 | 用于 |
|--------------|---------|------|
| `ECOSYSTEM` | 该生态的原生排序 | PyPI、npm、Maven…… |
| `SEMVER` | 严格 semver | 遵循 semver 的库 |
| `GIT` | 提交哈希（拓扑） | 按 commit 追踪的仓库 |

本工具箱暴露 `type` 和 `events`，但**不**自己做比较——你把事件喂给适合该生态的版本比较库。见 [版本范围语义](/zh/advanced/version-ranges)。

---

## 为什么用注册表而非自由文本？

固定的生态注册表意味着：

- **无拼写错误**——`PyPI` 是唯一合法形式；`Pypi` 不会静默匹配空
- **无重复**——每个数据库对同一生态用相同的字符串
- **工具可特化**——如 Maven 的 `groupId:artifactId` 拆分，正因为工具箱知道 Maven 约定才可能

---

## 另见

- [生态系统参考](/zh/reference/ecosystems) —— 完整常量清单
- [osv-filter 技能](/zh/guide/skills/filter) —— 过滤如何使用生态
- [版本范围语义](/zh/advanced/version-ranges) —— 每生态比较
- [官方 OSV 生态列表](https://ossf.github.io/osv-schema/#affectedpackage-field) —— 权威来源