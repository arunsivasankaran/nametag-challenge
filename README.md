# Nametag Self-Updating Go CLI

This project is a small Go application that demonstrates a self-updating CLI. At startup, it checks a remote update manifest, compares the current version to the newest published version, downloads the newer binary if needed, validates the binary checksum, and replaces the running executable.

The purpose of the project is to model the kind of update flow used by desktop applications and deployment tools: fetch new code, validate it, install it safely, and recover if anything fails.

## File structure

- `main.go` - the core CLI application and updater logic
  - version comparison
  - remote manifest fetch
  - checksum validation
  - download flow
  - binary replacement logic
  - basic rollback behavior
- `server.go` - a simple local mock HTTP server used to simulate an update source
  - `/manifest.json` returns a JSON version manifest
  - `/nametag.bin` serves the downloaded binary payload
- `version_test.go` - unit tests for version comparison and manifest validation
- `go.mod` - Go module definition
- `README.md` - project explanation and run instructions

## How it works

1. The app starts and prints its current version.
2. It requests a manifest from the configured update URL.
3. It compares the current binary version to the newest remote version.
4. If the remote version is newer:
   - downloads the replacement binary to a temp location
   - verifies its SHA256 checksum
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

### 2. Build the project for local development

If you just want a single local binary:

```bash
go build .
```

### 3. Run the app

```bash
go run .
```

This starts the CLI and checks for a newer version from GitHub Releases.

### 4. Run the mock update server

This project can also be paired with the mock release server in `server.go` for local testing scenarios. In a realistic setup, the app is expected to check a public GitHub repository instead of the local mock server.

## Testing

Run the automated tests:

```bash
go test ./...
```

This validates:

- version comparison logic
- GitHub release asset selection behavior
- Go tag format handling

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
Latest GitHub release: v1.1.0
Download matching asset for current OS/arch
Validate the downloaded payload
Backup old binary
Replace active binary
Restart or continue with new version
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
