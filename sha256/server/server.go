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
	fmt.Println("file upload request received")
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("File error")
		http.Error(w, "File error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dstPath := "data/uploads/" + handler.Filename
	dst, err := os.Create(dstPath)
	if err != nil {
		fmt.Println("Unable to save file")
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	hash := sha256.New()
	multiWriter := io.MultiWriter(dst, hash)

	_, err = io.Copy(multiWriter, file)
	if err != nil {
		fmt.Println("Error saving file")

		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	hashString := hex.EncodeToString(hash.Sum(nil))

	w.WriteHeader(http.StatusOK)
	fmt.Printf("File %s uploaded successfully with SHA256: %s\n", handler.Filename, hashString)
	w.Write([]byte("File uploaded successfully\nSHA256: " + hashString))
}
