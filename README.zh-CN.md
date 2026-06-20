# OSV Schema Skills

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/osv-schema.svg)](https://pkg.go.dev/github.com/scagogogo/osv-schema)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/osv-schema)](https://goreportcard.com/report/github.com/scagogogo/osv-schema)

**简体中文** | [English](README.md)

## 这是什么？

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

---

## 🤖 AI Agent 接入

本仓库被设计为 **Skills 仓库** — AI Agent 可以直接对接，无需任何自定义集成代码。当 Claude Code 打开此仓库时，6 个专用 Skills 自动可用：

| Skill | 用途 | 自动触发条件 |
|-------|------|-------------|
| `osv-parse` | 解析 & 展示 OSV JSON 数据 | 提到解析漏洞文件或提取 CVE/GHSA 数据 |
| `osv-validate` | 校验 OSV JSON 文件 | 要求检查 Schema 合规性或校验漏洞文件 |
| `osv-filter` | 按生态/引用类型/别名过滤 | 需要按 npm/PyPI/Maven 过滤或查找 FIX 引用 |
| `osv-query` | 提取 severity/maven/ranges/events | 需要 CVSS 评分、Maven groupId/artifactId 或版本范围 |
| `osv-severity` | CVSS 严重级别分析 | 评估漏洞风险或严重性 |
| `osv-affected` | 受影响包 & 版本分析 | 需要影响分析或版本范围检查 |
| `osv-installation` | 安装 & 设置指南 | 首次使用 Skills |

### Skills 工作原理

每个 Skill 是 `.claude/skills/<name>/SKILL.md` 文件，包含：

1. **YAML 前置数据** — 告诉 AI Agent *何时*触发、*能用什么工具*
2. **结构化正文** — 决策树、任务模式、API 参考、代码示例

示例 — `osv-parse` Skill 的前置数据：

```yaml
---
name: osv-parse
description: 解析 OSV JSON 文件并展示结构化漏洞数据。
             在提到 OSV 解析、CVE/GHSA 数据提取时触发...
allowed-tools: "Bash(osv:*)"
argument-hint: <osv-json-file>
---
```

当 AI Agent 遇到漏洞 JSON 文件时，它会自动调用 `osv parse <file>` — 无需手动提示。

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
1. → osv-parse:    解析 OSV JSON 文件
2. → osv-filter:   按 PyPI 生态系统过滤受影响包
3. → osv-severity: 提取 CVSS v3 评分
4. → 向用户报告发现
```

---

## 🖥️ CLI

### 安装

```bash
go install github.com/scagogogo/osv-schema/cmd/osv@latest
```

### 命令

```bash
# 解析 OSV JSON 文件
osv parse vulnerability.json           # 关键字段
osv parse -v vulnerability.json        # 完整详情（verbose）
osv parse -o json vulnerability.json   # JSON 输出

# 校验 OSV JSON 文件
osv validate vulnerability.json              # 单个文件
osv validate file1.json file2.json           # 批量
osv validate -o json vulnerability.json      # JSON 输出

# 按生态/引用类型/别名过滤
osv filter -e PyPI vulnerability.json        # 按生态
osv filter -r ADVISORY vulnerability.json    # 按引用类型
osv filter -a CVE vulnerability.json         # 按别名模式
osv filter -e PyPI -r FIX vulnerability.json # 组合

# 查询特定子信息
osv query --severity cvss3 vulnerability.json  # CVSS v3 评分
osv query --maven vulnerability.json           # Maven 分解
osv query --ranges vulnerability.json          # 版本范围
osv query --events vulnerability.json          # 事件时间线

# 显示版本
osv version
```

---

## 📦 Go SDK

### 安装

```bash
go get -u github.com/scagogogo/osv-schema
```

### 快速开始

```go
package main

import (
    "fmt"
    "log"

    osv "github.com/scagogogo/osv-schema"
)

func main() {
    // 从 JSON 文件解析 OSV 数据
    vulnerability, err := osv.UnmarshalFromJsonFile[any, any]("vulnerability.json")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("漏洞 ID: %s\n", vulnerability.ID)
    fmt.Printf("摘要: %s\n", vulnerability.Summary)

    // 从别名中获取 CVE
    if cve := vulnerability.Aliases.GetCVE(); cve != "" {
        fmt.Printf("CVE: %s\n", cve)
    }

    // 检查特定生态系统是否受影响
    if vulnerability.Affected.HasEcosystem("npm") {
        fmt.Println("影响 npm 包")
    }

    // 获取 CVSS v3 评分
    if cvss3 := vulnerability.Severity.GetCVSS3(); cvss3 != nil {
        fmt.Printf("CVSS v3: %.1f\n", cvss3.GetScore())
    }
}
```

### 关键方法

| 类型 | 方法 | 描述 |
|------|------|------|
| `OsvSchema` | `Affected.HasEcosystem(eco)` | 检查生态系统是否受影响 |
| `AffectedSlice` | `FilterByEcosystem(eco)` | 按生态过滤受影响包 |
| `AffectedSlice` | `Filter(fn)` | 自定义过滤谓词 |
| `Aliases` | `GetCVE()` | 获取第一个 CVE 标识符 |
| `Aliases` | `Filter(fn)` | 按模式过滤别名 |
| `SeveritySlice` | `GetCVSS3()` | 获取 CVSS v3 严重级别 |
| `SeveritySlice` | `GetCVSS2()` | 获取 CVSS v2 严重级别 |
| `Severity` | `GetScore()` | 解析评分为 float64 |
| `References` | `FilterByType(t)` | 按引用类型过滤 |
| `Package` | `IsMaven()` | 检查是否为 Maven 包 |
| `Package` | `GetGroupID()` | Maven groupId |
| `Package` | `GetArtifactID()` | Maven artifactId |
| `Event` | `IsIntroduced/IsFixed/...` | 事件类型检查 |

### 序列化支持

所有核心类型开箱即支持 JSON、YAML、mapstructure、数据库 (GORM) 和 MongoDB (BSON) 序列化。

### 生态系统支持

npm、PyPI、Maven、NuGet、RubyGems、Go、Cargo、Pub、Hex、Packagist 等 — 18+ 个生态系统定义为常量。

---

## 核心类型

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

泛型参数 `EcosystemSpecific` 和 `DatabaseSpecific` 允许你附加自定义数据。通用解析使用 `any` 即可。

---

## 构建与测试

```bash
go build ./...
go test ./...
go vet ./...
```

## 文档

- [OSV Schema 规范](https://ossf.github.io/osv-schema/)
- [Go 包文档](https://pkg.go.dev/github.com/scagogogo/osv-schema)

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

本项目根据 LICENSE 文件中指定的条款进行许可。
