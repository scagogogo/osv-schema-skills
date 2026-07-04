# Changelog

All notable releases are documented here. For the full commit history, see [GitHub Releases](https://github.com/scagogogo/osv-schema-skills/releases).

---

## v0.1.0 — 2026-07-05

**First release with full publishing capability.**

This marks the transition from a minimal prototype (v0.0.1, 2023-07) to a production-ready toolkit: 7 Claude Code Skills, full CLI, VitePress documentation site, and goreleaser multi-platform binaries.

### Added

- **7 Claude Code Skills** — 6 data skills (`osv-parse`, `osv-validate`, `osv-filter`, `osv-query`, `osv-severity`, `osv-affected`) plus `osv-installation` setup guide, all auto-triggering on vulnerability-related requests
- **Full CLI** — `osv parse`, `osv validate`, `osv filter`, `osv query`, `osv version` with `-o json` output and clear exit codes (0/1/2+)
- **VitePress website** — bilingual (EN + ZH), deployed to GitHub Pages, with Mermaid diagrams throughout ("one picture beats a thousand words")
- **goreleaser multi-platform binaries** — pre-built archives for:
  - Linux: `amd64`, `arm64`, `arm` (v7)
  - macOS: `amd64`, `arm64` (Apple Silicon)
  - Windows: `amd64`, `arm64`
  - Plus `checksums.txt` for integrity verification
- **19 ecosystems** — constants for `npm`, `PyPI`, `Maven`, `NuGet`, `RubyGems`, `Go`, `Cargo`, `Hex`, `Pub`, `Packagist`, and more
- **Type-safe Go SDK** — generic `OsvSchema[EcosystemSpecific, DatabaseSpecific]` with JSON, YAML, GORM, BSON serialization

### Changed

- README refactored for **AI First** — first screen is executable Quick Start, not narrative prose
- Documentation placeholders `<latest-tag>` now point to actual releases

### Fixed

- VitePress `config.ts` double-nested `themeConfig` — navigation and sidebar now render correctly
- Mermaid syntax validation — diagrams use only production-proven constructs

---

## v0.0.1 — 2023-07-xx

Initial prototype. No pre-built binaries; `go install` only. Minimal README. No Skills, no website, no release pipeline.