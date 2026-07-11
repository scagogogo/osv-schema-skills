# Website Bilingual UX Enhancement Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 在已实现的双语官网基础上，补齐两个体验缺口：首页 Hero 区增加显式语言切换按钮，内容页顶部增加"查看另一语言版本"的跨语言链接组件。

**Architecture:** 当前官网已完整支持 EN/ZH 双语（31 页一一对应、locale 配置、导航栏切换、localStorage 记忆、首屏重定向）。数据流：用户进入任一内容页 → Vue 组件 `LangSwitch` 在 `onMounted` 时读取 VitePress 路由 `useData().site.pathname` → 推导对应语言 URL（当前在 `/zh/` 下则去掉前缀得 EN 链接，反之加 `/zh/` 前缀得 ZH 链接）→ 渲染一个带链接的提示条。关键组件：新建 `LangSwitch.vue`（跨语言链接，挂载到 `Layout` 的 `home-hero-after` / `doc-before` 插槽），修改 `theme/index.ts`（注册组件到默认主题的 Layout 插槽），修改两个 `index.md`（首页 actions 加第 4 个按钮）。这样做是因为 VitePress 默认只在导航栏右上角提供小字语言切换，首页 Hero 区和内容页正文无显式入口，用户需要"找"才能切换——加显式入口降低发现成本。

**Tech Stack:** VitePress 1.5.0, Vue 3（VitePress 内置）, vitepress-plugin-mermaid 2.0.17, Node.js 20（CI）

**Risks:**
- Task 2 修改 `theme/index.ts` 注入 Vue 组件，若组件在 SSR/SSG 阶段访问 `window`/`location` 会构建失败 → 缓解：组件内用 `onMounted` 生命周期（仅客户端执行）+ 与现有 `theme/index.ts` 的 `if (typeof window === 'undefined') return` 守卫模式一致
- 跨语言 URL 推导需正确处理 base 路径 `/osv-schema-skills/` → 缓解：组件用 `useData()` 取 VitePress 标准化的 `page.relativePath`，不直接解析 `location.pathname`，避免 base 路径干扰
- VitePress Hero `actions` 加第 4 个按钮可能在窄屏换行 → 缓解：用 `theme: alt` 次要样式，不抢主按钮（brand）视觉；4 个按钮在常见桌面宽度仍单行显示

---

### Task 1: 首页 Hero 增加显式语言切换按钮

**Depends on:** None
**Files:**
- Modify: `website/index.md:11-20`
- Modify: `website/zh/index.md:11-20`

- [ ] **Step 1: 修改 EN 首页 hero actions — 追加"简体中文"切换按钮**
文件: `website/index.md:11-20`（hero.actions 块，在 GitHub 按钮之后追加第 4 项）

```markdown
  actions:
    - theme: brand
      text: 🤖 AI Agent 接入提示词
      link: /guide/ai-agent
    - theme: alt
      text: Quick Start
      link: /guide/quick-start
    - theme: alt
      text: GitHub
      link: https://github.com/scagogogo/osv-schema-skills
    - theme: alt
      text: 🌐 简体中文
      link: /zh/
```

- [ ] **Step 2: 修改 ZH 首页 hero actions — 追加"English"切换按钮**
文件: `website/zh/index.md:11-20`（hero.actions 块，在 GitHub 按钮之后追加第 4 项）

```markdown
  actions:
    - theme: brand
      text: 🤖 AI Agent 接入提示词
      link: /zh/guide/ai-agent
    - theme: alt
      text: 快速开始
      link: /zh/guide/quick-start
    - theme: alt
      text: GitHub
      link: https://github.com/scagogogo/osv-schema-skills
    - theme: alt
      text: 🌐 English
      link: /guide/introduction
```

- [ ] **Step 3: 本地构建验证首页按钮渲染**
Run: `cd website && npm run build 2>&1 | tail -5`
Expected:
  - Exit code: 0
  - Output contains: "building client" and "✨ built"
  - Output does NOT contain: "error" or "Error"

- [ ] **Step 4: 提交**
Run: `git add website/index.md website/zh/index.md && git commit -m "feat(website): add explicit language switch button to homepage hero (EN+ZH)"`

---

### Task 2: 内容页顶部增加跨语言链接组件

**Depends on:** None
**Files:**
- Create: `website/.vitepress/theme/LangSwitch.vue`
- Modify: `website/.vitepress/theme/index.ts`

- [ ] **Step 1: 创建 LangSwitch.vue — 读取当前页路由，推导并渲染对应语言版本链接**

```vue
<script setup>
import { ref, onMounted } from 'vue'
import { useData } from 'vitepress'

const { theme, page } = useData()
const otherLang = ref('')
const otherUrl = ref('')
const otherLabel = ref('')

onMounted(() => {
  if (typeof window === 'undefined') return
  // page.relativePath 形如 "guide/cli.md"（EN）或 "zh/guide/cli.md"（ZH）
  // 不受 base 路径干扰，VitePress 已标准化。
  const rel = page.value.relativePath
  const isZh = rel.startsWith('zh/')
  // 推导对应语言路径：去掉 zh/ 前缀得 EN，或加上 zh/ 得 ZH。
  // index.md 映射到目录根路径。
  let targetRel
  if (isZh) {
    targetRel = rel.slice(3) // "zh/guide/cli.md" → "guide/cli.md"
    otherLang.value = 'en'
    otherLabel.value = '🌐 English'
  } else {
    targetRel = 'zh/' + rel // "guide/cli.md" → "zh/guide/cli.md"
    otherLang.value = 'zh'
    otherLabel.value = '🌐 简体中文'
  }
  // index.md → 目录根：将 "index.md" 替换为 ""，保留目录路径。
  // .md 后缀去掉，转成 VitePress 路由路径（以 / 开头，不含 base）。
  let path = '/' + targetRel.replace(/index\.md$/, '').replace(/\.md$/, '')
  // 去掉尾部斜杠（除非是根路径），与 config.ts 中 link 风格一致。
  if (path.length > 1 && path.endsWith('/')) path = path.slice(0, -1)
  otherUrl.value = path
})
</script>

<template>
  <div v-if="otherUrl" class="lang-switch-banner">
    <a :href="otherUrl">{{ otherLabel }}</a>
  </div>
</template>

<style scoped>
.lang-switch-banner {
  margin: 0 0 1rem 0;
  padding: 0.6rem 1rem;
  border-radius: 6px;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  font-size: 0.9rem;
}
.lang-switch-banner a {
  color: var(--vp-c-brand);
  text-decoration: none;
  font-weight: 500;
}
.lang-switch-banner a:hover {
  text-decoration: underline;
}
</style>
```

- [ ] **Step 2: 修改 theme/index.ts — 用 h 函数包装默认 Layout 并在 doc-before 插槽注入 LangSwitch**
文件: `website/.vitepress/theme/index.ts`（完整替换整个文件）

```typescript
// 文件: website/.vitepress/theme/index.ts（完整替换）
import DefaultTheme from 'vitepress/theme'
import { watch, h } from 'vue'
import { useRouter } from 'vitepress'
import LangSwitch from './LangSwitch.vue'

const BASE = '/osv-schema-skills/'
const STORAGE_KEY = 'osv-lang'

function detectLang(pathname: string): 'zh' | 'en' {
  return pathname.startsWith(BASE + 'zh/') || pathname === BASE + 'zh'
    ? 'zh'
    : 'en'
}

export default {
  extends: DefaultTheme,
  // 包装默认 Layout，在 doc-before 插槽注入跨语言切换链接。
  // 仅 doc 布局生效（home 布局无 doc-before 插槽），不影响首页。
  Layout: () => {
    return h(DefaultTheme.Layout, null, {
      'doc-before': () => h(LangSwitch),
    })
  },
  setup() {
    if (typeof window === 'undefined') return

    const router = useRouter()

    const persist = () => {
      try {
        localStorage.setItem(STORAGE_KEY, detectLang(location.pathname))
      } catch {
        // localStorage 不可用（隐私模式等）时静默降级。
      }
    }

    persist()
    watch(() => router.route.path, persist)
  },
}
```

- [ ] **Step 3: 本地构建验证组件渲染无 SSR 错误**
Run: `cd website && npm run build 2>&1 | tail -10`
Expected:
  - Exit code: 0
  - Output contains: "✨ built" or "building"
  - Output does NOT contain: "error during build" or "window is not defined" or "ReferenceError"

- [ ] **Step 4: 本地预览验证跨语言链接正确性**
Run: `cd website && timeout 15 npx vitepress dev --port 5173 &> /tmp/osv-dev.log & sleep 8 && curl -s http://localhost:5173/guide/cli | grep -c "lang-switch-banner" ; curl -s http://localhost:5173/zh/guide/cli | grep -c "lang-switch-banner" ; pkill -f vitepress`
Expected:
  - Exit code: 0
  - 第一个 curl 输出 ≥ 1（EN 页面有跨语言链接容器）
  - 第二个 curl 输出 ≥ 1（ZH 页面有跨语言链接容器）

- [ ] **Step 5: 提交**
Run: `git add website/.vitepress/theme/LangSwitch.vue website/.vitepress/theme/index.ts && git commit -m "feat(website): add cross-language link banner atop content pages (EN+ZH)"`

---

### Task 3: 部署验证

**Depends on:** Task 1, Task 2
**Files:**
- None（仅运行验证）

- [ ] **Step 1: 推送触发 website.yml workflow**
Run: `git push origin main`
Expected:
  - Exit code: 0
  - `git push` 输出不含 "error" 或 "rejected"

- [ ] **Step 2: 监控 workflow 直至完成**
Run: `gh run watch $(gh run list --workflow=website.yml --limit=1 --json databaseId -q '.[0].databaseId') --exit-status`
Expected:
  - Exit code: 0
  - Output contains: "✓ website" and "success"

- [ ] **Step 3: curl 验证部署后的首页语言按钮与内容页跨语言链接**
Run: `curl -s https://scagogogo.github.io/osv-schema-skills/ | grep -c "简体中文" ; curl -s https://scagogogo.github.io/osv-schema-skills/zh/ | grep -c "English" ; curl -s https://scagogogo.github.io/osv-schema-skills/guide/cli.html | grep -c "lang-switch-banner"`
Expected:
  - Exit code: 0
  - 首页 EN 输出 ≥ 1（含"简体中文"按钮）
  - 首页 ZH 输出 ≥ 1（含"English"按钮）
  - 内容页输出 ≥ 1（含跨语言链接容器）
