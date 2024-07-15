//go:build linux
// +build linux

package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net/http"
	"os"
)

// UploadHandler handles file upload and saves it using sendfile
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// 10mb以上存临时文件，否则存储内存。
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a destination file
	dstPath := handler.Filename
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("Error creating destination file: %v", err)
		http.Error(w, "Error creating destination file", http.StatusInternalServerError)
		return
	}
	defer dstFile.Close()

	// Get file descriptors
	srcFd := int(file.(*os.File).Fd())
	dstFd := int(dstFile.Fd())

	// Get the file size
	fi, err := file.(*os.File).Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}
	fileSize := fi.Size()

	// Use sendfile to transfer data
	_, err = unix.Sendfile(dstFd, srcFd, nil, int(fileSize))
	if err != nil {
		log.Printf("Error using sendfile: %v", err)
		http.Error(w, "Error using sendfile", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func main() {
	http.HandleFunc("/upload", UploadHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
