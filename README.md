# Secure File Transfer (SHA-256)

This project demonstrates a small secure file transfer workflow in Go. A client uploads a file to a local HTTPS server, the server streams the upload to disk, and the server returns the SHA-256 digest of the saved content so the client can verify file integrity.

## Features

- Accepts file uploads over TLS at `https://localhost:8080/upload`.
- Stores uploads under `sha256/server/data/uploads/` when the server is run from `sha256/server`.
- Computes SHA-256 while streaming uploaded content to disk.
- Returns the SHA-256 digest in the HTTP response.
- Includes automated handler tests for upload success and common error paths.

## Project structure

```text
.
├── docs/
│   ├── README.md
│   ├── testing.md
│   └── usage.md
├── sha256/
│   ├── certs/          # Local development TLS certificate and key
│   ├── client/         # Multipart HTTPS upload client
│   └── server/         # HTTPS upload server and tests
└── test/README.md      # Test plan summary
```

## Requirements

- Go 1.22 or newer.
- Localhost access to port `8080`.
- The included development certificate files in `sha256/certs/`.

> The bundled certificate and key are for local development only. Do not use them in production.

## Run the demo

### 1. Start the server

```bash
cd sha256/server
go run server.go
```

The server listens on `https://localhost:8080` and registers the `/upload` endpoint.

### 2. Run the client in another terminal

```bash
cd sha256/client
go run client.go
```

The client uploads `sha256/client/data/input-file.txt` using the local certificate authority bundle at `sha256/certs/server.crt`.

### 3. Expected output

```text
File uploaded successfully
SHA256: <64-character-hex-digest>
```

A timestamped copy of the uploaded file is written to `sha256/server/data/uploads/`.

## Run tests

From the repository root:

```bash
go test ./...
```

The current automated tests exercise the server upload handler without starting a real TLS listener. See [Testing Guide](docs/testing.md) for details.

## Documentation

- [Usage Guide](docs/usage.md) explains setup, running the server/client, and troubleshooting.
- [Testing Guide](docs/testing.md) describes the test suite and how to add coverage.
