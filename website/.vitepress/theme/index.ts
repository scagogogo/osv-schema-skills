// 自定义主题：在默认主题基础上，"记住用户的语言选择"。
//
// VitePress 默认只按 URL 路径决定语言（/ = English，/zh/ = 简体中文），
// 刷新或下次访问落到根路径时永远是英文，不会记住上次的选择。
//
// 这里做两件事：
//   1. 每次路由变化时，把当前语言写入 localStorage（键 osv-lang）。
//   2. 首屏落在英文根路径时，若上次选择是 zh，则由 config.ts 注入的 head 内联脚本
//      在 Vue 挂载前就重定向到 /zh/（避免英文闪一下）。这里的 watch 负责"记录选择"，
//      head 脚本负责"应用选择"，两者配合且不会造成重定向死循环：
//      用户从中文手动切回英文时是 SPA 跳转（不刷新），watch 立即把偏好改为 en，
//      因此不会被重定向弹回中文。
import DefaultTheme from 'vitepress/theme'
import { watch } from 'vue'
import { useRouter } from 'vitepress'

const BASE = '/osv-schema-skills/'
const STORAGE_KEY = 'osv-lang'

function detectLang(pathname: string): 'zh' | 'en' {
  return pathname.startsWith(BASE + 'zh/') || pathname === BASE + 'zh'
    ? 'zh'
    : 'en'
}

export default {
  extends: DefaultTheme,
  setup() {
    // SSR/SSG 阶段没有 window/localStorage，直接跳过。
    if (typeof window === 'undefined') return

    const router = useRouter()

    const persist = () => {
      try {
        localStorage.setItem(STORAGE_KEY, detectLang(location.pathname))
      } catch {
        // localStorage 不可用（隐私模式等）时静默降级。
      }
    }

    // 首次进入立即记录一次，之后每次路由变化再记录。
    persist()
    watch(() => router.route.path, persist)
  },
}
