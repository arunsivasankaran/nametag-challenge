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

### 1. Build the project

From the project root:

```bash
go build ./...
```

This compiles the CLI and any related Go files in the module.

### 2. Run the app

```bash
go run .
```

This starts the CLI and checks for a newer version from the configured update URL.

### 3. Run the mock update server

The app is configured to look for updates at:

```text
http://127.0.0.1:8080/manifest.json
```

To simulate the update service locally, run the mock server in a separate terminal:

```bash
go run .
```

However, since the app itself is the CLI and not the HTTP server, the project is currently intended for a demo flow where the update endpoint is served by a separate component. In a fuller implementation, you would typically run the server in a separate process or a dedicated release host.

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
