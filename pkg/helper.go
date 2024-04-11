package pkg

import (
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
