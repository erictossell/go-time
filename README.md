## Go-Time

Another CLI tool for time tracking. This one is written in Go.

Track your time and categorize it with tags. You can start and stop timers, or use the TUI to manage your time entries.

#### Built with

- [Cobra](https://github.com/spf13/cobra)
- [Bubbletea](https://github.com/charmbracelet/bubbletea/)
- [Huh?](https://github.com/charmbracelet/huh/)
- SQLite and [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

### Usage

```bash
  go-time [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  del         Delete an existing time entry
  edit        Edit an existing time entry
  help        Help about any command
  read        List all active timers or time entries
  start       Start a new timer with optional tags
  stop        Stop the current timer and add tags
  tui         Launch the Text-based User Interface

Flags:
  -h, --help   help for go-time

Use "go-time [command] --help" for more information about a command.

```

### NixOS Flakes Installation

In `flake.nix` inputs add:

```nix
inputs = {
  go-time.url = "github:erictossell/go-time";
}; 
```

Import as a `module.nix`:

```nix
{ pkgs, go-time, ... }: 
{
  environment.systemPackages = with pkgs; [
    go-time.packages.${system}.default
  ];
}
```

