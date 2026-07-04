# CI/CD 集成

在持续集成流水线中用 `osv` CLI 做校验闸门——如果漏洞记录格式错误则阻断构建。

---

## GitHub Actions

添加一个 workflow 步骤，合并前校验所有 OSV JSON 文件：

```yaml
# .github/workflows/validate-osv.yml
name: Validate OSV records
on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install osv CLI
        run: |
          curl -fsSL https://github.com/scagogogo/osv-schema-skills/releases/download/v0.1.0/osv_v0.1.0_linux_amd64.tar.gz \
            | tar -xz osv && sudo mv osv /usr/local/bin/

      - name: Validate all OSV files
        run: osv validate advisories/*.json
```

**原理**：`osv validate` 在任一文件未通过 schema 检查时以退出码 `1` 退出。GitHub Actions 视非零退出为任务失败，从而阻断 PR。

---

## GitLab CI

```yaml
# .gitlab-ci.yml
validate-osv:
  stage: test
  image: golang:1.22
  before_script:
    - go install github.com/scagogogo/osv-schema-skills/cmd/osv@latest
  script:
    - osv validate advisories/*.json
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
```

---

## Jenkins Pipeline

```groovy
pipeline {
  agent any
  stages {
    stage('Validate OSV') {
      steps {
        sh '''
          curl -fsSL https://github.com/scagogogo/osv-schema-skills/releases/download/v0.1.0/osv_v0.1.0_linux_amd64.tar.gz \
            | tar -xz osv && chmod +x osv
          ./osv validate advisories/*.json
        '''
      }
    }
  }
}
```

---

## Pre-commit hook

推送前在本地校验：

```bash
# .git/hooks/pre-commit
#!/usr/bin/env bash
osv validate advisories/*.json
```

```bash
chmod +x .git/hooks/pre-commit
```

每次提交都会执行校验。如有文件无效，提交会被阻止。

---

## 生成校验报告

用 `-o json` 产出机器可读报告供下游工具消费：

```bash
osv validate -o json advisories/*.json > validation-report.json
```

**示例输出**：

```json
[
  { "file": "advisories/CVE-2021-1234.json", "valid": true },
  { "file": "advisories/GHSA-xyz.json", "valid": false, "errors": ["missing required field: id"] }
]
```

可将其作为 artifact 上传、以 PR 评论发布、或喂给安全仪表盘。

---

## 另见

- [osv-validate 技能](/zh/guide/skills/validate) —— 技能级文档
- [CLI 参考](/zh/guide/cli#osv-validate) —— 标志与退出码
- [实战示例：CI 闸门](/zh/guide/examples#1-ci-校验闸门) —— 最小示例
