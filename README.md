# socket-cli

A Go CLI for the [Socket.dev](https://socket.dev) supply chain security platform.

Unlike the official Socket CLI (written in TypeScript with a large dependency tree), this is a single statically-linked binary with zero runtime dependencies.

## Installation

```bash
go install github.com/scottbrown/socket-cli/cmd/socket@latest
```

Or build from source:

```bash
git clone https://github.com/scottbrown/socket-cli.git
cd socket-cli
task build
# Binary is at .build/socket
```

## Authentication

The CLI needs a Socket.dev API token. Get one from the [Socket dashboard](https://socket.dev/dashboard).

Set it via environment variable (preferred):

```bash
export SOCKET_API_TOKEN=sktsec_...
```

Or write it to a file:

```bash
mkdir -p ~/.config/socket
echo "sktsec_..." > ~/.config/socket/token
chmod 600 ~/.config/socket/token
```

## Usage

```
socket [command]

Available Commands:
  orgs        Manage organizations
  repos       Manage repositories
  fullscans   Manage full scans
  diffscans   Manage diff scans
  alerts      View and manage alerts
  packages    Query package information
  quota       Show API quota information
```

### Examples

```bash
# List your organizations
socket orgs list

# List repos in an org
socket repos list --org my-org

# View supply chain alerts
socket alerts list --org my-org --per-page 10

# Look up a package by purl
socket packages lookup pkg:npm/express@4.18.2

# Upload manifests for a full scan
socket fullscans create --org my-org --repo my-repo --branch main \
  --file package.json --file package-lock.json

# Check API quota
socket quota
```

### Commands

| Command | Description |
|---------|-------------|
| `socket orgs list` | List organizations linked to your token |
| `socket repos list --org <slug>` | List repositories |
| `socket repos get --org <slug> --repo <repo>` | Get repo details |
| `socket repos delete --org <slug> --repo <repo>` | Delete a repository |
| `socket fullscans list --org <slug>` | List full scans (with filtering) |
| `socket fullscans get --org <slug> --id <id>` | Get scan results |
| `socket fullscans create --org <slug> --repo <name> --file <manifest>` | Upload manifests for scanning |
| `socket fullscans delete --org <slug> --id <id>` | Delete a scan |
| `socket fullscans metadata --org <slug> --id <id>` | Get scan metadata |
| `socket diffscans list --org <slug>` | List diff scans |
| `socket diffscans get --org <slug> --id <id>` | Get a diff scan |
| `socket diffscans delete --org <slug> --id <id>` | Delete a diff scan |
| `socket alerts list --org <slug>` | List alerts |
| `socket alerts triage --org <slug>` | List triaged alerts |
| `socket packages lookup <purl> [--org <slug>]` | Look up packages by purl |
| `socket quota` | Show remaining API quota |

## Development

Requires [Go](https://go.dev) 1.26+ and [Task](https://taskfile.dev).

```bash
task build      # Compile to .build/socket
task test       # Run unit tests
task coverage   # Generate coverage report (.test/coverage.html)
task vet        # Run go vet
task clean      # Remove build artifacts
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SOCKET_API_TOKEN` | API token for authentication | — |
| `SOCKET_API_URL` | Override the API base URL | `https://api.socket.dev/v0` |

## AI Skill

A Claude Code skill file is included at [`docs/skill.md`](docs/skill.md) for interactive use with the Socket CLI. Copy it to `~/.claude/skills/socket-cli/skill.md` to enable the `/socket-cli` slash command.

## License

See [LICENSE](LICENSE).
