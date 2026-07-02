# osv-parse

解析 OSV JSON 文件并展示结构化的漏洞数据。

> **触发条件：** 提到 OSV 解析、漏洞 JSON 读取、CVE/GHSA 数据提取，或用户提供了 OSV JSON 文件路径。
> **技能源码：** [`.claude/skills/osv-parse/SKILL.md`](https://github.com/scagogogo/osv-schema-skills/blob/main/.claude/skills/osv-parse/SKILL.md)

## CLI

```bash
osv parse vulnerability.json           # 关键字段（文本）
osv parse -v vulnerability.json        # 全字段（日期、详情、范围、鸣谢）
osv parse -o json vulnerability.json   # JSON 输出
```

| 标志 | 说明 |
|------|------|
| `-v, --verbose` | 展示全字段 |
| `-o, --output` | `text`（默认）或 `json` |

## SDK 等价

```go
v, err := osv.UnmarshalFromJsonFile[any, any]("vulnerability.json")
fmt.Println(v.ID, v.Summary, v.Aliases.GetCVE())
```

## 决策树

```mermaid
flowchart TD
  Q["拿到 OSV JSON 文件？"] --> P["osv parse file.json"]
  P --> Need{"需要全部字段？"}
  Need -->|"否"| Done["展示关键字段"]
  Need -->|"是"| V["osv parse -v file.json"]
  V --> Done2["展示全字段"]
```

## 输出结构

```mermaid
flowchart TD
  OUT["parse 输出"] --> ID["ID"]
  OUT --> SV["schema 版本"]
  OUT --> SUM["摘要 / aliases / CVE"]
  OUT --> SEV["severity（CVSS）"]
  OUT --> AFF["受影响包"]
  OUT --> REF["引用"]
  OUT --> VV["-v 时额外:<br/>published/modified/withdrawn/<br/>related/details/各范围事件/credits"]
```

## 它打印什么

ID、schema 版本、摘要、aliases/CVE、severity、受影响包、引用。加 `-v` 还会展示 published/modified 日期、withdrawn、related、details、每范围事件和 credits。

## 底层发生了什么

`osv parse` 只是 SDK `UnmarshalFromJsonFile` 之上的一层薄壳——和你在 Go 里会写的调用一模一样。与 SDK 路径唯一的区别就是文本/JSON 的渲染。

```mermaid
sequenceDiagram
  participant U as 你 / 智能体
  participant CLI as osv parse
  participant SDK as UnmarshalFromJsonFile[any,any]
  participant R as 渲染器
  U->>CLI: osv parse file.json [-v] [-o json]
  CLI->>SDK: 读文件 → 解码 JSON
  SDK-->>CLI: *OsvSchema（带类型内核）
  CLI->>R: -o text → 关键字段（-v 则全部）
  CLI->>R: -o json → 把带类型内核重新 marshal
  R-->>U: 打印输出
```

## 文本 vs JSON：怎么选

```mermaid
flowchart TD
  WHO{"谁来读输出？"} -->|"终端前的人"| T["-o text（默认）<br/>紧凑、关键字段"]
  WHO -->|"脚本 / 智能体 / 管道"| J["-o json<br/>完整结构、机器可解析"]
  T --> TV{"需要每个字段？"}
  TV -->|是| ADDV["加 -v"]
  TV -->|否| OK["完成"]
```

`-o json` 把带类型的 `OsvSchema` 内核重新 marshal，所以输出就是标准 OSV 记录——字段名与 [OSV Schema](/zh/reference/osv-schema) 完全一致：

```bash
osv parse -o json vulnerability.json | jq '{id, summary, severity, affected}'
```

```json
{
  "id": "GHSA-vxv8-r8q2-63xw",
  "summary": "TensorFlow vulnerable to `CHECK` fail in `FractionalMaxPoolGrad`",
  "severity": [{ "type": "CVSS_V3", "score": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H" }],
  "affected": [{ "package": { "ecosystem": "PyPI", "name": "tensorflow", "purl": "" }, "ranges": [...] }]
}
```

::: tip 解析绝不修改文件
`parse` 只读。它把数据解码进带类型内核再打印——绝不写回。想在解析前先确认文件*格式正确*，用 [[osv-validate]]；格式错误的文件会让 `parse` 以非零码退出并给出解码错误。
:::

## 交叉引用

- [[osv-validate]] — 先确认文件 schema 合规
- [[osv-filter]] / [[osv-query]] — 对解析后的数据做收窄或提取
