package pkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
)

// CloseFile closes the given file and logs an error if one occurs.
func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Printf("Error closing file: %v", err)
	}
}

// CloseMultipartFile closes the given multipart file and logs an error if one occurs.
func CloseMultipartFile(file multipart.File) {
	err := file.Close()
	if err != nil {
		log.Printf("Error closing multipart file: %v", err)
	}
}

// ComputeMD5Hash reads a file and computes its MD5 hash.
// alternate approach use linux shell utility but since this method is uses a streaming approach
// to read the file and compute the MD5 hash,t doesn't need to load the entire file into memory at once.
func ComputeMD5Hash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer CloseFile(file)

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	return fmt.Sprintf("%x", hashInBytes), nil
}
