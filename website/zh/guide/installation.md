# 安装

安装 `osv` CLI、Go SDK，并启用 Claude Code 技能。

## 环境要求

- **Go 1.18+**（用于 SDK 和源码构建）
- 仅 `go get` / `go install` / 下载二进制时需要联网

## 安装方式一图览

```mermaid
flowchart TD
  START["要装什么？"] --> Q1{"选一种"}
  Q1 -->|"预编译二进制<br/>（最快，无需 Go）"| BIN["从 Release 下载 tar.gz"]
  Q1 -->|"go install<br/>（需 Go）"| GI["go install ...@latest"]
  Q1 -->|"源码构建"| SRC["git clone + go build"]
  BIN --> VER["osv version 验证"]
  GI --> VER
  SRC --> VER
  VER --> OK["就绪 ✓"]
```

## CLI

::: tabs
== 预编译二进制

每个 tag 都通过 goreleaser 发布预编译二进制：

| 操作系统 | 架构 |
|----------|------|
| Linux | amd64、arm64、arm (v7) |
| macOS | amd64、arm64 |
| Windows | amd64、arm64 |

```bash
# Linux amd64 示例——按你的情况替换版本号/平台
VERSION=v0.1.0
curl -fsSL -o osv.tar.gz \
  https://github.com/scagogogo/osv-schema-skills/releases/download/${VERSION}/osv_${VERSION}_linux_amd64.tar.gz
tar -xzf osv.tar.gz osv
chmod +x osv && sudo mv osv /usr/local/bin/
osv version
```

用自带的 `checksums.txt` 校验完整性：

```bash
sha256sum -c checksums.txt --ignore-missing
```

Release 地址：<https://github.com/scagogogo/osv-schema-skills/releases>

== go install

```bash
go install github.com/scagogogo/osv-schema-skills/cmd/osv@latest
osv version
```

== 源码构建

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills
go build -o osv ./cmd/osv/
./osv version
```
:::

## Go SDK

```bash
go get -u github.com/scagogogo/osv-schema-skills
```

```go
import osv "github.com/scagogogo/osv-schema-skills"
```

用法见 [Go SDK 指南](/zh/guide/sdk)。

## Claude Code 技能

当 Claude Code 打开本仓库时，6 个技能自动激活——无需安装步骤：

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills
claude   # 技能已生效
```

或作为插件安装（即将推出）：

```bash
claude plugin add scagogogo/osv-schema-skills
```

见 [技能总览](/zh/guide/skills)。

## 验证

```bash
osv version                                   # CLI + schema 版本
osv parse test_data/GHSA-vxv8-r8q2-63xw.json  # 解析一条样例记录
```
