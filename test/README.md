# Test Plan

The repository includes automated Go tests and a manual smoke test for the local HTTPS transfer flow.

## Automated tests

Run from the repository root:

```bash
go test ./...
```

Covered behavior:

- `GET /upload` is rejected because only `POST` uploads are supported.
- Malformed or incomplete upload requests are rejected when the `file` multipart field is missing.
- A valid multipart upload is saved and the response includes the SHA-256 digest for the uploaded bytes.

## Manual smoke test

1. Start the server from `sha256/server` with `go run server.go`.
2. Run the client from `sha256/client` with `go run client.go`.
3. Verify that the client prints `File uploaded successfully` and a 64-character `SHA256:` value.
4. Verify that a timestamped file appears in `sha256/server/data/uploads/`.

See `docs/testing.md` for the detailed testing guide.
