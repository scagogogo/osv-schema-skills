# Introduction

**OSV Schema Skills** is an **AI-native** Go library + CLI + Skills bundle for the [OSV (Open Source Vulnerability) Schema](https://ossf.github.io/osv-schema/). It lets you parse, validate, filter, and query vulnerability data — through a **Go SDK**, a **CLI tool**, or directly via **AI agent skills**.

## Architecture at a glance

![Three-layer architecture](/architecture.svg)

All three access layers share **one Go core**, so behavior is identical whether an AI agent, a shell script, or a Go program is driving.

## Why this exists

Working with vulnerability data is tedious:

- **OSV JSON** carries rich nested structures (affected packages, CVSS scores, version ranges, references) that are hard to inspect by hand.
- **Filtering** by ecosystem, severity, or reference type usually means writing throwaway code.
- **Validation** against the schema is error-prone without tooling.
- **AI agents** (like Claude Code) had no structured way to interact with vulnerability data — until now.

This repo closes that last gap: when Claude Code opens it, 6 specialized skills become automatically available.

## Design principles

| Principle | How it shows up |
|-----------|-----------------|
| AI First | Skills auto-trigger; README leads with copy-paste commands an agent can run |
| One picture beats a thousand words | This site leans on Mermaid diagrams instead of prose |
| Type safety | Generic `OsvSchema[EcosystemSpecific, DatabaseSpecific any]` |
| Broad serialization | Every core type supports JSON, YAML, mapstructure, GORM, BSON |
| Never nil from constructors | Unmarshal functions return errors explicitly |

## What's next

- 🤖 **[AI Agent 接入](/guide/ai-agent)** — copy one prompt into Claude Code / Codex, done
- [Quick Start](/guide/quick-start) — running in 30 seconds
- [Skills Overview](/guide/skills) — the 6 auto-triggering skills
- [CLI](/guide/cli) / [Go SDK](/guide/sdk)
