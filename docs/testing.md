# Testing Guide

## Run all tests

Run the full Go test suite from the repository root:

```bash
go test ./...
```

## Current automated coverage

The server tests in `sha256/server/server_test.go` exercise the upload handler directly with `httptest`:

- Rejects non-`POST` requests with `405 Method Not Allowed`.
- Rejects requests that do not include the required multipart `file` field with `400 Bad Request`.
- Accepts a valid multipart upload, writes the file to a temporary upload directory, and returns the SHA-256 digest of the uploaded bytes.

The tests use a temporary working directory for successful uploads, so they do not modify the repository's real `sha256/server/data/uploads/` directory.

## Manual smoke test

Use this when validating the end-to-end TLS flow:

1. Start the server:

   ```bash
   cd sha256/server
   go run server.go
   ```

2. In another terminal, run the client:

   ```bash
   cd sha256/client
   go run client.go
   ```

3. Confirm the client prints `File uploaded successfully` and a 64-character SHA-256 digest.

## Suggested future tests

- Client request construction tests after extracting reusable client functions from `main`.
- Filename sanitization and path traversal tests before accepting untrusted filenames in production.
- Upload size limit tests after adding maximum body size enforcement.
- Integration tests that start an HTTPS test server with a generated certificate.
