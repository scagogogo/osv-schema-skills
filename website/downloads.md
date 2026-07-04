# Downloads

Pre-built binaries for the `osv` CLI are published with every release. No Go toolchain needed — download, extract, run.

::: tip Latest version
**v0.1.0** — released 2026-07-05. See the [Changelog](/guide/changelog) for what's new.
:::

---

## Pre-built binaries

| OS | Architecture | File | Size |
|----|--------------|------|------|
| 🐧 Linux | amd64 (x86_64) | `osv_v0.1.0_linux_amd64.tar.gz` | ~0.9 MB |
| 🐧 Linux | arm64 (aarch64) | `osv_v0.1.0_linux_arm64.tar.gz` | ~0.9 MB |
| 🐧 Linux | arm v7 | `osv_v0.1.0_linux_arm.tar.gz` | ~0.9 MB |
| 🍎 macOS | amd64 (Intel) | `osv_v0.1.0_darwin_amd64.tar.gz` | ~1.0 MB |
| 🍎 macOS | arm64 (Apple Silicon) | `osv_v0.1.0_darwin_arm64.tar.gz` | ~1.0 MB |
| 🪟 Windows | amd64 | `osv_v0.1.0_windows_amd64.zip` | ~1.0 MB |
| 🪟 Windows | arm64 | `osv_v0.1.0_windows_arm64.zip` | ~0.9 MB |

All files are on the [GitHub Release page](https://github.com/scagogogo/osv-schema-skills/releases/tag/v0.1.0). Each release also ships `checksums.txt` for integrity verification.

---

## Quick install (one-liner)

### Linux / macOS

```bash
# Set your platform here. Example: Linux amd64
VERSION=v0.1.0
OS=linux
ARCH=amd64
curl -fsSL -o osv.tar.gz \
  https://github.com/scagogogo/osv-schema-skills/releases/download/${VERSION}/osv_${VERSION}_${OS}_${ARCH}.tar.gz
tar -xzf osv.tar.gz osv
chmod +x osv && sudo mv osv /usr/local/bin/
osv version
```

### Windows (PowerShell)

```powershell
$VERSION = "v0.1.0"
$ARCH = "amd64"
Invoke-WebRequest -Uri "https://github.com/scagogogo/osv-schema-skills/releases/download/$VERSION/osv_${VERSION}_windows_${ARCH}.zip" -OutFile "osv.zip"
Expand-Archive -Path "osv.zip" -DestinationPath "."
.\osv.exe version
```

---

## Verify the checksum

Every release publishes `checksums.txt`. Verify your download before trusting it:

```bash
# Download both the archive and checksums.txt
curl -fsSL -O https://github.com/scagogogo/osv-schema-skills/releases/download/v0.1.0/osv_v0.1.0_linux_amd64.tar.gz
curl -fsSL -O https://github.com/scagogogo/osv-schema-skills/releases/download/v0.1.0/checksums.txt

# Verify (only check the file you downloaded)
sha256sum -c checksums.txt --ignore-missing
```

Expected output:

```
osv_v0.1.0_linux_amd64.tar.gz: OK
```

---

## Fall back to `go install`

If a release has no pre-built assets for your platform, or you already have Go 1.18+:

```bash
go install github.com/scagogogo/osv-schema-skills/cmd/osv@latest
```

This installs to `$GOPATH/bin` (or `$HOME/go/bin`). Make sure that directory is on your `PATH`.

---

## Build from source

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills
go build -o osv ./cmd/osv
./osv version
```

To inject a specific version string:

```bash
go build -ldflags "-X main.cliVersion=v0.1.0" -o osv ./cmd/osv
./osv version
# osv-cli version: v0.1.0
# OSV schema version: 1.4.0
```

---

## AI Agent: auto-detect platform

An AI agent can pick the right binary by inspecting the OS/arch:

```bash
#!/usr/bin/env bash
# Auto-install osv CLI for the current platform
VERSION=v0.1.0
OS=$(uname -s | tr '[:upper:]' '[:lower:]')   # linux / darwin
ARCH=$(uname -m)                                # x86_64 / arm64
case "$ARCH" in
  x86_64) ARCH=amd64 ;;
  aarch64) ARCH=arm64 ;;
esac
curl -fsSL "https://github.com/scagogogo/osv-schema-skills/releases/download/${VERSION}/osv_${VERSION}_${OS}_${ARCH}.tar.gz" \
  | tar -xz osv && chmod +x osv && sudo mv osv /usr/local/bin/
osv version
```

---

## See also

- [Installation guide](/guide/installation) — full setup walkthrough
- [Quick Start](/guide/quick-start) — running against a real record in 30 seconds
- [Changelog](/guide/changelog) — what changed in each release
- [GitHub Releases](https://github.com/scagogogo/osv-schema-skills/releases) — full release history
