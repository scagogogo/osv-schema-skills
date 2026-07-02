# 技能总览

本仓库被设计为一个 **技能仓库**。当 Claude Code 打开它时，7 个专用技能自动可用——无需集成代码。

## 技能

六个**数据技能**干真正的活——解析、校验、过滤、查询、评分、影响面。第七个 `osv-installation` 是**安装向导**，首次使用时触发；它不调用 `osv`（其 `allowed-tools` 是 `Bash(go:*)`，而非 `Bash(osv:*)`）。

| 技能 | 用途 | 何时自动触发 |
|------|------|--------------|
| [`osv-parse`](/zh/guide/skills/parse) | 解析并展示 OSV JSON 数据 | 你提到解析漏洞文件或提取 CVE/GHSA 数据 |
| [`osv-validate`](/zh/guide/skills/validate) | 校验 OSV JSON 文件 | 你要检查 schema 合规性或验证漏洞文件 |
| [`osv-filter`](/zh/guide/skills/filter) | 按生态 / 引用类型 / 别名过滤 | 你想要 npm/PyPI/Maven 过滤或 FIX 引用 |
| [`osv-query`](/zh/guide/skills/query) | 提取 severity、Maven、ranges、events | 你需要 CVSS 分数、Maven GAV 或版本范围 |
| [`osv-severity`](/zh/guide/skills/severity) | CVSS 严重程度分析 | 你在评估漏洞风险或严重程度 |
| [`osv-affected`](/zh/guide/skills/affected) | 受影响包与版本分析 | 你需要影响分析或版本范围检查 |
| [`osv-installation`](/zh/guide/installation) | 安装与设置指南 | 你是第一次使用这些技能 |

## 什么措辞落到哪个技能

你从不点名技能——你只描述意图，智能体拿你的话去匹配每个技能的 `description`。以下是路由到各技能的典型措辞：

```mermaid
mindmap
  root((你的话))
    osv-parse
      "读一下这条公告"
      "这个文件里有什么"
      "把 CVE 提出来"
    osv-validate
      "这是合规的 OSV 吗"
      "检查 schema"
      "校验这些文件"
    osv-filter
      "只看 npm 的"
      "PyPI 包"
      "只要 FIX 链接"
    osv-query
      "CVSS 分数"
      "Maven groupId"
      "版本范围"
    osv-severity
      "有多严重"
      "风险等级"
      "算不算严重"
    osv-affected
      "哪些版本"
      "我受影响吗"
      "受影响的包"
    osv-installation
      "安装 CLI"
      "设置 osv"
      "第一次用"
```

## 技能如何接线

每个技能是 `.claude/skills/<name>/` 下的一个 `SKILL.md` 文件：

```mermaid
graph LR
  FM["YAML frontmatter<br/>name · description · allowed-tools"] --> Body["结构化正文<br/>决策树 · 模式 · API 参考"]
  Body --> CLI["允许的工具<br/>Bash(osv:*)"]
  CLI --> Core["Go 内核"]
```

1. **YAML frontmatter**——告诉智能体 *何时* 触发以及 *能用什么工具*。
2. **结构化正文**——决策树、任务模式、API 参考、代码示例。

示例——`osv-parse` frontmatter：

```yaml
---
name: osv-parse
description: Parse an OSV JSON file and display structured vulnerability data.
             Triggers on mentions of OSV parsing, CVE/GHSA data extraction...
allowed-tools: "Bash(osv:*)"
argument-hint: <osv-json-file>
---
```

## 技能决策树

当智能体遇到一个漏洞任务，它经由技能路由：

```mermaid
flowchart TD
  Start["用户提到<br/>漏洞数据"] --> Q{"需要什么？"}
  Q -->|"解析 / 读取"| P["osv-parse"]
  Q -->|"检查 schema"| V["osv-validate"]
  Q -->|"缩小范围"| F["osv-filter"]
  Q -->|"提取字段"| Q2["osv-query"]
  Q -->|"风险等级"| SEV["osv-severity"]
  Q -->|"影响 / 版本"| AFF["osv-affected"]

  P --> Need{"还需要更多？"}
  V --> Need
  F --> Need
  Q2 --> Need
  SEV --> Need
  AFF --> Need
  Need -->|"是"| Q
  Need -->|"否"| Report["汇报结果"]
```

## 技能之间的能力边界

```mermaid
graph TD
  subgraph 只读["技能都是只读契约"]
    P["osv-parse<br/>展示"]
    V["osv-validate<br/>校验"]
    F["osv-filter<br/>缩小"]
    Q["osv-query<br/>提取"]
    SEV["osv-severity<br/>评分"]
    AFF["osv-affected<br/>影响面"]
  end
  P -.调用.-> CLI["osv CLI"]
  V -.调用.-> CLI
  F -.调用.-> CLI
  Q -.调用.-> CLI
  SEV -.调用.-> CLI
  AFF -.调用.-> CLI
  CLI --> CORE["共享 Go 内核<br/>真实逻辑所在"]
```

技能本身不含逻辑——它们只声明 *何时触发* 和 *调哪个命令*。

## 真实工作流

```
用户："检查 GHSA-vxv8-r8q2-63xw 是否影响任何 PyPI 包，有多严重"

智能体工作流：
1. → osv-parse:     解析 OSV JSON 文件
2. → osv-filter:    按 PyPI 生态过滤受影响包
3. → osv-severity:  提取 CVSS v3 分数
4. → 向用户汇报结果
```

```mermaid
sequenceDiagram
  participant U as 用户
  participant A as 智能体
  participant S as 技能链
  U->>A: GHSA-vxv8 影响哪些 PyPI 包？多严重？
  A->>S: 1) osv-parse 解析
  S-->>A: 结构化数据
  A->>S: 2) osv-filter -e PyPI
  S-->>A: 过滤后的 affected
  A->>S: 3) osv-severity 取 CVSS
  S-->>A: 7.5 (高)
  A-->>U: 影响 3 个 PyPI 包，CVSS v3=7.5 高危
```

## 在你的项目里使用技能

**方案一——克隆本仓库。** Claude Code 打开目录时技能自动激活：

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills
```

**方案二——作为 Claude Code 插件安装**（即将推出）：

```bash
claude plugin add scagogogo/osv-schema-skills
```

::: tip
技能是只读契约——它们只声明 *何时触发* 和 *调哪个 CLI 命令*。所有真实逻辑都在共享 Go 内核里，所以技能、CLI、SDK 三者行为完全一致。
:::
