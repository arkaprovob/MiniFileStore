package pkg

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
func getFileStorePath(filename string) (string, error) {
	config, err := GetConfig()
	if err != nil {
		log.Println("Error getting the config:", err)
		return "", err
	}
	return filepath.Join(config.FileStore, filename), nil
}

func getFileStoreDir() (string, error) {
	config, err := GetConfig()
	if err != nil {
		log.Println("Error getting the config:", err)
		return "", err
	}
	return config.FileStore, nil
}

func RecordStorePath(filename string) (string, error) {
	config, err := GetConfig()
	if err != nil {
		log.Println("Error getting the config:", err)
		return "", err
	}
	return filepath.Join(config.RecordStore, filename), nil
}

func ManageFileUpdate(duplicate bool, newFileName string, previousFileDetails FileDetails) error {

	newFileDetails := FileDetails{
		Filename:  newFileName,
		FileSize:  previousFileDetails.FileSize,
		FileHash:  previousFileDetails.FileHash,
		WordCount: previousFileDetails.WordCount,
	}

	// if duplicate is true, then duplicate an existing file with the newFileName
	if duplicate {
		log.Println("Duplicating the file")
		err := DuplicateFile(previousFileDetails.Filename, newFileName)
		if err != nil {
			return err
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
	err = updateInCSV(previousFileDetails.Filename, newFileDetails)
	if err != nil {
		log.Println("Error updating the record in the CSV:", err)
		return err
	}
	return nil
}

func UpdateFileName(prevFilename string, newName string) error {

	previousFile, err := getFileStorePath(prevFilename)
	if err != nil {
		return err
	}
	newFile, err := getFileStorePath(newName)
	if err != nil {
		return err
	}

	// Rename the file to the new name
	err = os.Rename(previousFile, newFile)
	if err != nil {
		return err
	}
	return nil
}

func DuplicateFile(src string, dst string) error {
	// Open the source file for reading
	sourceFile, err := getFileStorePath(src)
	if err != nil {
		log.Println("Error finding the path of the source file:", err)
		return err
	}
	destFile, err := getFileStorePath(dst)
	if err != nil {
		log.Println("Error finding the path of the destination file:", err)
		return err
	}

	srcFile, err := os.Open(sourceFile)
	if err != nil {
		log.Println("Error opening the source file:", err)
		return err
	}
	defer CloseFile(srcFile)

	// Create the destination file
	dstFile, err := os.Create(destFile)
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
	// Check if the old file name is different from the new file name
	if oldRecord.Filename != newRecord.Filename {
		// If the file names are different, delete the old file
		err = deleteFile(oldRecord.Filename)
		if err != nil {
			log.Println("Error deleting the file:", err)
			return err
		}
	}
	return nil
}

func deleteFile(filename string) error {

	filePath, err := getFileStorePath(filename)
	if err != nil {
		log.Println("Error finding the path of the file:", err)
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		log.Println("Error deleting the file:", err)
		return err
	}
	return nil
}

func countWordsInFile(fileLocation string) (int, error) {

	filePath, err := getFileStorePath(fileLocation)
	if err != nil {
		log.Println("Error finding the path of the file:", err)
		return 0, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer CloseFile(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	wordCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		wordCount += len(words)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return wordCount, nil
}
