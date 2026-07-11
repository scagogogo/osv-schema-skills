<script setup>
import { computed } from 'vue'
import { useData } from 'vitepress'

const { page } = useData()

// 用 computed 而非 onMounted：page.relativePath 在 SSR/SSG 阶段即为已知值，
// 这样静态 HTML 初始就含 banner（爬虫、禁用 JS 用户、curl 均可见），
// 无需等待客户端 hydration。无 window/location 访问，SSR 安全。
const other = computed(() => {
  const rel = page.value.relativePath
  if (!rel) return { url: '', label: '' }
  const isZh = rel.startsWith('zh/')
  let targetRel
  if (isZh) {
    targetRel = rel.slice(3) // "zh/guide/cli.md" → "guide/cli.md"
  } else {
    targetRel = 'zh/' + rel // "guide/cli.md" → "zh/guide/cli.md"
  }
  // index.md → 目录根；.md 后缀去掉；转成 VitePress 路由路径（/ 开头，不含 base）。
  let path = '/' + targetRel.replace(/index\.md$/, '').replace(/\.md$/, '')
  if (path.length > 1 && path.endsWith('/')) path = path.slice(0, -1)
  return {
    url: path,
    label: isZh ? '🌐 English' : '🌐 简体中文',
  }
})
</script>

<template>
  <div v-if="other.url" class="lang-switch-banner">
    <a :href="other.url">{{ other.label }}</a>
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
