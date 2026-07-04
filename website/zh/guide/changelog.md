# 更新日志

所有重要版本的更新记录于此。完整 commit 历史 见 [GitHub Releases](https://github.com/scagogogo/osv-schema-skills/releases)。

---

## v0.1.0 — 2026-07-05

**首个具备完整发布能力的版本。**

这次发布标志着从最小原型（v0.0.1，2023-07）到生产可用工具包的跨越：7 个 Claude Code Skills、完整 CLI、VitePress 文档站、goreleaser 多平台二进制。

### 新增

- **7 个 Claude Code Skills** — 6 个数据技能（`osv-parse`、`osv-validate`、`osv-filter`、`osv-query`、`osv-severity`、`osv-affected`）外加 `osv-installation` 安装指南，遇到漏洞相关请求时自动触发
- **完整 CLI** — `osv parse`、`osv validate`、`osv filter`、`osv query`、`osv version`，支持 `-o json` 输出与清晰的退出码（0/1/2+）
- **VitePress 官网** — 中英双语，部署到 GitHub Pages，全程使用 Mermaid 图表（"一图抵千言")
- **goreleaser 多平台二进制** — 预编译 archive 覆盖：
  - Linux：`amd64`、`arm64`、`arm`（v7）
  - macOS：`amd64`、`arm64`（Apple Silicon）
  - Windows：`amd64`、`arm64`
  - 附带 `checksums.txt` 供完整性校验
- **19 个生态系统常量** — `npm`、`PyPI`、`Maven`、`NuGet`、`RubyGems`、`Go`、`Cargo`、`Hex`、`Pub`、`Packagist` 等
- **类型安全的 Go SDK** — 泛型 `OsvSchema[EcosystemSpecific, DatabaseSpecific]`，支持 JSON、YAML、GORM、BSON 序列化

### 变更

- README 重构为 **AI First** — 首屏是可执行的 Quick Start，而非叙事段落
- 文档占位符 `<latest-tag>` 指向实际发布的版本

### 修复

- VitePress `config.ts` 双重嵌套 `themeConfig` — 导航与侧边栏现可正确渲染
- Mermaid 语法验证 — 图表仅使用经生产验证的语法元素

---

## v0.0.1 — 2023-07-xx

初始原型。无预编译二进制；仅支持 `go install`。README 极简。无 Skills、无官网、无发布流水线。