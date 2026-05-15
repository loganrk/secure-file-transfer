package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadHandlerRejectsNonPostMethods(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	rec := httptest.NewRecorder()

	uploadHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "Only POST allowed") {
		t.Fatalf("expected method error response, got %q", rec.Body.String())
	}
}

func TestUploadHandlerRequiresFileField(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader("not multipart"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	uploadHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "File error") {
		t.Fatalf("expected file error response, got %q", rec.Body.String())
	}
}

func TestUploadHandlerSavesFileAndReturnsSHA256(t *testing.T) {
	withWorkingUploadDir(t)

	content := []byte("secure transfer test payload\n")
	fileName := "payload.txt"
	body, contentType := multipartUploadBody(t, "file", fileName, content)
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	uploadHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %q", http.StatusOK, rec.Code, rec.Body.String())
	}

	expectedHashBytes := sha256.Sum256(content)
	expectedHash := hex.EncodeToString(expectedHashBytes[:])
	if !strings.Contains(rec.Body.String(), expectedHash) {
		t.Fatalf("expected response to contain hash %s, got %q", expectedHash, rec.Body.String())
	}

	savedPath := filepath.Join("data", "uploads", fileName)
	savedContent, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("expected uploaded file at %s: %v", savedPath, err)
	}

	if !bytes.Equal(savedContent, content) {
		t.Fatalf("saved file content mismatch: got %q, want %q", savedContent, content)
	}
}

func multipartUploadBody(t *testing.T, fieldName, fileName string, content []byte) (io.Reader, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("create multipart file field: %v", err)
	}

	if _, err := part.Write(content); err != nil {
		t.Fatalf("write multipart content: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	return body, writer.FormDataContentType()
}

func withWorkingUploadDir(t *testing.T) {
	t.Helper()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempDir, "data", "uploads"), 0o755); err != nil {
		t.Fatalf("create upload directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change to temp directory: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})
}
