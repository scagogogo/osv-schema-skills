---
layout: home

hero:
  name: OSV Schema Skills
  text: AI 原生的漏洞数据工具箱
  tagline: 把一段提示词粘贴进 Claude Code 或 Codex——智能体会自动安装，并开始解析、校验、查询 OSV 漏洞数据。Go SDK + CLI + 6 个技能，共用同一个 Go 内核。
  image:
    src: /architecture.svg
    alt: 三层架构
  actions:
    - theme: brand
      text: 🤖 AI Agent 接入提示词
      link: /zh/guide/ai-agent
    - theme: alt
      text: 快速开始
      link: /zh/guide/quick-start
    - theme: alt
      text: GitHub
      link: https://github.com/scagogogo/osv-schema-skills

features:
  - icon: 🤖
    title: AI Agent 优先
    details: 6 个专用 Claude Code 技能会在漏洞任务上自动触发。无需自定义集成——克隆仓库即可激活。
    link: /zh/guide/skills
    linkText: 浏览技能
  - icon: 🖥️
    title: 跨平台 CLI
    details: Linux / macOS / Windows（amd64、arm64、arm）全平台预编译二进制。一条 curl 即可下载，AI 智能体随处可运行。
    link: /zh/guide/cli
    linkText: CLI 命令
  - icon: 📦
    title: 类型安全的 Go SDK
    details: 泛型 OsvSchema[Eco, DB]，开箱即用支持 JSON、YAML、GORM、BSON 序列化。通用解析用 any 即可。
    link: /zh/guide/sdk
    linkText: SDK 参考
  - icon: 🛡️
    title: 19 个生态系统
    details: npm、PyPI、Maven、NuGet、RubyGems、Go、Cargo、Hex、Pub、Packagist 等——全部定义为常量。
    link: /zh/reference/ecosystems
    linkText: 生态列表
  - icon: 🔍
    title: 过滤与查询
    details: 按生态系统、引用类型（FIX/ADVISORY）或别名模式过滤。查询 CVSS、Maven GAV、版本范围与事件时间线。
    link: /zh/guide/cli
  - icon: 🚀
    title: 自动化发布
    details: 打 tag 即由 goreleaser 构建全平台二进制。Release 附带二进制 + 校验和——可供流水线与智能体校验。
    link: https://github.com/scagogogo/osv-schema-skills/releases
    linkText: 下载
---

### 🤖 一段提示词，让你的 AI 智能体上岗

最快路径：**复制下方提示词，粘贴进 Claude Code 或 Codex，回车。** 智能体会安装 CLI、发现技能，随即开始处理 OSV 漏洞数据。完整版见 [AI Agent 接入](/zh/guide/ai-agent) 页。

```text
You now have access to the OSV Schema Skills toolkit
(https://github.com/scagogogo/osv-schema-skills), an AI-native Go library + CLI + Claude Code
Skills bundle for the OSV vulnerability schema. Set it up now:
1. Install the `osv` CLI — download a pre-built binary from the GitHub Release matching my
   OS/arch, or `go install github.com/scagogogo/osv-schema-skills/cmd/osv@latest`. Verify `osv version`.
2. Commands: `osv parse [-v] <file>`, `osv validate <file>…`, `osv filter -e/-r/-a <file>`,
   `osv query --severity cvss3|cvss2 --maven --ranges --events <file>`. Use `-o json` for parsing.
3. Clone the repo to activate the 6 Claude Code Skills (osv-parse/validate/filter/query/severity/affected).
When I ask about a vulnerability, pick the right command automatically, filter by ecosystem if I
name one, extract CVSS + affected ranges, and report concisely. Don't ask me which command to run.
```

→ [获取完整可复制提示词 →](/zh/guide/ai-agent)

---

### 它解决了什么问题

无论对人还是对 AI，处理漏洞数据都很痛苦：

- **OSV JSON 嵌套极深**——受影响包、CVSS 分数、版本范围、引用、事件时间线。手动看一条记录得翻 500 行。
- **过滤每次都要写一次性脚本**——"只看 PyPI 的包"或"只看 FIX 引用"，每次都变成一个临时脚本。
- **按 schema 校验**没有工具很容易出错（缺 `id`、范围格式错、severity 类型错）。
- **AI 智能体过去没有结构化入口**——在这个项目之前，智能体只能 `cat` JSON 然后一路瞎编。

### 解决方案：一个内核，三层访问

同一套解析/过滤/查询逻辑，可从任何地方触达：

```mermaid
graph TD
  A["🤖 AI Agent 技能<br/>6 个自动触发技能"]
  C["🖥️ CLI<br/>osv parse/validate/filter/query"]
  S["📦 Go SDK<br/>OsvSchema 泛型"]
  CORE["Go 内核库<br/>parse · validate · filter · query"]
  OSV["OSV Schema<br/>CVE · GHSA · CVSS · affected · ranges"]

  A --> CORE
  C --> CORE
  S --> CORE
  CORE --> OSV
```

| 层 | 最适合 | 示例 |
|----|--------|------|
| 🤖 **技能** | Claude Code、AI 工作流 | 你一提漏洞文件，智能体自动触发 `osv-parse` |
| 🖥️ **CLI** | Shell、CI 流水线 | `osv filter -e PyPI -o json vuln.json` |
| 📦 **SDK** | Go 应用 | `v.Affected.FilterByEcosystem(osv.EcosystemPyPI)` |

### 工作原理——意图到技能的路由

关键在于 **intent-to-skill routing（意图→技能路由）**：每个 `SKILL.md` 声明 *何时* 触发（`description`）以及 *能调用什么工具*（`allowed-tools: Bash(osv:*)`）。智能体把你的请求与这些描述匹配，挑出正确的 `osv` 子命令——你永远不必点名要跑哪个。

```mermaid
flowchart TD
  U["用户：GHSA-... 严重吗？<br/>影响 PyPI 吗？"] --> MATCH{"将意图匹配到<br/>技能描述"}
  MATCH -->|解析| P["osv-parse"]
  MATCH -->|按生态过滤| F["osv-filter"]
  MATCH -->|严重程度| SEV["osv-query --severity"]
  P --> CORE["Go 内核<br/>UnmarshalFromJson"]
  F --> CORE
  SEV --> CORE
  CORE --> R["结构化结果 → 智能体总结"]
```

底层上，每条命令都调用同一个带类型的 Go 内核（`OsvSchema[EcosystemSpecific, DatabaseSpecific any]`）——所以 CLI、SDK、技能三者绝不可能彼此不一致。技能只是薄薄的 **声明式契约**；所有真实逻辑都集中在一处。

### 数据如何流转：从 JSON 到报告

```mermaid
flowchart LR
  F["OSV JSON 文件"] --> U["UnmarshalFromJson"]
  U --> V["OsvSchema 结构体<br/>（带类型）"]
  V --> OP{"操作"}
  OP -->|显示| P["parse 输出关键字段"]
  OP -->|校验| VAL["validate 校验 id/版本"]
  OP -->|过滤| FLT["filter 按 -e/-r/-a"]
  OP -->|提取| QRY["query 取 CVSS/范围/事件"]
  P --> R["人类可读 / JSON"]
  VAL --> R
  FLT --> R
  QRY --> R
```

### 一个典型的 AI 智能体工作流

```mermaid
flowchart LR
  U["用户提到<br/>一个漏洞"] --> P["osv-parse<br/>解析 JSON"]
  P --> F["osv-filter<br/>按生态过滤"]
  F --> SEV["osv-severity<br/>CVSS 分数"]
  SEV --> R["汇报结果"]
```

---

准备好接入你的智能体了吗？**[复制提示词 →](/zh/guide/ai-agent)**
