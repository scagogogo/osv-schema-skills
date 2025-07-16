# OSV Schema

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/osv-schema.svg)](https://pkg.go.dev/github.com/scagogogo/osv-schema)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/osv-schema)](https://goreportcard.com/report/github.com/scagogogo/osv-schema)

简体中文 | [English](README.md)

一个用于处理不同包管理器和漏洞数据库中漏洞数据的 OSV（开源漏洞）Schema 的 Go 语言实现。

## 什么是 OSV Schema？

OSV Schema 是描述开源软件漏洞的标准化格式。这个 Go 语言库提供了 OSV Schema 规范的完整实现，允许您以类型安全的方式解析、操作和处理漏洞数据。

## 特性

- **完整的 OSV Schema 实现**：全面支持所有 OSV Schema 字段和结构
- **泛型类型支持**：灵活的生态系统特定和数据库特定数据处理
- **多种序列化格式**：支持 JSON、YAML 和数据库存储
- **丰富的查询方法**：内置的过滤和查询漏洞数据的方法
- **CVSS 支持**：解析和处理 CVSS v2 和 v3 评分
- **数据库集成**：内置支持带有 GORM 标签的 SQL 数据库

## 安装

```bash
go get -u github.com/scagogogo/osv-schema
```

## 快速开始

### 基本用法

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
    cve := vulnerability.Aliases.GetCVE()
    if cve != "" {
        fmt.Printf("CVE: %s\n", cve)
    }
    
    // 检查特定生态系统是否受影响
    if vulnerability.Affected.HasEcosystem("npm") {
        fmt.Println("此漏洞影响 npm 包")
    }
}
```

### 处理严重性评分

```go
// 获取 CVSS v3 评分
if cvss3 := vulnerability.Severity.GetCVSS3(); cvss3 != nil {
    score := cvss3.GetScore()
    fmt.Printf("CVSS v3 评分: %.1f\n", score)
}

// 获取 CVSS v2 评分
if cvss2 := vulnerability.Severity.GetCVSS2(); cvss2 != nil {
    score := cvss2.GetScore()
    fmt.Printf("CVSS v2 评分: %.1f\n", score)
}
```

### 过滤受影响的包

```go
// 按生态系统过滤
npmAffected := vulnerability.Affected.FilterByEcosystem("npm")
for _, affected := range npmAffected {
    fmt.Printf("包名: %s\n", affected.Package.Name)
    for _, version := range affected.Versions {
        fmt.Printf("  受影响版本: %s\n", version)
    }
}

// 自定义过滤
criticalAffected := vulnerability.Affected.Filter(func(affected *osv.Affected[any, any]) bool {
    // 自定义逻辑来判断是否为严重漏洞
    return len(affected.Versions) > 10
})
```

### 处理参考资料

```go
// 按类型过滤参考资料
advisories := vulnerability.References.FilterByType(osv.ReferenceTypeAdvisory)
for _, ref := range advisories {
    fmt.Printf("公告: %s\n", ref.URL)
}

fixes := vulnerability.References.FilterByType(osv.ReferenceTypeFix)
for _, ref := range fixes {
    fmt.Printf("修复: %s\n", ref.URL)
}
```

## 核心类型

### OsvSchema

表示漏洞的主要结构：

```go
type OsvSchema[EcosystemSpecific, DatabaseSpecific any] struct {
    SchemaVersion    string                                              `json:"schema_version"`
    ID              string                                              `json:"id"`
    Modified        time.Time                                           `json:"modified"`
    Published       time.Time                                           `json:"published"`
    Withdrawn       string                                              `json:"withdrawn"`
    Aliases         Aliases                                             `json:"aliases"`
    Related         Related                                             `json:"related"`
    Summary         string                                              `json:"summary"`
    Details         string                                              `json:"details"`
    Severity        SeveritySlice                                       `json:"severity"`
    Affected        AffectedSlice[EcosystemSpecific, DatabaseSpecific]  `json:"affected"`
    References      References                                          `json:"references"`
    DatabaseSpecific DatabaseSpecific                                   `json:"database_specific"`
    Credits         *Credits                                            `json:"credits"`
}
```

### 关键组件

- **Aliases**：漏洞标识符（CVE、GHSA 等）
- **Affected**：受漏洞影响的包和版本范围
- **Severity**：CVSS 评分和其他严重性指标
- **References**：指向公告、修复和其他资源的链接
- **Range**：使用语义版本控制的版本范围

## 生态系统支持

该库支持所有主要的包生态系统：

- npm (Node.js)
- PyPI (Python)
- Maven (Java)
- NuGet (.NET)
- RubyGems (Ruby)
- Go modules
- Cargo (Rust)
- 以及更多...

## 数据库集成

该库包含对带有 GORM 标签的数据库存储的内置支持：

```go
// 结构体可以直接与 GORM 一起使用
db.AutoMigrate(&osv.OsvSchema[any, any]{})

// 将漏洞保存到数据库
db.Create(&vulnerability)

// 查询漏洞
var vulnerabilities []osv.OsvSchema[any, any]
db.Where("summary LIKE ?", "%critical%").Find(&vulnerabilities)
```

## 测试

运行测试套件：

```bash
go test ./...
```

## 文档

- [OSV Schema 规范](https://ossf.github.io/osv-schema/)
- [Go 包文档](https://pkg.go.dev/github.com/scagogogo/osv-schema)

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

本项目根据 LICENSE 文件中指定的条款进行许可。
