import { defineConfig } from 'vitepress'
import { withMermaid } from 'vitepress-plugin-mermaid'

// GitHub Pages 子路径：仓库为 scagogogo/osv-schema-skills
// 部署后访问 https://scagogogo.github.io/osv-schema-skills/
// 双语：英文为默认（root），中文在 /zh/ 下
export default withMermaid(
  defineConfig({
    lang: 'en-US', // 站点默认；zh/ locale 在客户端切换 html lang
    title: 'OSV Schema Skills',
    description:
      'AI-native OSV schema toolkit — Go SDK + CLI + 7 Claude Code Skills for parsing, validating, filtering, querying vulnerability data and setup.',
    base: '/osv-schema-skills/',
    cleanDist: true,
    lastUpdated: true,

    head: [
      ['meta', { name: 'theme-color', content: '#3c8c4c' }],
      ['meta', { property: 'og:title', content: 'OSV Schema Skills' }],
      [
        'meta',
        {
          property: 'og:description',
          content: 'AI-native OSV schema toolkit: Go SDK + CLI + Claude Code Skills.',
        },
      ],
      // 语言记忆：在 Vue 挂载前尽早执行。若上次选择了简体中文，且此刻正落在
      // 英文根路径上，则直接重定向到 /zh/，避免"先显示英文再跳转"的闪烁。
      // 写入偏好由 theme/index.ts 的路由 watch 负责，两者配合，不会死循环。
      [
        'script',
        {},
        `;(function(){try{var b='/osv-schema-skills/';var p=location.pathname;if((p===b||p===b+'index.html')&&localStorage.getItem('osv-lang')==='zh'){location.replace(b+'zh/');}}catch(e){}})();`,
      ],
    ],

    // 顶层 locales：控制 per-locale 的 html lang/dir（themeConfig.locales 控制 nav/sidebar/UI 文案）
    locales: {
      root: {
        lang: 'en-US',
      },
      'zh/': {
        lang: 'zh-CN',
        link: '/zh/',
      },
    },

    themeConfig: {
      siteTitle: 'OSV Schema Skills',
      logo: '/logo.svg',

      socialLinks: [
        { icon: 'github', link: 'https://github.com/scagogogo/osv-schema-skills' },
        { icon: 'npm', link: 'https://pkg.go.dev/github.com/scagogogo/osv-schema-skills' },
      ],

      footer: {
        message: 'Released under the MIT License.',
        copyright: 'Copyright © 2024-present scagogogo',
      },

      // —— 英文（默认，root）——
      locales: {
        root: {
          label: 'English',
          nav: [
            // Guide：入门与访问层（下拉分组扁平，按主题排序）
            { text: 'Guide', items: [
              { text: '📖 Introduction', link: '/guide/introduction' },
              { text: '🤖 AI Agent Integration', link: '/guide/ai-agent' },
              { text: 'Installation', link: '/guide/installation' },
              { text: 'Quick Start', link: '/guide/quick-start' },
              { text: 'CLI', link: '/guide/cli' },
              { text: 'Go SDK', link: '/guide/sdk' },
            ]},
            // Skills：7 个技能的总览与逐项
            { text: 'Skills', items: [
              { text: 'Overview', link: '/guide/skills' },
              { text: 'osv-parse', link: '/guide/skills/parse' },
              { text: 'osv-validate', link: '/guide/skills/validate' },
              { text: 'osv-filter', link: '/guide/skills/filter' },
              { text: 'osv-query', link: '/guide/skills/query' },
              { text: 'osv-severity', link: '/guide/skills/severity' },
              { text: 'osv-affected', link: '/guide/skills/affected' },
            ]},
            // Reference：字段、生态、方法
            { text: 'Reference', items: [
              { text: 'OSV Schema', link: '/reference/osv-schema' },
              { text: 'Ecosystems', link: '/reference/ecosystems' },
              { text: 'Methods', link: '/reference/methods' },
            ]},
            // More：实战、FAQ、版本
            { text: 'More', items: [
              { text: '📚 Examples & Cookbook', link: '/guide/examples' },
              { text: '❓ FAQ & Troubleshooting', link: '/guide/faq' },
              { text: '📝 Changelog', link: '/guide/changelog' },
              { text: '⬇️ Releases & Downloads', link: 'https://github.com/scagogogo/osv-schema-skills/releases' },
            ]},
            {
              text: 'GitHub',
              link: 'https://github.com/scagogogo/osv-schema-skills',
            },
          ],

          sidebar: {
            '/guide/': [
              {
                text: 'Getting Started',
                items: [
                  { text: 'Introduction', link: '/guide/introduction' },
                  { text: '🤖 AI Agent', link: '/guide/ai-agent' },
                  { text: 'Installation', link: '/guide/installation' },
                  { text: 'Quick Start', link: '/guide/quick-start' },
                ],
              },
              {
                text: 'Skills',
                collapsed: false,
                items: [
                  { text: 'Overview', link: '/guide/skills' },
                  { text: 'osv-parse', link: '/guide/skills/parse' },
                  { text: 'osv-validate', link: '/guide/skills/validate' },
                  { text: 'osv-filter', link: '/guide/skills/filter' },
                  { text: 'osv-query', link: '/guide/skills/query' },
                  { text: 'osv-severity', link: '/guide/skills/severity' },
                  { text: 'osv-affected', link: '/guide/skills/affected' },
                  { text: 'osv-installation', link: '/guide/installation' },
                ],
              },
              {
                text: 'Access Layers',
                items: [
                  { text: 'CLI', link: '/guide/cli' },
                  { text: 'Go SDK', link: '/guide/sdk' },
                ],
              },
              {
                text: 'Cookbook & Help',
                items: [
                  { text: 'Examples & Cookbook', link: '/guide/examples' },
                  { text: 'FAQ & Troubleshooting', link: '/guide/faq' },
                  { text: 'Changelog', link: '/guide/changelog' },
                ],
              },
            ],
            '/reference/': [
              {
                text: 'Reference',
                items: [
                  { text: 'OSV Schema', link: '/reference/osv-schema' },
                  { text: 'Ecosystems', link: '/reference/ecosystems' },
                  { text: 'Methods', link: '/reference/methods' },
                ],
              },
            ],
          },

          outline: { label: 'On this page' },
          lastUpdated: { text: 'Last updated' },
          docFooter: { prev: 'Previous', next: 'Next' },
          returnToTopLabel: 'Back to top',
          sidebarMenuLabel: 'Menu',
          darkModeSwitchLabel: 'Appearance',
          langMenuLabel: 'Change language',
        },

        // —— 简体中文 ——
        'zh/': {
          label: '简体中文',
          link: '/zh/',
          nav: [
            { text: '指南', items: [
              { text: '📖 项目介绍', link: '/zh/guide/introduction' },
              { text: '🤖 AI Agent 接入', link: '/zh/guide/ai-agent' },
              { text: '安装', link: '/zh/guide/installation' },
              { text: '快速开始', link: '/zh/guide/quick-start' },
              { text: 'CLI 命令行', link: '/zh/guide/cli' },
              { text: 'Go SDK', link: '/zh/guide/sdk' },
            ]},
            { text: 'Skills 技能', items: [
              { text: '总览', link: '/zh/guide/skills' },
              { text: 'osv-parse', link: '/zh/guide/skills/parse' },
              { text: 'osv-validate', link: '/zh/guide/skills/validate' },
              { text: 'osv-filter', link: '/zh/guide/skills/filter' },
              { text: 'osv-query', link: '/zh/guide/skills/query' },
              { text: 'osv-severity', link: '/zh/guide/skills/severity' },
              { text: 'osv-affected', link: '/zh/guide/skills/affected' },
            ]},
            { text: '参考', items: [
              { text: 'OSV Schema', link: '/zh/reference/osv-schema' },
              { text: '生态系统', link: '/zh/reference/ecosystems' },
              { text: '方法清单', link: '/zh/reference/methods' },
            ]},
            { text: '更多', items: [
              { text: '📚 实战示例', link: '/zh/guide/examples' },
              { text: '❓ 常见问题与排错', link: '/zh/guide/faq' },
              { text: '📝 更新日志', link: '/zh/guide/changelog' },
              { text: '⬇️ 版本与下载', link: 'https://github.com/scagogogo/osv-schema-skills/releases' },
            ]},
            {
              text: 'GitHub',
              link: 'https://github.com/scagogogo/osv-schema-skills',
            },
          ],

          sidebar: {
            '/zh/guide/': [
              {
                text: '入门',
                items: [
                  { text: '项目介绍', link: '/zh/guide/introduction' },
                  { text: '🤖 AI Agent 接入', link: '/zh/guide/ai-agent' },
                  { text: '安装', link: '/zh/guide/installation' },
                  { text: '快速开始', link: '/zh/guide/quick-start' },
                ],
              },
              {
                text: 'Skills 技能',
                collapsed: false,
                items: [
                  { text: '总览', link: '/zh/guide/skills' },
                  { text: 'osv-parse', link: '/zh/guide/skills/parse' },
                  { text: 'osv-validate', link: '/zh/guide/skills/validate' },
                  { text: 'osv-filter', link: '/zh/guide/skills/filter' },
                  { text: 'osv-query', link: '/zh/guide/skills/query' },
                  { text: 'osv-severity', link: '/zh/guide/skills/severity' },
                  { text: 'osv-affected', link: '/zh/guide/skills/affected' },
                  { text: 'osv-installation', link: '/zh/guide/installation' },
                ],
              },
              {
                text: '访问层',
                items: [
                  { text: 'CLI 命令行', link: '/zh/guide/cli' },
                  { text: 'Go SDK', link: '/zh/guide/sdk' },
                ],
              },
              {
                text: '实战与帮助',
                items: [
                  { text: '实战示例', link: '/zh/guide/examples' },
                  { text: '常见问题与排错', link: '/zh/guide/faq' },
                  { text: '更新日志', link: '/zh/guide/changelog' },
                ],
              },
            ],
            '/zh/reference/': [
              {
                text: '参考',
                items: [
                  { text: 'OSV Schema', link: '/zh/reference/osv-schema' },
                  { text: '生态系统', link: '/zh/reference/ecosystems' },
                  { text: '方法清单', link: '/zh/reference/methods' },
                ],
              },
            ],
          },

          outline: { label: '本页目录' },
          lastUpdated: { text: '最后更新' },
          docFooter: { prev: '上一页', next: '下一页' },
          returnToTopLabel: '回到顶部',
          sidebarMenuLabel: '菜单',
          darkModeSwitchLabel: '外观',
          langMenuLabel: '语言',
        },
      },

      search: {
        provider: 'local',
      },
    },

    mermaid: {
      // Mermaid 图配置，遵循"一图抵千言"
      theme: 'default',
    },
  }),
)
