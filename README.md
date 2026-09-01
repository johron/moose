# Moose
Terminal-based text editor written in Go, inspired by Emacs, Vim and VS Code.

## In-dev screenshot
<img width="480" height="262" alt="image" src="https://github.com/user-attachments/assets/99b80245-5044-4a30-97db-731d5dca6a67" />

## Features (+Wanted features)
- Modal editing
- Advanced layout with splitting and workspaces
- Undo and redo (currently no grouping, single changes only atm.)
- Lua-based extensions
- Command palette for commands, find (basically unimplemented), replace (unimplemented)
- Quick and easy keyboard- and command-oriented text manipulation

## Requirements
- Go 1.27
- Make
- (See [flake.nix](./flake.nix))

## Building
Clone repo and build:
```bash
git clone https://github.com/johron/moose.git
cd moose
# Activate flake if nix
make release
```

Debug build: `build/moose-debug`
Release build: `build/moose`

### Build targets
* `make build`: debug binary `build/moose-debug`
* `make release`: release binary `build/moose`
* `make run`: run debug version through `go run`
* `make clean`: remove `build/` directory
