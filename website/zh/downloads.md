# 下载

`osv` CLI 的预编译二进制随每个版本发布。无需 Go 工具链——下载、解压、运行。

::: tip 最新版本
**v0.1.0** —— 2026-07-05 发布。新功能见 [更新日志](/zh/guide/changelog)。
:::

---

## 预编译二进制

| 系统 | 架构 | 文件 | 大小 |
|------|------|------|------|
| 🐧 Linux | amd64 (x86_64) | `osv_v0.1.0_linux_amd64.tar.gz` | ~0.9 MB |
| 🐧 Linux | arm64 (aarch64) | `osv_v0.1.0_linux_arm64.tar.gz` | ~0.9 MB |
| 🐧 Linux | arm v7 | `osv_v0.1.0_linux_arm.tar.gz` | ~0.9 MB |
| 🍎 macOS | amd64 (Intel) | `osv_v0.1.0_darwin_amd64.tar.gz` | ~1.0 MB |
| 🍎 macOS | arm64 (Apple Silicon) | `osv_v0.1.0_darwin_arm64.tar.gz` | ~1.0 MB |
| 🪟 Windows | amd64 | `osv_v0.1.0_windows_amd64.zip` | ~1.0 MB |
| 🪟 Windows | arm64 | `osv_v0.1.0_windows_arm64.zip` | ~0.9 MB |

所有文件在 [GitHub Release 页面](https://github.com/scagogogo/osv-schema-skills/releases/tag/v0.1.0)。每个版本还附带 `checksums.txt` 供完整性校验。

---

## 一行快速安装

### Linux / macOS

```bash
# 在此设置你的平台。示例：Linux amd64
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

## 校验 checksum

每个版本都发布 `checksums.txt`。信任下载前请先校验：

```bash
# 同时下载 archive 和 checksums.txt
curl -fsSL -O https://github.com/scagogogo/osv-schema-skills/releases/download/v0.1.0/osv_v0.1.0_linux_amd64.tar.gz
curl -fsSL -O https://github.com/scagogogo/osv-schema-skills/releases/download/v0.1.0/checksums.txt

# 校验（只检查你下载的那个文件）
sha256sum -c checksums.txt --ignore-missing
```

期望输出：

```
osv_v0.1.0_linux_amd64.tar.gz: OK
```

---

## 回退到 `go install`

若某版本未提供你平台的预编译资产，或你已有 Go 1.18+：

```bash
go install github.com/scagogogo/osv-schema-skills/cmd/osv@latest
```

安装到 `$GOPATH/bin`（或 `$HOME/go/bin`）。确保该目录在你的 `PATH` 上。

---

## 从源码构建

```bash
git clone https://github.com/scagogogo/osv-schema-skills.git
cd osv-schema-skills
go build -o osv ./cmd/osv
./osv version
```

注入指定版本号：

```bash
go build -ldflags "-X main.cliVersion=v0.1.0" -o osv ./cmd/osv
./osv version
# osv-cli version: v0.1.0
# OSV schema version: 1.4.0
```

---

## AI Agent：自动检测平台

AI Agent 可通过检测 OS/arch 自动选择正确的二进制：

```bash
#!/usr/bin/env bash
# 为当前平台自动安装 osv CLI
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

## 另见

- [安装指南](/zh/guide/installation) —— 完整安装步骤
- [快速开始](/zh/guide/quick-start) —— 30 秒内对真实记录跑起来
- [更新日志](/zh/guide/changelog) —— 每个版本的变更
- [GitHub Releases](https://github.com/scagogogo/osv-schema-skills/releases) —— 完整版本历史
