# GORM 与 BSON 序列化

SDK 中每个核心类型都携带 JSON、YAML、mapstructure、GORM、BSON 序列化标签。本页展示如何用于持久化。

---

## 序列化标签

每个字段有多重标签：

```go
type OsvSchema[EcosystemSpecific, DatabaseSpecific any] struct {
    Id          string `json:"id" yaml:"id" mapstructure:"id" db:"id" bson:"id" gorm:"column:id"`
    Summary     string `json:"summary,omitempty" yaml:"summary,omitempty" mapstructure:"summary" db:"summary" bson:"summary" gorm:"column:summary"`
    Affected    AffectedSlice[EcosystemSpecific, DatabaseSpecific] `json:"affected" yaml:"affected" mapstructure:"affected" db:"affected" bson:"affected" gorm:"column:affected;serializer:json"`
    ...
}
```

| 标签 | 用途 |
|------|------|
| `json` | API 响应、文件读取 |
| `yaml` | YAML 配置文件 |
| `mapstructure` | Viper/配置解析 |
| `db` / `gorm` | GORM 访问 SQL 数据库 |
| `bson` | MongoDB |

---

## GORM：简单字段是列

`id`、`schema_version`、`summary` 等字段存为普通 SQL 列：

```go
import "gorm.io/gorm"

db, _ := gorm.Open(postgres.Open("dsn"), &gorm.Config{})
db.AutoMigrate(&osv_schema.OsvSchema[any, any]{})

v, _ := osv_schema.UnmarshalFromJsonFile[any, any]("vuln.json")
db.Create(&v)

// 按 id 查询
var loaded osv_schema.OsvSchema[any, any]
db.First(&loaded, "id = ?", "GHSA-vxv8-r8q2-63xw")
```

---

## GORM：嵌套结构是 JSON 列

复杂嵌套结构（`AffectedSlice`、`SeveritySlice`、`ranges[]`）通过 `serializer:json` 存为 JSON 字符串：

```go
// Affected 存为 JSON 列
Affected AffectedSlice[Eco, DB] `gorm:"column:affected;serializer:json"`
```

这意味着你不能用原始 SQL 查询单个受影响包——需要用 PostgreSQL 的 JSON 操作符或改用不同 schema 设计。但对大多数用例（插入 + 按 id 读取），这样就够了。

---

## MongoDB via BSON

`bson` 标签与 MongoDB 文档字段对齐：

```go
import "go.mongodb.org/mongo-driver/mongo"

coll := db.Collection("vulnerabilities")
v, _ := osv_schema.UnmarshalFromJsonFile[any, any]("vuln.json")

// 插入
coll.InsertOne(ctx, v)

// 按 id 查询
var loaded osv_schema.OsvSchema[any, any]
coll.FindOne(ctx, bson.M{"id": "GHSA-vxv8-r8q2-63xw"}).Decode(&loaded)
```

---

## 自定义类型 + 序列化

如果你自定义了 `EcosystemSpecific` 和 `DatabaseSpecific` 结构体，它们的字段**必须**也携带序列化标签：

```go
type PyPISpecific struct {
    AffectedArchitectures []string `json:"affected_architectures" bson:"affected_architectures" gorm:"column:affected_architectures;serializer:json"`
}
```

否则它们无法在 GORM/BSON 往返中存活。

---

## CLI：`-o json` 输出

CLI 的 `-o json` 标志通过同一个带标签结构体重新 marshal：

```bash
osv parse -o json vuln.json | jq '.affected'
```

字段名与 OSV schema 完全一致，因为 `json` 标签用的是同一套键名。

---

## 另见

- [自定义特定字段](/zh/advanced/custom-specifics) —— 如何定义自己的类型
- [SDK 指南](/zh/guide/sdk) —— SDK 基本用法
- [OSV Schema 参考](/zh/reference/osv-schema) —— 字段定义