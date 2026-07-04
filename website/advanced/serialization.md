# GORM & BSON Serialization

Every core type in the SDK carries serialization tags for JSON, YAML, mapstructure, GORM, and BSON. This page shows how to use them for persistence.

---

## Serialization tags

Each field has multiple tags:

```go
type OsvSchema[EcosystemSpecific, DatabaseSpecific any] struct {
    Id          string `json:"id" yaml:"id" mapstructure:"id" db:"id" bson:"id" gorm:"column:id"`
    Summary     string `json:"summary,omitempty" yaml:"summary,omitempty" mapstructure:"summary" db:"summary" bson:"summary" gorm:"column:summary"`
    Affected    AffectedSlice[EcosystemSpecific, DatabaseSpecific] `json:"affected" yaml:"affected" mapstructure:"affected" db:"affected" bson:"affected" gorm:"column:affected;serializer:json"`
    ...
}
```

| Tag | Use case |
|-----|----------|
| `json` | API responses, file reading |
| `yaml` | YAML config files |
| `mapstructure` | Viper/config parsing |
| `db` / `gorm` | SQL databases via GORM |
| `bson` | MongoDB |

---

## GORM: simple fields are columns

Fields like `id`, `schema_version`, `summary` are stored as regular SQL columns:

```go
import "gorm.io/gorm"

db, _ := gorm.Open(postgres.Open("dsn"), &gorm.Config{})
db.AutoMigrate(&osv_schema.OsvSchema[any, any]{})

v, _ := osv_schema.UnmarshalFromJsonFile[any, any]("vuln.json")
db.Create(&v)

// Query by id
var loaded osv_schema.OsvSchema[any, any]
db.First(&loaded, "id = ?", "GHSA-vxv8-r8q2-63xw")
```

---

## GORM: nested structures are JSON columns

Complex nested structures (`AffectedSlice`, `SeveritySlice`, `ranges[]`) are stored as JSON strings via `serializer:json`:

```go
// Affected is stored as a JSON column
Affected AffectedSlice[Eco, DB] `gorm:"column:affected;serializer:json"`
```

This means you can't query individual affected packages with raw SQL — you'd need to use PostgreSQL's JSON operators or a different schema design. But for most use cases (insert + read by id), it works fine.

---

## MongoDB via BSON

The `bson` tags align with MongoDB document fields:

```go
import "go.mongodb.org/mongo-driver/mongo"

coll := db.Collection("vulnerabilities")
v, _ := osv_schema.UnmarshalFromJsonFile[any, any]("vuln.json")

// Insert
coll.InsertOne(ctx, v)

// Query by id
var loaded osv_schema.OsvSchema[any, any]
coll.FindOne(ctx, bson.M{"id": "GHSA-vxv8-r8q2-63xw"}).Decode(&loaded)
```

---

## Custom types + serialization

If you define custom `EcosystemSpecific` and `DatabaseSpecific` structs, their fields **must** also carry serialization tags:

```go
type PyPISpecific struct {
    AffectedArchitectures []string `json:"affected_architectures" bson:"affected_architectures" gorm:"column:affected_architectures;serializer:json"`
}
```

Otherwise they won't survive the round-trip through GORM/BSON.

---

## CLI: `-o json` output

The CLI's `-o json` flag re-marshals through the same tagged struct:

```bash
osv parse -o json vuln.json | jq '.affected'
```

Field names match the OSV schema exactly because the `json` tags use the same keys.

---

## See also

- [Custom specifics](/advanced/custom-specifics) — how to define your own types
- [SDK guide](/guide/sdk) — basic SDK usage
- [OSV Schema reference](/reference/osv-schema) — field definitions