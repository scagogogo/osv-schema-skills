# Research: osv-schema 仓库概览

**Question:** 当前仓库的完整结构、功能、技术栈和代码现状是什么？
**Context:** 用户需要了解仓库全貌，以便后续开发决策
**Deliverable:** 结构化的仓库概览报告
**Scope:** Small

---

## 1. 项目基本信息

| 属性 | 值 |
|------|------|
| **模块名** | `github.com/scagogogo/osv-schema` |
| **语言** | Go (最低版本 1.18) |
| **包名** | `osv_schema` |
| **许可证** | LICENSE 文件（需查看具体类型） |
| **规范标准** | [OSV Schema](https://ossf.github.io/osv-schema/) — 开源漏洞描述标准格式 |
| **定位** | Go 语言实现的 OSV Schema，用于解析、操作漏洞数据 |

## 2. 核心设计理念

项目使用 **Go 泛型** 来处理 OSV Schema 中两个可扩展的字段：

- `EcosystemSpecific` — 由包管理器生态决定的特定数据（如 npm/maven 特有的字段）
- `DatabaseSpecific` — 由漏洞数据库决定的特定数据（如 OSV.dev、GitHub Advisory 特有的字段）

核心结构体 `OsvSchema[EcosystemSpecific, DatabaseSpecific any]` 让使用者可以灵活定制这两个维度。

## 3. 文件结构与职责

| 文件 | 行数 | 核心类型 | 职责 |
|------|------|---------|------|
| `osv_schema.go` | 48 | `OsvSchema[EcosystemSpecific, DatabaseSpecific]` | 主结构体，包含 OSV 所有顶级字段 |
| `affected.go` | 145 | `Affected[E,S,D]`, `AffectedSlice[E,D]` | 漏洞影响范围：包名、版本区间、严重级别 |
| `package.go` | 159 | `Package`, `Ecosystem` (18种常量) | 包管理器类型和包信息，含 Maven 的 GroupID/ArtifactID 解析 |
| `range.go` | 86 | `Range[DatabaseSpecific]`, `RangeType` | 版本范围：SEMVER/ECOSYSTEM/GIT 三种类型 |
| `event.go` | 76 | `Event`, `Events` | 版本事件：introduced/fixed/last_affected/limit |
| `severity.go` | 154 | `Severity`, `SeveritySlice`, `SeverityType` | CVSS 严重级别，含分数解析和数据库序列化 |
| `aliases.go` | 65 | `Aliases` | 漏洞别名编号，含 GetCVE 快捷方法 |
| `credits.go` | 69 | `Credits`, `CreditsType` (10种常量) | 漏洞致谢信息 |
| `references.go` | 147 | `Reference`, `References`, `ReferenceType` (10种) | 参考链接，含按类型过滤 |
| `related.go` | 41 | `Related` | 相关漏洞编号 |
| `errors.go` | 12 | `wrapScanError` | Scan 错误包装辅助函数 |
| `unmarshal.go` | 26 | `UnmarshalFromJson`, `UnmarshalFromJsonFile` | JSON 反序列化入口函数 |
| `unmarshal_test.go` | 21 | 1个测试 | 基础 JSON 解析测试 |

## 4. 序列化支持矩阵

每个核心类型都实现了 `sql.Scanner` + `driver.Valuer`，支持数据库存储：

| 格式 | 标签 | 用途 |
|------|------|------|
| `json` | `json:"..."` | JSON 序列化/反序列化 |
| `yaml` | `yaml:"..."` | YAML 支持 |
| `mapstructure` | `mapstructure:"..."` | 配置文件解析 (如 Viper) |
| `db` | `db:"..."` | 通用数据库标签 |
| `bson` | `bson:"..."` | MongoDB 支持 |
| `gorm` | `gorm:"column:...;serializer:json"` | GORM ORM (复杂字段用 JSON serializer) |

**数据库策略：** 简单字段直接存为列，复杂嵌套结构（如 `AffectedSlice`、`SeveritySlice`）存为 JSON 字符串。

## 5. 便捷方法

| 类型 | 方法 | 作用 |
|------|------|------|
| `OsvSchema` | `AffectedHasEcosystem` | 检查是否影响某生态 |
| `AffectedSlice` | `HasEcosystem`, `Filter`, `FilterByEcosystem` | 按生态过滤受影响包 |
| `Aliases` | `GetCVE`, `Filter` | 快捷获取 CVE 编号 |
| `SeveritySlice` | `GetCVSS3`, `GetCVSS2` | 快捷获取 CVSS 分数 |
| `Severity` | `GetScore`, `GetScoreAsFloat`, `GetScoreAsPointer` | 解析 CVSS 分数字符串为 float64 |
| `References` | `FilterByType` | 按引用类型过滤 |
| `Package` | `IsMaven`, `GetGroupID`, `GetArtifactID` | Maven 包名解析 |
| `Event` | `IsIntroduced`, `IsFixed`, `IsLastAffected`, `IsLimit` | 事件类型判断 |

## 6. 依赖

| 依赖 | 版本 | 用途 |
|------|------|------|
| `go-pointer` | v0.0.2 | 指针工具（`ToPointer`, `FromPointer`） |
| `testify` | v1.8.3 | 测试断言 |

## 7. 测试现状

- **测试文件：** `unmarshal_test.go`（1 个测试）
- **测试数据：** `test_data/GHSA-vxv8-r8q2-63xw.json`
- **覆盖率：** 极低，仅覆盖 JSON 反序列化基础路径
- **未覆盖：** Scan/Value、CVSS 解析、过滤方法、Maven 解析等

## 8. 已注释掉的代码

`osv_schema.go:50-69` — `OsvSchema` 的 `sql.Scanner`/`driver.Valuer` 实现已被注释掉，说明顶级结构暂不支持直接数据库 Scan/Value。

## 9. 代码风格特征

- 注释混合中英文
- 有 TODO 注释（`osv_schema.go:21` — 对 Withdrawn 字段的疑问）
- 分隔线风格：`// ------------------------------------------------- --------------------------------------------------------------------`
- 所有类型都用泛型参数实现可扩展性
- 常量定义完整，符合 OSV 规范

## 10. 关键发现与行动建议

1. **测试覆盖率极低** — 仅 1 个基础测试，大量方法（Scan/Value、GetCVE、GetScore 等）未测试
2. **Withdrawn 字段类型不一致** — OSV 规范要求 `time.Time`，但当前是 `string`（代码中有 TODO 疑问）
3. **OsvSchema 未实现 Scanner/Valuer** — 顶级结构无法直接存入数据库
4. **Ecosystem 常量未覆盖所有 OSV 规范** — OSV 规范还定义了 `SUSE`, `Photon`, `Windows` 等生态
5. **Event 结构体与 OSV 规范不完全匹配** — OSV 规范中 Event 是 "oneof" 语义（只有其中一个字段有值），当前是扁平结构体

---

**信息源：**
- 源码文件全量阅读
- OSV Schema 规范 (https://ossf.github.io/osv-schema/)
- Git 历史 (20 commits)