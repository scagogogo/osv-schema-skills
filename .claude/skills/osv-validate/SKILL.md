---
name: osv-validate
description: Validate OSV (Open Source Vulnerability) JSON files against the schema. Triggers on mentions of OSV validation, vulnerability format checking, schema compliance, or when user wants to verify an OSV JSON file is well-formed.
allowed-tools: "Bash(osv:*)"
argument-hint: <osv-json-file> [file...]
---

# OSV Validate

> **Setup:** See `/osv-installation` for one-time CLI/SDK install.
> **Layers:** SDK (Go) → CLI (shell) — pick your entry point.

## When to Use

- Validate one or more OSV JSON files against the schema
- Check if a vulnerability JSON file is well-formed and parseable
- Verify required fields (id, schema_version) are present
- Batch validate multiple OSV files in a directory

## Decision Tree

```
OSV JSON file(s) → what do you need?
├─ Quick validity check?        → osv validate <file>
├─ Batch validate?              → osv validate file1.json file2.json
├─ Machine-readable result?     → osv validate -o json <file>
└─ Programmatic validation?     → SDK: UnmarshalFromJson() + field checks
```

## Task Patterns

### Validate a single file

**Goal:** Check if `vulnerability.json` is a valid OSV file.

| Layer | Approach |
|-------|----------|
| CLI | `osv validate vulnerability.json` |

Output: `✓ vulnerability.json (id=GHSA-xxx, schema_version=1.4.0)` or `✗ vulnerability.json` with error details.

### Validate multiple files

**Goal:** Batch validate all OSV files in a directory.

| Layer | Approach |
|-------|----------|
| CLI | `osv validate file1.json file2.json file3.json` |

### Validate with JSON output

**Goal:** Get structured validation results for automation.

| Layer | Approach |
|-------|----------|
| CLI | `osv validate -o json vulnerability.json` |

Returns: `[{file, valid, errors[], id, schema_version}]`

### Programmatic validation

**Goal:** Validate an OSV file in Go code.

```go
raw, err := os.ReadFile("vulnerability.json")
if err != nil {
    return err
}
if !json.Valid(raw) {
    return fmt.Errorf("not valid JSON")
}
osvData, err := osv.UnmarshalFromJson[any, any](raw)
if err != nil {
    return fmt.Errorf("OSV parse error: %w", err)
}
if osvData.ID == "" {
    return fmt.Errorf("missing required field: id")
}
if osvData.SchemaVersion == "" {
    return fmt.Errorf("missing required field: schema_version")
}
```

## API Reference

### CLI Commands

```bash
osv validate <file>                    # Validate single file (text output)
osv validate <file1> <file2> ...       # Validate multiple files
osv validate -o json <file>            # JSON output format
```

### Validation Checks

| Check | Description |
|-------|-------------|
| File exists | Can the file be read? |
| Valid JSON | Is the content valid JSON? |
| OSV parseable | Can it be unmarshalled into OsvSchema? |
| Required: `id` | Is the `id` field non-empty? |
| Required: `schema_version` | Is the `schema_version` field non-empty? |

### Validation Result Format (JSON)

```json
[
  {
    "file": "vulnerability.json",
    "valid": true,
    "id": "GHSA-vxv8-r8q2-63xw",
    "schema_version": "1.4.0"
  }
]
```

Error case:
```json
[
  {
    "file": "bad.json",
    "valid": false,
    "errors": ["missing required field: id"]
  }
]
```

## Cross-References

- [[osv-parse]] — parse and display OSV data
- [[osv-filter]] — filter validated data by various criteria
- [[osv-query]] — extract specific sub-information

## Important Notes

- Exit code 0 = all files valid, exit code 1 = at least one file invalid
- Validation checks required fields per OSV Schema spec: `id` and `schema_version`
- The validator does NOT check against the full JSON Schema — it verifies parseability and required fields
- For full schema validation, use the OSV Schema specification at https://ossf.github.io/osv-schema/
