# manage-agent-skills

A command-line tool to simplify the installation and management of agent skills from GitHub repositories.

## Background

- Installing agent skills is cumbersome (requires git clone for each skill and creating symbolic links)
- Reusing skills across multiple agents is tedious (requires creating symbolic links in each agent's directory)

## Purpose

Simplify the installation and management of agent skills.

## Features

- **download**: Clone a GitHub repository (e.g., `tsubasaogawa/semantic-commit-helper`) to `~/.local/bin/manage-agent-skills`
- **install <agent name>**: Create symbolic links to cloned repositories in the specified agent's skill directory
  - Agent skill directory list is managed in `config.toml`
- **list**: Display a list of skills cloned under `~/.local/bin/manage-agent-skills`

## Installation

### Prerequisites

- Go 1.22 or later
- Git

### Build from source

```bash
git clone https://github.com/tsubasaogawa/manage-agent-skills.git
cd manage-agent-skills
go build -o manage-agent-skills .
```

Optionally, move the binary to your PATH:

```bash
sudo mv manage-agent-skills /usr/local/bin/
```

Or install with `go install`:

```bash
go install github.com/tsubasaogawa/manage-agent-skills@latest
```

## Configuration

Create a configuration file at `~/.config/manage-agent-skills/config.toml`:

```toml
[agents]
my-agent = "~/.local/share/my-agent/skills"
another-agent = "/opt/agents/another-agent/skills"
```

Each entry maps an agent name to its skills directory path.

## Usage

### Download a skill

Download a skill from a GitHub repository:

```bash
manage-agent-skills download tsubasaogawa/semantic-commit-helper
```

This clones the repository to `~/.local/bin/manage-agent-skills/semantic-commit-helper`.

### List downloaded skills

Display all downloaded skills:

```bash
manage-agent-skills list
```

Output:
```
Downloaded skills:
  - semantic-commit-helper
  - another-skill
```

### Install skills to an agent

Install all downloaded skills to a specific agent:

```bash
manage-agent-skills install my-agent
```

This creates symbolic links in the agent's skill directory (as defined in `config.toml`) for all downloaded skills.

Output:
```
  Installed semantic-commit-helper
  Installed another-skill
Successfully installed skills to agent 'my-agent'
```

### Help

Get help for any command:

```bash
manage-agent-skills --help
manage-agent-skills download --help
```

## Directory Structure

```
~/.local/bin/manage-agent-skills/
├── semantic-commit-helper/   # Downloaded skill repository
├── another-skill/             # Downloaded skill repository
└── ...

~/.config/manage-agent-skills/
└── config.toml                # Configuration file

~/.local/share/my-agent/skills/
├── semantic-commit-helper -> ~/.local/bin/manage-agent-skills/semantic-commit-helper
├── another-skill -> ~/.local/bin/manage-agent-skills/another-skill
└── ...
```

## Development

### Run tests

```bash
go test ./...
```

### Build

```bash
go build -o manage-agent-skills .
```

## License

See [LICENSE](LICENSE) file.