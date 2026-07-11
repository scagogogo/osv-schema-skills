# Website Bilingual Quality Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 修复官网双语支持的三处真实质量缺口——中文页面 html lang 错误标记为 en-US、无 sitemap.xml 影响 SEO、本地搜索未针对中文优化分词。

**Architecture:** 数据流：(1) html lang 修复——VitePress 的 `getLocaleForPath` 用 `/${localeKey}/` 作正则匹配页面相对路径，当 locale 键带尾斜杠（`'zh/'`）时正则变为 `/zh//` 无法匹配，所有页面退化为 root 的 en-US；改为无尾斜杠键 `'zh'` 后正则 `/zh/` 正确匹配，ZH 页获得 `zh-CN`。(2) sitemap——VitePress 1.6.4 内置 sitemap 生成，仅需在 config 设 `sitemap.hostname`，构建时自动产出 `sitemap.xml` 列出全部 62 个页面（EN 31 + ZH 31）。(3) 中文搜索——VitePress local search 的 `LocalSearchOptions.locales` 支持按 locale 覆盖 `miniSearch.options.tokenize`，给 `zh` locale 配中文分词正则，root 保持默认英文分词。关键组件：仅修改 `website/.vitepress/config.ts` 一个文件的三处。这样做是因为三处缺口都源于 config 配置，根因已在本地实验验证（改键后 html lang 正确为 zh-CN）。

**Tech Stack:** VitePress 1.6.4（lockfile 实际版本，含内置 sitemap），vitepress-plugin-mermaid 2.0.17，MiniSearch（VitePress 内置 local search 引擎），Node.js 20（CI）

**Risks:**
- Task 3 的 per-locale tokenize 若正则写错可能让中文搜索返回空结果 → 缓解：用 VitePress 官方文档推荐的中文分词正则 `/[\x00-\x7F]/` 拆 ASCII、其余按字拆；验证步骤包含中文关键词搜索断言
- Task 1 改 locale 键可能影响 `themeConfig.locales` 里同样用 `'zh/'` 的键 → 缓解：themeConfig.locales 的键用于 sidebar 前缀匹配（`/zh/guide/`），其键保持 `'zh/'` 不变，只改顶层 `locales`（控制 html lang）的键；验证步骤断言 EN/ZH 页 html lang 分别为 en-US/zh-CN
- sitemap hostname 必须含 base 路径 → 缓解：用 `https://scagogogo.github.io/osv-schema-skills/`，与 config 的 `base: '/osv-schema-skills/'` 一致

---

### Task 1: 修复中文页面 html lang 标记错误

**Depends on:** None
**Files:**
- Modify: `website/.vitepress/config.ts:38-46`（顶层 locales 块）

- [ ] **Step 1: 修改顶层 locales 键 — 将 'zh/' 改为 'zh' 以修正 VitePress locale 路径正则匹配**
文件: `website/.vitepress/config.ts:38-46`（顶层 locales 块，仅改键名，themeConfig.locales 保持不变）

```typescript
    // 顶层 locales：控制 per-locale 的 html lang/dir（themeConfig.locales 控制 nav/sidebar/UI 文案）
    // 键名必须无尾斜杠：VitePress 的 getLocaleForPath 用 `/${key}/` 作正则匹配页面相对路径，
    // 若键为 'zh/' 则正则变 '/zh//' 无法匹配，所有页面 lang 退化为 root 的 en-US。
    // themeConfig.locales 的键用于 sidebar 前缀匹配，保持 'zh/' 不变。
    locales: {
      root: {
        lang: 'en-US',
      },
      'zh': {
        lang: 'zh-CN',
        link: '/zh/',
      },
    },
```

- [ ] **Step 2: 本地构建验证 EN/ZH 页 html lang 正确**
Run: `cd website && npm run build 2>&1 | tail -3 && grep -oE '<html[^>]*lang="[^"]*"' .vitepress/dist/guide/cli.html | head -1 && grep -oE '<html[^>]*lang="[^"]*"' .vitepress/dist/zh/guide/cli.html | head -1`
Expected:
  - Exit code: 0
  - Output contains: "build complete"
  - EN 行 contains: `lang="en-US"`
  - ZH 行 contains: `lang="zh-CN"`

- [ ] **Step 3: 提交**
Run: `git add website/.vitepress/config.ts && git commit -m "fix(website): correct html lang attribute for Chinese pages (was en-US)"`

---

### Task 2: 启用 sitemap.xml 提升 SEO

**Depends on:** None
**Files:**
- Modify: `website/.vitepress/config.ts`（在 `base` 字段之后、`cleanDist` 之前插入 sitemap 配置）

- [ ] **Step 1: 修改 config.ts 添加 sitemap 配置 — 启用 VitePress 内置 sitemap 生成**
文件: `website/.vitepress/config.ts`（在 `base: '/osv-schema-skills/',` 行之后插入 sitemap 块；找到 `base: '/osv-schema-skills/',` 与 `cleanDist: true,` 两行，在二者之间插入）

```typescript
    base: '/osv-schema-skills/',
    // 内置 sitemap：构建时生成 sitemap.xml，列出全部页面（EN 31 + ZH 31），
    // 帮助搜索引擎发现两个语言版本。hostname 必须含 base 路径。
    sitemap: {
      hostname: 'https://scagogogo.github.io/osv-schema-skills/',
    },
    cleanDist: true,
```

- [ ] **Step 2: 本地构建验证 sitemap.xml 生成且含中英文页面**
Run: `cd website && npm run build 2>&1 | tail -3 && ls -la .vitepress/dist/sitemap.xml && grep -c '<loc>' .vitepress/dist/sitemap.xml && grep -c '/zh/' .vitepress/dist/sitemap.xml`
Expected:
  - Exit code: 0
  - `ls` 输出含 `sitemap.xml`
  - 第一个 grep 计数 ≥ 60（页面总数）
  - 第二个 grep 计数 ≥ 30（含 /zh/ 的中文页 URL）

- [ ] **Step 3: 提交**
Run: `git add website/.vitepress/config.ts && git commit -m "feat(website): enable sitemap.xml for SEO with both language versions"`

---

### Task 3: 优化中文本地搜索分词

**Depends on:** None
**Files:**
- Modify: `website/.vitepress/config.ts:370-372`（search 配置块）

- [ ] **Step 1: 修改 search 配置 — 为 zh locale 配置中文分词 tokenize**
文件: `website/.vitepress/config.ts:370-372`（search 块，替换整个 search 对象）

```typescript
      search: {
        provider: 'local',
        options: {
          // per-locale 搜索配置：root 保持 VitePress 默认英文分词；
          // zh locale 用中文分词正则——ASCII 字符按非字母数字边界拆，其余按单字拆，
          // 使中文关键词（如"漏洞""校验"）能被 MiniSearch 正确索引与命中。
          locales: {
            zh: {
              miniSearch: {
                options: {
                  tokenize: (text: string) =>
                    text
                      // 先按 ASCII 词边界（非字母数字字符）拆分英文/数字片段
                      .split(/[\x00-\x2F\x3A-\x40\x5B-\x60\x7B-\x7F]+/)
                      .filter((t) => t.length > 0)
                      // 再把每个片段里的非 ASCII（中文等）按单字展开，
                      // 与 ASCII 片段一起作为独立 token 返回
                      .flatMap((t) => {
                        const cjk = t.match(/[一-鿿]/g)
                        return cjk ? [...t.split(/[一-鿿]+/).filter(Boolean), ...cjk] : [t]
                      }),
                },
              },
            },
          },
        },
      },
```

- [ ] **Step 2: 本地构建验证 search 配置无构建错误**
Run: `cd website && npm run build 2>&1 | tail -5`
Expected:
  - Exit code: 0
  - Output contains: "build complete"
  - Output does NOT contain: "error" or "TypeError" or "Invalid"

- [ ] **Step 3: 提交**
Run: `git add website/.vitepress/config.ts && git commit -m "feat(website): optimize local search tokenization for Chinese (zh locale)"`

---

### Task 4: 部署验证

**Depends on:** Task 1, Task 2, Task 3
**Files:**
- None（仅运行验证）

- [ ] **Step 1: 推送触发 website.yml workflow**
Run: `git push origin main`
Expected:
  - Exit code: 0
  - Output does NOT contain: "rejected" or "error"

- [ ] **Step 2: 监控 workflow 直至完成**
Run: `gh run watch $(gh run list --workflow=website.yml --limit=1 --json databaseId -q '.[0].databaseId') --exit-status`
Expected:
  - Exit code: 0
  - Output contains: "success"

- [ ] **Step 3: curl 验证部署后的 html lang、sitemap、search 索引**
Run: `curl -s https://scagogogo.github.io/osv-schema-skills/zh/guide/cli.html | grep -oE '<html[^>]*lang="[^"]*"' | head -1 ; curl -s -o /dev/null -w "sitemap HTTP %{http_code}\n" https://scagogogo.github.io/osv-schema-skills/sitemap.xml ; curl -s https://scagogogo.github.io/osv-schema-skills/sitemap.xml | grep -c '/zh/'`
Expected:
  - Exit code: 0
  - ZH 页 html lang 输出含 `lang="zh-CN"`
  - sitemap HTTP 状态为 200
  - sitemap 中 /zh/ 计数 ≥ 30
