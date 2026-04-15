package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {

	filePath := "data/input-file.txt"

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// ✅ Extract original filename
	originalName := filepath.Base(filePath)

	ext := filepath.Ext(originalName)
	name := originalName[:len(originalName)-len(ext)]

	// ✅ Add timestamp
	fileName := fmt.Sprintf("%s-%s%s",
		name,
		time.Now().Format("20060102-150405"),
		ext,
	)

	// ✅ FIXED syntax
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	writer.Close()

	// TLS config
	certPool := x509.NewCertPool()
	certData, err := os.ReadFile("../certs/server.crt")
	if err != nil {
		panic(err)
	}

	certPool.AppendCertsFromPEM(certData)

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", "https://localhost:8080/upload", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	fmt.Println(string(respBody))
}
