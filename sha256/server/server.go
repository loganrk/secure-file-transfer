package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("HTTPS Server running on https://localhost:8080")
	err := http.ListenAndServeTLS(":8080", "../certs/server.crt", "../certs/server.key", nil)
	if err != nil {
		panic(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dstPath := "data/uploads/" + handler.Filename
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	hash := sha256.New()
	multiWriter := io.MultiWriter(dst, hash)

	_, err = io.Copy(multiWriter, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	hashString := hex.EncodeToString(hash.Sum(nil))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully\nSHA256: " + hashString))
}
