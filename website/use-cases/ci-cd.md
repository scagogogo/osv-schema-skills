# CI/CD Integration

Use the `osv` CLI as a validation gate in your continuous integration pipeline — fail the build if any vulnerability record is malformed.

---

## GitHub Actions

Add a workflow step that validates all OSV JSON files before merging:

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

**How it works**: `osv validate` exits with code `1` if any file fails the schema check. GitHub Actions treats non-zero exits as job failures, blocking the PR.

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

Validate locally before you push:

```bash
# .git/hooks/pre-commit
#!/usr/bin/env bash
osv validate advisories/*.json
```

```bash
chmod +x .git/hooks/pre-commit
```

Now every commit runs validation. If any file is invalid, the commit is blocked.

---

## Generate a validation report

Use `-o json` to produce a machine-readable report for downstream tools:

```bash
osv validate -o json advisories/*.json > validation-report.json
```

**Sample output**:

```json
[
  { "file": "advisories/CVE-2021-1234.json", "valid": true },
  { "file": "advisories/GHSA-xyz.json", "valid": false, "errors": ["missing required field: id"] }
]
```

You can then upload this report as an artifact, post it as a PR comment, or feed it to a security dashboard.

---

## See also

- [osv-validate skill](/guide/skills/validate) — skill-level documentation
- [CLI reference](/guide/cli#osv-validate) — flags and exit codes
- [Examples: CI gate](/guide/examples#1-ci-validation-gate) — minimal example
