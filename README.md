# Zellijinator

A powerful session manager for [Zellij](https://zellij.dev/) - inspired by [tmuxinator](https://github.com/tmuxinator/tmuxinator) but designed specifically for Zellij's unique features.

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

Zellijinator helps you manage complex Zellij sessions with pre-configured layouts, making it easy to start development environments with a single command. Define your project's layout in YAML, and Zellijinator will create a Zellij session with all your tabs, panes, and commands ready to go.

## Features

- 🚀 **Quick Start**: Launch complex Zellij sessions with `zellijinator project-name`
- 📝 **YAML Configuration**: Simple, readable project definitions
- 🎨 **Beautiful CLI**: Styled output using [Charm](https://charm.sh) libraries
- 🔧 **Flexible Layouts**: Support for custom and predefined pane layouts
- 📁 **Project Management**: Create, edit, list, and delete projects easily
- 🎯 **Interactive Selection**: Choose projects interactively when no name is provided
- 🔄 **Session Handling**: Intelligent session management with automatic attachment to existing sessions

## Installation

### Homebrew (macOS and Linux)

```bash
brew tap dphaener/tap
brew install zellijinator
```

### Download Binary

Download the latest release from the [releases page](https://github.com/dphaener/zellijinator/releases).

### From Source

```bash
git clone https://github.com/dphaener/zellijinator.git
cd zellijinator
go build -o zellijinator
sudo mv zellijinator /usr/local/bin/
```

### Prerequisites

- [Zellij](https://zellij.dev/) must be installed
- Go 1.23+ (for building from source)

## Quick Start

1. Create a new project:
```bash
zellijinator new myproject
```

2. Start the project:
```bash
zellijinator myproject
# or just
zellijinator  # for interactive selection
```

3. List all projects:
```bash
zellijinator list
```

## Usage

### Commands

- `zellijinator [project]` - Start a project (default command)
- `zellijinator new <project>` - Create a new project configuration
- `zellijinator edit [project]` - Edit an existing project
- `zellijinator list` - List all projects
- `zellijinator delete [project]` - Delete a project

### Configuration

Project configurations are stored in `~/.zellijinator/` as YAML files. Here's an example configuration:

```yaml
# Project name (required)
name: myproject

# Root directory for the project (supports ~ and environment variables)
root: ~/projects/myapp

# Session name (optional, defaults to project name)
session_name: myapp-dev

# Default layout preset for tabs without explicit layout
default_layout: even-horizontal

# Environment variables
env:
  NODE_ENV: development
  DATABASE_URL: postgresql://localhost/myapp_dev

# Tabs configuration
tabs:
  - name: editor
    root: .  # Relative to project root
    layout: even-horizontal
    panes:
      - commands:
          - nvim

  - name: server
    # Using predefined layout
    layout: even-vertical
    panes:
      - name: backend
        commands:
          - npm run dev
      - name: frontend
        root: ./frontend
        commands:
          - npm run dev

  - name: terminals
    # Custom layout with nested panes
    panes:
      - commands:
          - echo "Main pane"
      - split_direction: vertical
        panes:
          - commands:
              - echo "Top right"
          - commands:
              - echo "Bottom right"
```

### Predefined Layouts

Zellijinator supports several predefined layouts for common pane arrangements:

- `even-horizontal` - Panes split evenly horizontally
- `even-vertical` - Panes split evenly vertically  
- `main-vertical` - Large pane on the left, others stacked vertically on the right
- `main-horizontal` - Large pane on top, others split horizontally below
- `tiled` - Panes arranged in a grid

### Environment Variables

You can use environment variables in your configuration:

```yaml
root: $HOME/projects/$PROJECT_NAME
env:
  API_KEY: $SECRET_API_KEY
```

## Advanced Features

### Custom Zellij Layouts

If you need more control, you can specify a custom KDL layout file:

```yaml
name: custom-project
layout: ~/.config/zellij/layouts/custom.kdl
```

### Focus Control

Set which pane should be focused when the session starts:

```yaml
tabs:
  - name: main
    panes:
      - commands: ["htop"]
      - focus: true  # This pane will be focused
        commands: ["nvim"]
```

### Compact Mode

Use Zellij's compact mode to save screen space:

```yaml
layout: compact
```

## Tips

1. **Default Command**: Since `start` is the default command, you can launch projects with just `zellijinator myproject`

2. **Interactive Mode**: Run `zellijinator` without arguments to interactively select a project

3. **Session Names**: If not specified, the session name defaults to the project name

4. **Relative Paths**: Paths in tab/pane configurations are relative to the project root

5. **Shell Integration**: Commands run in your default shell (`$SHELL`)

## Troubleshooting

### Session Already Exists

Zellijinator will automatically attach to existing sessions. Dead sessions are resurrected when you try to start them.

### Commands Not Running

Ensure your commands are valid and that any required dependencies are installed. Commands run in your default shell.

### Layout Issues

If panes aren't arranged as expected, check that you're using the correct layout syntax and that nested panes are properly structured.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [tmuxinator](https://github.com/tmuxinator/tmuxinator)
- Built with [Cobra](https://github.com/spf13/cobra) for CLI management
- Styled with [Charm](https://charm.sh) libraries for beautiful terminal output
- Designed for [Zellij](https://zellij.dev/) terminal multiplexer# Triggering auto-release
