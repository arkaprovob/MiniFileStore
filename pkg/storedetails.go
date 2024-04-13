package pkg

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

var CsvFileLocation = "fileDetails.csv"
var TempCsvFileLocation = "fileDetailsTemp.csv"

type FileDetails struct {
	Filename string
	FileSize int64
	FileHash string
}

func storeInCSV(details FileDetails) error {
	file, err := os.OpenFile(CsvFileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening the file:", err)
		return err
	}
	defer CloseFile(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{details.Filename, strconv.FormatInt(details.FileSize, 10), details.FileHash}
	err = writer.Write(record)
	if err != nil {
		log.Println("Error writing record to the file:", err)
		return err
	}

	return nil
}

func deleteFromCSV(fileName string) error {
	file, err := os.Open(CsvFileLocation)
	if err != nil {
		log.Println("Error opening the file:", err)
		return err
	}
	defer CloseFile(file)

	temp, err := os.Create(TempCsvFileLocation)
	if err != nil {
		log.Println("Error creating the file:", err)
		return err
	}
	defer CloseFile(temp)

	reader := csv.NewReader(file)
	writer := csv.NewWriter(temp)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading the file:", err)
			return err
		}

		if record[0] != fileName {
			err = writer.Write(record)
			if err != nil {
				log.Println("Error writing record to the file:", err)
				return err
			}
		}
	}

	err = os.Remove(CsvFileLocation)
	if err != nil {
		log.Println("Error removing the file:", err)
		return err
	}

	err = os.Rename(TempCsvFileLocation, CsvFileLocation)
	if err != nil {
		log.Println("Error renaming the file:", err)
		return err
	}

	return nil
}

func getAllEntries() ([]FileDetails, error) {
	file, err := os.Open(CsvFileLocation)
	if err != nil {
		log.Println("Error opening the file:", err)
		return nil, err
	}
	defer CloseFile(file)

	reader := csv.NewReader(file)
	var entries []FileDetails

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading the file:", err)
			return nil, err
		}

		fileSize, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Println("Error parsing the file size:", err)
			return nil, err
		}

		entries = append(entries, FileDetails{
			Filename: record[0],
			FileSize: fileSize,
			FileHash: record[2],
		})
	}

	return entries, nil
}
func updateInCSV(fileName string, newDetails FileDetails) error {
	file, err := os.Open(CsvFileLocation)
	if err != nil {
		log.Println("Error opening the file:", err)
		return err
	}
	defer CloseFile(file)

	temp, err := os.Create(TempCsvFileLocation)
	if err != nil {
		log.Println("Error creating the file:", err)
		return err
	}
	defer CloseFile(temp)

	reader := csv.NewReader(file)
	writer := csv.NewWriter(temp)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading the file:", err)
			return err
		}

		if record[0] == fileName {
			record[0] = newDetails.Filename
			record[1] = strconv.FormatInt(newDetails.FileSize, 10)
			record[2] = newDetails.FileHash
		}

		err = writer.Write(record)
		if err != nil {
			log.Println("Error writing to the file:", err)
			return err
		}
	}

	err = os.Remove(CsvFileLocation)
	if err != nil {
		log.Println("Error removing the file:", err)
		return err
	}

	err = os.Rename(TempCsvFileLocation, CsvFileLocation)
	if err != nil {
		log.Println("Error renaming the file:", err)
		return err
	}

	return nil
}

// TODO: Improve the findByHash function by implementing a hashmap for O(1) lookup time.
func findByHash(hash string) (*FileDetails, error) {
	entries, err := getAllEntries()
	if err != nil {
		log.Println("Error getting all entries:", err)
		return nil, err
	}

	for _, entry := range entries {
		if entry.FileHash == hash {
			return &entry, nil
		}
	}

	return nil, nil
}

// TODO: Improve the findByName function by implementing a hashmap for O(1) lookup time.
func findByName(name string) (*FileDetails, error) {
	entries, err := getAllEntries()
	if err != nil {
		log.Println("Error getting all entries:", err)
		return nil, err
	}

	for _, entry := range entries {
		if entry.Filename == name {
			return &entry, nil
		}
	}

	return nil, nil
}

func findByHashOrName(hash string, name string) (*FileDetails, error) {
	// First, try to find by hash
	record, err := findByHash(hash)
	if err != nil {
		log.Println("Error finding the hash:", err)
		return nil, err
	}

	// If a record is found by hash, return it
	if record != nil {
		return record, nil
	}

	// If no record is found by hash, try to find by name
	record, err = findByName(name)
	if err != nil {
		log.Println("Error finding the name:", err)
		return nil, err
	}

	// Return the record found by name (could be nil if no record is found)
	return record, nil
}

func cleanCSV() error {
	file, err := os.Create(CsvFileLocation)
	if err != nil {
		log.Println("Error creating the file:", err)
		return err
	}
	defer CloseFile(file)

	return nil
}
