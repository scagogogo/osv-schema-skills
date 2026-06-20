# OSV Schema Skills

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/osv-schema.svg)](https://pkg.go.dev/github.com/scagogogo/osv-schema)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/osv-schema)](https://goreportcard.com/report/github.com/scagogogo/osv-schema)

[简体中文](#简体中文) | **English**

## What Is This?

**OSV Schema Skills** is an **AI-native** Go library + CLI + Skills bundle for the [OSV (Open Source Vulnerability) Schema](https://ossf.github.io/osv-schema/). It lets you parse, validate, filter, and query vulnerability data — through a **Go SDK**, a **CLI tool**, or directly via **AI agent skills**.

### The Problem

Working with vulnerability data is tedious:

- **OSV JSON files** contain rich, nested structures (affected packages, CVSS scores, version ranges, references) that are hard to inspect manually
- **Filtering** by ecosystem, severity, or reference type requires writing custom code every time
- **Validating** OSV files against the schema is error-prone without tooling
- **AI agents** (like Claude Code) have no structured way to interact with vulnerability data

### The Solution

This repository provides **three layers of access**, all backed by the same Go library:

| Layer | Best For | Example |
|-------|----------|---------|
| 🤖 **AI Agent Skills** | Claude Code, AI workflows, automated analysis | Agent auto-triggers `osv-parse` when you mention a vulnerability file |
| 🖥️ **CLI** | Quick lookups, shell scripting, CI pipelines | `osv parse vulnerability.json` |
| 📦 **Go SDK** | Integration into Go applications | `osv.UnmarshalFromJsonFile[any, any]("vuln.json")` |

---

## 🤖 AI Agent Integration

This repository is designed as a **Skills repository** — AI agents can directly plug into it without any custom integration code. When Claude Code opens this repository, 6 specialized skills become automatically available:

| Skill | Purpose | Auto-triggers when... |
|-------|---------|----------------------|
| `osv-parse` | Parse & display OSV JSON data | You mention parsing a vulnerability file or extracting CVE/GHSA data |
| `osv-validate` | Validate OSV JSON files | You ask to check schema compliance or verify a vulnerability file |
| `osv-filter` | Filter by ecosystem / reference type / alias | You want to filter by npm/PyPI/Maven ecosystem or find FIX references |
| `osv-query` | Extract severity, Maven, ranges, events | You need CVSS scores, Maven groupId/artifactId, or version ranges |
| `osv-severity` | CVSS severity analysis | You're assessing vulnerability risk or severity |
| `osv-affected` | Affected package & version analysis | You need impact analysis or version range inspection |
| `osv-installation` | Setup & installation guide | It's your first time using the skills |

### How Skills Work

Each skill is a `SKILL.md` file in `.claude/skills/<name>/` with:

1. **YAML frontmatter** — tells the AI agent *when* to trigger and *what tools* it can use
2. **Structured body** — decision trees, task patterns, API reference, code examples

Example — the `osv-parse` skill frontmatter:

```yaml
---
name: osv-parse
description: Parse an OSV JSON file and display structured vulnerability data.
             Triggers on mentions of OSV parsing, CVE/GHSA data extraction...
allowed-tools: "Bash(osv:*)"
argument-hint: <osv-json-file>
---
```

When an AI agent encounters a vulnerability JSON file, it automatically knows to invoke `osv parse <file>` — no prompting required.

### Using Skills in Your Project

**Option 1: Clone this repo** — Skills are automatically available when Claude Code opens the directory:

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills
# Skills are now active in Claude Code
```

**Option 2: Install as a Claude Code plugin** (coming soon):

```bash
claude plugin add scagogogo/osv-schema-skills
```

### Skill Example Workflow

Here's how an AI agent would use the skills in a real workflow:

```
User: "Check if GHSA-vxv8-r8q2-63xw affects any PyPI packages and how severe it is"

Agent workflow:
1. → osv-parse: Parse the OSV JSON file
2. → osv-filter: Filter affected packages by PyPI ecosystem
3. → osv-severity: Extract CVSS v3 score
4. → Report findings to user
```

---

## 🖥️ CLI

### Installation

```bash
go install github.com/scagogogo/osv-schema/cmd/osv@latest
```

### Commands

```bash
# Parse an OSV JSON file
osv parse vulnerability.json           # Key fields
osv parse -v vulnerability.json        # All fields (verbose)
osv parse -o json vulnerability.json   # JSON output

# Validate OSV JSON files
osv validate vulnerability.json              # Single file
osv validate file1.json file2.json           # Batch
osv validate -o json vulnerability.json      # JSON output

# Filter by ecosystem, reference type, or alias
osv filter -e PyPI vulnerability.json        # By ecosystem
osv filter -r ADVISORY vulnerability.json    # By reference type
osv filter -a CVE vulnerability.json         # By alias pattern
osv filter -e PyPI -r FIX vulnerability.json # Combine

# Query specific sub-information
osv query --severity cvss3 vulnerability.json  # CVSS v3 score
osv query --maven vulnerability.json           # Maven decomposition
osv query --ranges vulnerability.json          # Version ranges
osv query --events vulnerability.json          # Event timeline

# Show version
osv version
```

---

## 📦 Go SDK

### Installation

```bash
go get -u github.com/scagogogo/osv-schema
```

### Quick Start

```go
package main

import (
    "fmt"
    "log"

    osv "github.com/scagogogo/osv-schema"
)

func main() {
    // Parse OSV data from JSON file
    vulnerability, err := osv.UnmarshalFromJsonFile[any, any]("vulnerability.json")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s\n", vulnerability.ID)
    fmt.Printf("Summary: %s\n", vulnerability.Summary)

    // Get CVE from aliases
    if cve := vulnerability.Aliases.GetCVE(); cve != "" {
        fmt.Printf("CVE: %s\n", cve)
    }

    // Check if specific ecosystem is affected
    if vulnerability.Affected.HasEcosystem("npm") {
        fmt.Println("Affects npm packages")
    }

    // Get CVSS v3 score
    if cvss3 := vulnerability.Severity.GetCVSS3(); cvss3 != nil {
        fmt.Printf("CVSS v3: %.1f\n", cvss3.GetScore())
    }
}
```

### Key Methods

| Type | Method | Description |
|------|--------|-------------|
| `OsvSchema` | `Affected.HasEcosystem(eco)` | Check if ecosystem is affected |
| `AffectedSlice` | `FilterByEcosystem(eco)` | Filter affected packages |
| `AffectedSlice` | `Filter(fn)` | Custom filter predicate |
| `Aliases` | `GetCVE()` | Get first CVE identifier |
| `Aliases` | `Filter(fn)` | Filter aliases by pattern |
| `SeveritySlice` | `GetCVSS3()` | Get CVSS v3 severity entry |
| `SeveritySlice` | `GetCVSS2()` | Get CVSS v2 severity entry |
| `Severity` | `GetScore()` | Parse score as float64 |
| `References` | `FilterByType(t)` | Filter by reference type |
| `Package` | `IsMaven()` | Check if Maven package |
| `Package` | `GetGroupID()` | Maven groupId |
| `Package` | `GetArtifactID()` | Maven artifactId |
| `Event` | `IsIntroduced/IsFixed/...` | Event type checks |

### Serialization Support

Every core type supports JSON, YAML, mapstructure, database (GORM), and MongoDB (BSON) serialization out of the box.

### Ecosystem Support

npm, PyPI, Maven, NuGet, RubyGems, Go, Cargo, Pub, Hex, Packagist, and more — all 18+ ecosystems defined as constants.

---

## Core Types

```go
type OsvSchema[EcosystemSpecific, DatabaseSpecific any] struct {
    SchemaVersion    string
    ID               string
    Modified         time.Time
    Published        time.Time
    Withdrawn        string
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

Generic type parameters `EcosystemSpecific` and `DatabaseSpecific` let you attach custom data per ecosystem or vulnerability database. Use `any` for general-purpose parsing.

---

## Build & Test

```bash
go build ./...
go test ./...
go vet ./...
```

## Documentation

- [OSV Schema Specification](https://ossf.github.io/osv-schema/)
- [Go Package Documentation](https://pkg.go.dev/github.com/scagogogo/osv-schema)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms specified in the LICENSE file.

---

<a id="简体中文"></a>

## 简体中文

[English](#english) | **简体中文**

### 这是什么？

**OSV Schema Skills** 是一个 **AI 原生** 的 Go 语言库 + CLI + Skills 套件，用于 [OSV（开源漏洞）Schema](https://ossf.github.io/osv-schema/)。它可以让你通过 **Go SDK**、**CLI 工具** 或直接通过 **AI Agent Skills** 来解析、校验、过滤和查询漏洞数据。

### 解决的问题

处理漏洞数据通常很繁琐：

- **OSV JSON 文件**包含复杂的嵌套结构（受影响包、CVSS 评分、版本范围、参考资料），手动检查很困难
- **过滤**按生态系统、严重性或引用类型需要每次写自定义代码
- **校验** OSV 文件是否符合 Schema 规范缺乏工具支持
- **AI Agent**（如 Claude Code）没有结构化的方式与漏洞数据交互

### 解决方案

本仓库提供 **三层访问方式**，都基于同一个 Go 库：

| 层级 | 适用场景 | 示例 |
|------|----------|------|
| 🤖 **AI Agent Skills** | Claude Code、AI 工作流、自动化分析 | Agent 在你提到漏洞文件时自动触发 `osv-parse` |
| 🖥️ **CLI** | 快速查询、Shell 脚本、CI 流水线 | `osv parse vulnerability.json` |
| 📦 **Go SDK** | 集成到 Go 应用程序中 | `osv.UnmarshalFromJsonFile[any, any]("vuln.json")` |

### 🤖 AI Agent 接入

本仓库被设计为 **Skills 仓库** — AI Agent 可以直接对接，无需任何自定义集成代码。当 Claude Code 打开此仓库时，6 个专用 Skills 自动可用：

| Skill | 用途 | 自动触发条件 |
|-------|------|-------------|
| `osv-parse` | 解析 & 展示 OSV JSON 数据 | 提到解析漏洞文件或提取 CVE/GHSA 数据 |
| `osv-validate` | 校验 OSV JSON 文件 | 要求检查 Schema 合规性或校验漏洞文件 |
| `osv-filter` | 按生态/引用类型/别名过滤 | 需要按 npm/PyPI/Maven 过滤或查找 FIX 引用 |
| `osv-query` | 提取 severity/maven/ranges/events | 需要 CVSS 评分、Maven groupId/artifactId 或版本范围 |
| `osv-severity` | CVSS 严重级别分析 | 评估漏洞风险或严重性 |
| `osv-affected` | 受影响包 & 版本分析 | 需要影响分析或版本范围检查 |

### Skills 工作原理

每个 Skill 是 `.claude/skills/<name>/SKILL.md` 文件，包含：

1. **YAML 前置数据** — 告诉 AI Agent *何时*触发、*能用什么工具*
2. **结构化正文** — 决策树、任务模式、API 参考、代码示例

### 在你的项目中使用 Skills

**方式一：克隆仓库** — Claude Code 打开目录后 Skills 自动生效：

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills
# Skills 在 Claude Code 中已激活
```

**方式二：安装为 Claude Code 插件**（即将推出）：

```bash
claude plugin add scagogogo/osv-schema-skills
```

### Skill 实际工作流程

```
用户: "检查 GHSA-vxv8-r8q2-63xw 是否影响 PyPI 包，严重程度如何"

Agent 工作流:
1. → osv-parse:  解析 OSV JSON 文件
2. → osv-filter:  按 PyPI 生态系统过滤受影响包
3. → osv-severity: 提取 CVSS v3 评分
4. → 向用户报告发现
```

### 🖥️ CLI 使用

```bash
go install github.com/scagogogo/osv-schema/cmd/osv@latest

osv parse vulnerability.json           # 解析关键字段
osv parse -v vulnerability.json        # 完整详情
osv validate vulnerability.json        # 校验
osv filter -e PyPI vulnerability.json  # 按生态过滤
osv query --severity cvss3 vulnerability.json  # 查询 CVSS
osv query --ranges vulnerability.json  # 查询版本范围
```

### 📦 Go SDK 使用

```go
import osv "github.com/scagogogo/osv-schema"

vulnerability, _ := osv.UnmarshalFromJsonFile[any, any]("vuln.json")
fmt.Println(vulnerability.ID)
fmt.Println(vulnerability.Aliases.GetCVE())
fmt.Println(vulnerability.Affected.HasEcosystem("npm"))
```

### 贡献

欢迎贡献！请随时提交 Pull Request。

### 许可证

本项目根据 LICENSE 文件中指定的条款进行许可。
