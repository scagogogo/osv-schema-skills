# osv-validate

校验 OSV JSON 文件是否符合 schema。

> **触发条件：** 提到 OSV 校验、漏洞格式检查、schema 合规性，或验证文件是否规范。
> **技能源码：** [`.claude/skills/osv-validate/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-validate/SKILL.md)

## CLI

```bash
osv validate vulnerability.json              # 单文件
osv validate file1.json file2.json           # 批量
osv validate -o json vulnerability.json      # JSON 输出
```

若有文件无效则以退出码 `1` 退出——对 CI 友好。

| 标志 | 说明 |
|------|------|
| `-o, --output` | `text`（默认）或 `json` |

## 它检查什么

- 文件可读且是合法 JSON
- 能作为 OSV 解析（`UnmarshalFromJson`）
- 必需字段存在：`id` 和 `schema_version`

## 校验流程

```mermaid
flowchart TD
  F["输入文件"] --> R{"可读 & 合法 JSON?"}
  R -->|"否"| E1["✗ 报错"]
  R -->|"是"| P{"能解析为 OSV?"}
  P -->|"否"| E2["✗ 报错"]
  P -->|"是"| C{"id & schema_version<br/>都存在?"}
  C -->|"否"| E3["✗ 报错"]
  C -->|"是"| OK["✓ 有效"]
```

## 决策树

```mermaid
flowchart TD
  Q["有 OSV 文件吗？"] --> V["osv validate *.json"]
  V --> R{"退出码？"}
  R -->|"0"| OK["✓ 全部有效"]
  R -->|"1"| Fail["✗ 有无效——列出错误"]
  Fail --> Gate["CI 闸门失败"]
```

## 在 CI 中的位置

```mermaid
flowchart LR
  PUSH["提交 / PR"] --> CI["CI 流水线"]
  CI --> VAL["osv validate *.json"]
  VAL --> RC{"exit code"}
  RC -->|"0"| MERGE["可合并 ✓"]
  RC -->|"1"| BLK["阻止合并 ✗"]
```

## SDK 等价

```go
raw, _ := os.ReadFile("vulnerability.json")
if !json.Valid(raw) { /* 不是 JSON */ }
v, err := osv.UnmarshalFromJson[any, any](raw)
// 然后检查 v.ID != "" && v.SchemaVersion != ""
```

## 交叉引用

- [[osv-parse]] — 展示有效文件的内容
- [[osv-installation]] — 先安装 CLI
