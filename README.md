# Nametag Self-Updating Go CLI

NOTE: This project depends on released posted to https://github.com/arunsivasankaran/nametag-challenge (public repo)

This project is a small Go application that demonstrates a self-updating CLI. At startup, it checks the latest public GitHub release for a configured repository, compares the local version to the newest published tag, downloads the matching release asset for the current OS and architecture, and replaces the running executable when a newer version is available.

The purpose of the project is to model the kind of update flow used by desktop applications and deployment tools: fetch the newest release, choose the correct asset, install it safely, and recover if the replacement fails.

## File structure

- `main.go` - the core CLI application and updater logic
  - version comparison
  - GitHub release discovery
  - platform-specific asset selection
  - binary download and install flow
  - rollback behavior
- `build.sh` - builds binaries for Windows, Linux, and macOS
- `version_test.go` - unit tests for version comparison and release asset selection
- `go.mod` - Go module definition
- `README.md` - project overview and usage instructions

## How it works

1. The app starts and prints its current version.
2. It calls the GitHub Releases API for the configured repository.
3. It compares the current binary version to the newest remote tag.
4. If the remote version is newer:
   - selects the release asset that matches the current OS and architecture
   - downloads the replacement binary to a temp location
   - renames the current binary to a backup
   - swaps in the new binary
   - removes the backup after a successful install
5. If the install fails, it attempts to restore the previous binary.

## Build and run

### Prerequisites

- Go 1.22 or newer installed
- A terminal with access to the project directory

### 1. Build all platform binaries

From the project root, run:

```bash
chmod +x build.sh
./build.sh
```

This compiles the project for the following targets and writes the artifacts to the `dist/` directory:

- Windows amd64
- Linux amd64
- Linux arm64
- macOS amd64
- macOS arm64

Example output:

```text
dist/
  nametag-windows-amd64.exe
  nametag-linux-amd64
  nametag-linux-arm64
  nametag-darwin-amd64
  nametag-darwin-arm64
```

### 2. Build a local binary for development

```bash
go build .
```

### 3. Run the app

```bash
export GITHUB_OWNER=your-org
export GITHUB_REPO=your-repo
go run .
```

If the environment variables are absent, the app falls back to the default repository config in the code.

## Testing

Run the automated tests:

```bash
go test ./...
```

This validates:

- version comparison logic
- GitHub release tag parsing
- platform-specific asset selection

## Notes

This is a minimal but production-minded MVP. A real-world version would likely add:

- signed release checks
- stronger cross-platform install safety
- more robust rollback and restart behavior
- richer logging and telemetry
- a more explicit release asset naming convention

## Example update workflow

A typical flow for this project looks like this:

```text
Current binary: v1.0.0
Latest GitHub release: v1.1.0
Select asset for current OS/arch
Download new binary
Backup old binary
Replace active binary
Continue with updated version
```

This keeps the challenge focused on the core idea: a program that updates itself safely.

## Testing

Run the automated tests:

```bash
go test ./...
```

This validates:

- version comparison logic
- manifest validation logic

## Notes

This is a minimal but production-minded MVP. A real-world version would likely add:

- signed release checks
- stronger cross-platform installation logic
- more robust rollback and restart behavior
- a proper release server or artifact repository
- logging and telemetry

## Example update workflow

A typical flow for this project looks like this:

```text
Current binary: v1.0.0
Remote manifest: v1.1.0
Download new binary
Validate checksum
Backup old binary
Replace active binary
Restart or continue with new version
```

This keeps the challenge focused on the core idea: a program that updates itself safely.
