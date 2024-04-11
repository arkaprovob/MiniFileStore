package pkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

// checkErrorAndRespond is a helper function that checks if an error occurred (err is not nil).
// If an error occurred, it logs the message and responds with the given code and message.
func checkErrorAndRespond(err error, message string, code int, w http.ResponseWriter) {
	if err != nil {
		log.Println(message, err)
		http.Error(w, message, code)
		return
	}
}

func validateRequiredField(fieldName, fieldValue string) error {
	if fieldValue == "" {
		log.Printf("%s is required", fieldName)
		return fmt.Errorf("field %s is missing", fieldName)
	}
	return nil
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

// todo get teh filepath from the environment variable `os.Getenv("FILES_DIR")`
func getFilePath(filename string) string {
	return filepath.Join("files", filename)
}

func ManageFileUpdate(duplicate bool, newFileName string, previousFileDetails FileDetails) error {

	// if duplicate is true, then duplicate an existing file with the newFileName
	if duplicate {
		err := DuplicateFile(previousFileDetails.Filename, newFileName)
		if err != nil {
			return err
		}

		newFileDetails := FileDetails{
			Filename: newFileName,
			FileSize: previousFileDetails.FileSize,
			FileHash: previousFileDetails.FileHash,
		}

		err = storeInCSV(newFileDetails)
		if err != nil {
			return err
		}

		return nil
	}

	// if duplicate is false, then update the existing file with the newFileName
	err := UpdateFileName(previousFileDetails.Filename, newFileName)
	if err != nil {
		return err
	}

	return nil
}

func UpdateFileName(prevFilename string, newName string) error {
	// Rename the file to the new name
	err := os.Rename(getFilePath(prevFilename), getFilePath(newName))
	if err != nil {
		return err
	}
	return nil
}

func DuplicateFile(src string, dst string) error {
	// Open the source file for reading
	srcFile, err := os.Open(getFilePath(src))
	if err != nil {
		log.Println("Error opening the source file:", err)
		return err
	}
	defer CloseFile(srcFile)

	// Create the destination file
	dstFile, err := os.Create(getFilePath(dst))
	if err != nil {
		log.Println("Error creating the destination file:", err)
		return err
	}
	defer CloseFile(dstFile)

	// Copy the contents from the source file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Println("Error copying the file:", err)
		return err
	}

	// Ensure the data is written to the disk
	err = dstFile.Sync()
	if err != nil {
		log.Println("Error syncing the file:", err)
		return err
	}

	return nil
}

func modifyRecordAndFile(oldRecord FileDetails, newRecord FileDetails) error {
	err := deleteFromCSV(oldRecord.Filename)
	if err != nil {
		log.Println("Error deleting the record from the CSV:", err)
		return err
	}
	err = storeInCSV(newRecord)
	if err != nil {
		log.Println("Error storing the record in the CSV:", err)
		return err
	}
	err = deleteFile(oldRecord.Filename)
	if err != nil {
		log.Println("Error deleting the file:", err)
		return err
	}
	return nil
}

func deleteFile(filename string) error {
	err := os.Remove(getFilePath(filename))
	if err != nil {
		log.Println("Error deleting the file:", err)
		return err
	}
	return nil
}
