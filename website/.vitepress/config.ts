import { defineConfig } from 'vitepress'
import { withMermaid } from 'vitepress-plugin-mermaid'

// GitHub Pages 子路径：仓库为 scagogogo/osv-schema-skills
// 部署后访问 https://scagogogo.github.io/osv-schema-skills/
export default withMermaid(
  defineConfig({
    lang: 'en-US',
    title: 'OSV Schema Skills',
    description:
      'AI-native OSV schema toolkit — Go SDK + CLI + 6 Claude Code Skills for parsing, validating, filtering and querying vulnerability data.',
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
    ],

    themeConfig: {
      siteTitle: 'OSV Schema Skills',
      logo: '/logo.svg',

      nav: [
        { text: 'Guide', items: [
          { text: 'Introduction', link: '/guide/introduction' },
          { text: '🤖 AI Agent 接入', link: '/guide/ai-agent' },
          { text: 'Installation', link: '/guide/installation' },
          { text: 'Quick Start', link: '/guide/quick-start' },
          { text: 'Skills', link: '/guide/skills' },
          { text: 'CLI', link: '/guide/cli' },
          { text: 'Go SDK', link: '/guide/sdk' },
        ]},
        { text: 'Reference', items: [
          { text: 'OSV Schema', link: '/reference/osv-schema' },
          { text: 'Ecosystems', link: '/reference/ecosystems' },
          { text: 'Methods', link: '/reference/methods' },
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
              { text: '🤖 AI Agent 接入', link: '/guide/ai-agent' },
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
            ],
          },
          {
            text: 'Access Layers',
            items: [
              { text: 'CLI', link: '/guide/cli' },
              { text: 'Go SDK', link: '/guide/sdk' },
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

      socialLinks: [
        { icon: 'github', link: 'https://github.com/scagogogo/osv-schema-skills' },
        { icon: 'npm', link: 'https://pkg.go.dev/github.com/scagogogo/osv-schema-skills' },
      ],

      footer: {
        message: 'Released under the MIT License.',
        copyright: 'Copyright © 2024-present scagogogo',
      },

      search: {
        provider: 'local',
      },

      outline: {
        label: 'On this page',
      },

      lastUpdated: {
        text: 'Last updated',
      },
    },

    mermaid: {
      // Mermaid 图配置，遵循"一图抵千言"
      theme: 'default',
    },
  }),
)
