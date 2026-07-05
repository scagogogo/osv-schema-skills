# Ecosystem Naming Standard

OSV uses **ecosystem-scoped package names** to identify vulnerable components. This page explains the naming convention, why case-sensitivity matters, and which ecosystems this toolkit supports.

---

## Why ecosystems?

The name `requests` alone is ambiguous — it could be a PyPI package, a RubyGem, or an npm module. OSV requires every package to carry an `ecosystem` field, making the `ecosystem` + `name` pairing globally unambiguous:

```json
{
  "package": {
    "ecosystem": "PyPI",
    "name": "requests"
  }
}
```

This is the foundation of `osv filter -e PyPI` — without ecosystems, filtering "PyPI packages" would require guessing from naming conventions.

---

## Case-sensitivity

Ecosystem names are **case-sensitive** per the OSV spec. This toolkit matches exactly:

| Ecosystem | Correct | Wrong |
|-----------|---------|-------|
| Python | `PyPI` | ~~`pypi`~~, ~~`PyPi`~~ |
| Node.js | `npm` | ~~`NPM`~~, ~~`Npm`~~ |
| Go | `Go` | ~~`go`~~, ~~`GO`~~ |
| Java | `Maven` | ~~`maven`~~ |
| .NET | `NuGet` | ~~`nuget`~~ |
| Rust | `crates.io` | ~~`Crates.io`~~ |

The `osv filter -e` flag passes the value through **as-is** — it does not auto-normalize case. This is deliberate: the spec says case-sensitive, so the toolkit stays faithful.

::: tip Reference types and aliases ARE case-insensitive
Only `ecosystem` is case-sensitive. The `-r` (reference type) and `-a` (alias pattern) flags auto-uppercase before matching, so `fix`, `Fix`, `FIX` all work; `cve`, `Cve`, `CVE` all work.
:::

---

## Supported ecosystems

This toolkit defines constants for the ecosystems standardized by OSV:

| Constant | Value | Notes |
|----------|-------|-------|
| `EcosystemGo` | `Go` | Name is a Go module path |
| `EcosystemNpm` | `npm` | Name is an npm package name |
| `EcosystemOSSFuzz` | `OSS-Fuzz` | OSS-Fuzz project reports |
| `EcosystemPyPI` | `PyPI` | Name is a normalized PyPI package name |
| `EcosystemRubyGems` | `RubyGems` | Name is a gem name |
| `EcosystemCratesIo` | `crates.io` | Name is a Rust crate name |
| `EcosystemPackagist` | `Packagist` | PHP package manager |
| `EcosystemMaven` | `Maven` | Name is `groupId:artifactId` |
| `EcosystemNuGet` | `NuGet` | .NET package name |
| `EcosystemLinux` | `Linux` | Only `name: Kernel` |
| `EcosystemDebian` | `Debian` | Source package name; optional `:<RELEASE>` suffix (e.g. `Debian:7`) |
| `EcosystemAlpine` | `Alpine` | Source package name; requires `:v<RELEASE>` suffix (e.g. `Alpine:v3.16`) |
| `EcosystemHex` | `Hex` | Erlang/Elixir package |
| `EcosystemAndroid` | `Android` | Android component name (Framework, Media Framework, Kernel, …) |
| `EcosystemGitHubActions` | `GitHub Actions` | Name is `{owner}/{repo}` of the action |
| `EcosystemPub` | `Pub` | Dart package name |
| `EcosystemConanCenter` | `ConanCenter` | C/C++ Conan package name |
| `EcosystemRocky` | `Rocky` | Source package name; optional `:<RELEASE>` suffix |
| `EcosystemAlmaLinux` | `AlmaLinux` | Source package name; optional `:<RELEASE>` suffix |

See the [Ecosystems reference](/reference/ecosystems) for the full list with details.

---

## Distribution release scoping

Some Linux-distribution ecosystems carry a `:<RELEASE>` suffix on the ecosystem string itself (not the package name) to scope a record to a particular distro release:

| Ecosystem | Suffix form | Example | Required? |
|-----------|-------------|---------|-----------|
| `Debian` | `:<RELEASE>` (numeric) | `Debian:7` | Optional |
| `Alpine` | `:v<RELEASE>` (numeric, `v` prefix) | `Alpine:v3.16` | **Required** |
| `Rocky` | `:<RELEASE>` (numeric) | `Rocky:9` | Optional |
| `AlmaLinux` | `:<RELEASE>` (numeric) | `AlmaLinux:9` | Optional |

This is part of the OSV spec, not a toolkit convention — the suffix travels inside the `ecosystem` field, so `osv filter -e "Debian:7"` matches only Debian 7 records. Note Alpine is the strict one: its release suffix is **mandatory** and must carry the `v` prefix.

---

## Maven's special naming

Maven package names are `groupId:artifactId`:

```json
{
  "package": {
    "ecosystem": "Maven",
    "name": "org.apache.logging.log4j:log4j-core"
  }
}
```

The CLI's `--maven` flag splits this into separate `groupId` and `artifactId` fields:

```bash
osv query --maven -o json vuln.json
```

```json
{
  "maven": {
    "groupId": "org.apache.logging.log4j",
    "artifactId": "log4j-core"
  }
}
```

---

## Version comparison per ecosystem

Each ecosystem has its own version-ordering rules. OSV's `range.type` field says which rule to use:

| `range.type` | Comparison | Used by |
|--------------|-----------|---------|
| `ECOSYSTEM` | The ecosystem's native ordering | PyPI, npm, Maven, ... |
| `SEMVER` | Strict semver | Libraries that follow semver |
| `GIT` | Commit hash (topological) | Repositories tracked by commit |

This toolkit exposes the `type` and `events` but does **not** perform the comparison itself — you feed the events to a version-comparison library appropriate for the ecosystem. See [Version Range Semantics](/advanced/version-ranges).

---

## Why a registry, not free text?

A fixed ecosystem registry means:

- **No typos** — `PyPI` is the only valid form; `Pypi` won't silently match nothing
- **No duplicates** — every database uses the same string for the same ecosystem
- **Tooling can special-case** — e.g. Maven's `groupId:artifactId` split is only possible because the toolkit knows Maven's convention

---

## See also

- [Ecosystems reference](/reference/ecosystems) — full constant list
- [osv-filter skill](/guide/skills/filter) — how filtering uses ecosystems
- [Version Range Semantics](/advanced/version-ranges) — per-ecosystem comparison
- [Official OSV ecosystem list](https://ossf.github.io/osv-schema/#affectedpackage-field) — canonical source