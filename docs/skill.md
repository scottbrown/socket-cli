---
name: socket-cli
description: >
  Interact with the Socket.dev supply chain security platform using the socket
  CLI. Use for investigating alerts, listing repos, running scans, looking up
  packages, checking quota, and triaging supply chain risks.
when_to_use: >
  Socket.dev queries, supply chain security, dependency alerts, package
  scoring, full scans, diff scans, SBOM, purl lookups, Socket API,
  dependency risk assessment, npm/pypi/go package security
allowed-tools: Bash Read Grep Glob
user-invocable: true
---

# Socket CLI Assistant

This skill provides interactive access to the Socket.dev supply chain security
platform via the `socket` CLI.

## Prerequisites

The CLI authenticates via the `SOCKET_API_TOKEN` environment variable or a
token file at `~/.config/socket/token`. Verify with:
```bash
test -n "$SOCKET_API_TOKEN" && echo "Token set" || test -f ~/.config/socket/token && echo "Token file exists" || echo "No auth configured"
```

The `socket` binary must be on `$PATH`. If not, build it from source:
```bash
go build -o socket ./cmd/socket
```

## Available Commands

### Organizations
```bash
socket orgs list
```

### Repositories
```bash
socket repos list --org <slug>
socket repos get --org <slug> --repo <repo-slug>
socket repos delete --org <slug> --repo <repo-slug>
```

### Full Scans
```bash
socket fullscans list --org <slug> [--repo <name>] [--branch <branch>] [--per-page N] [--page N]
socket fullscans get --org <slug> --id <scan-id>
socket fullscans create --org <slug> --repo <name> --file <manifest> [--branch <branch>] [--commit <hash>]
socket fullscans delete --org <slug> --id <scan-id>
socket fullscans metadata --org <slug> --id <scan-id>
```

### Diff Scans
```bash
socket diffscans list --org <slug> [--per-page N] [--page N]
socket diffscans get --org <slug> --id <scan-id>
socket diffscans delete --org <slug> --id <scan-id>
```

### Alerts
```bash
socket alerts list --org <slug> [--per-page N] [--page N]
socket alerts triage --org <slug>
```

### Packages
```bash
socket packages lookup <purl> [<purl>...] [--org <slug>]
```
Package URLs follow the purl spec: `pkg:<ecosystem>/<name>@<version>`
Examples: `pkg:npm/express@4.18.2`, `pkg:pypi/requests@2.31.0`, `pkg:golang/github.com/spf13/cobra@1.8.0`

### Quota
```bash
socket quota
```

## Discovering Your Org Slug

If you don't know your org slug, run `socket orgs list` first and use the
`slug` field from the response.

## Interpreting Results

### Alert Severity
- **critical** — immediate action required, likely active exploit or malware
- **high** — significant supply chain risk (obfuscated code, known malicious patterns)
- **medium** — concerning patterns worth investigating
- **low** — informational, minor quality concerns

### Alert Categories
- `supplyChainRisk` — risks related to package provenance, publishing patterns, or code integrity
- `vulnerability` — known CVEs or security vulnerabilities
- `quality` — code quality and maintenance signals
- `license` — license compatibility issues

### Package Scores (0.0 to 1.0)
- `overall` — composite score
- `supplyChain` — publishing practices, contributor patterns
- `quality` — code quality signals
- `maintenance` — update frequency, responsiveness
- `vulnerability` — known CVE exposure
- `license` — license permissiveness

### Common Alert Types
- `obfuscatedFile` — package contains likely obfuscated code
- `installScripts` — package runs scripts during install
- `networkAccess` — unexpected network calls
- `shellAccess` — shell command execution
- `envVars` — reads environment variables
- `filesystemAccess` — unexpected file I/O

## Workflow Guidance

When investigating alerts:
1. Start with `socket alerts list --org <slug> --per-page 10` for an overview
2. Note the affected repos and packages
3. Use `socket packages lookup` to get scores for concerning packages
4. Cross-reference with `socket fullscans get` for full context
5. Summarize findings with severity, affected repos, and recommended actions

When reviewing a repository's security posture:
1. `socket repos get --org <slug> --repo <repo-slug>` for repo details
2. Use the `head_full_scan_id` to get the latest scan: `socket fullscans get --org <slug> --id <id>`
3. `socket alerts list --org <slug>` and filter by repo in the response

When looking up a specific package:
1. `socket packages lookup pkg:<ecosystem>/<name>@<version> --org <slug>`
2. Interpret scores and flag anything below 0.7 as worth investigating
