# Secure File Transfer (SHA-256)

This project demonstrates a simple secure file transfer flow in Go using HTTPS and SHA-256 file integrity hashing.

## What the project does
- Accepts file uploads over TLS (`https://localhost:8080/upload`).
- Saves uploaded files to `sha256/server/data/uploads/`.
- Computes a SHA-256 hash while streaming the file to disk.
- Returns the SHA-256 digest to the client after upload.

## Project structure
- `sha256/server/server.go` – HTTPS upload server and SHA-256 computation.
- `sha256/client/client.go` – HTTPS client that uploads `sha256/client/data/input-file.txt`.
- `sha256/certs/` – local TLS certificate and key for localhost development.

## How to run

### 1) Start the server
```bash
cd sha256/server
go run server.go
```

### 2) Run the client in another terminal
```bash
cd sha256/client
go run client.go
```

### 3) Expected result
The client prints a success message and SHA-256 value returned by the server, for example:
```text
File uploaded successfully
SHA256: <64-character-hex-digest>
```

## Notes
- This setup is intended for local development/demo usage.
- The client currently trusts the provided local certificate in `sha256/certs/server.crt`.
