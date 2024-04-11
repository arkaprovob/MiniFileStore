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
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{details.Filename, strconv.FormatInt(details.FileSize, 10), details.FileHash}
	err = writer.Write(record)
	if err != nil {
		return err
	}

	return nil
}

func deleteFromCSV(fileName string) error {
	file, err := os.Open(CsvFileLocation)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	temp, err := os.Create(TempCsvFileLocation)
	if err != nil {
		return err
	}
	defer func(temp *os.File) {
		err := temp.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(temp)

	reader := csv.NewReader(file)
	writer := csv.NewWriter(temp)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if record[0] != fileName {
			err = writer.Write(record)
			if err != nil {
				return err
			}
		}
	}

	err = os.Remove(CsvFileLocation)
	if err != nil {
		return err
	}

	err = os.Rename(TempCsvFileLocation, CsvFileLocation)
	if err != nil {
		return err
	}

	return nil
}

func getAllEntries() ([]FileDetails, error) {
	file, err := os.Open(CsvFileLocation)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	reader := csv.NewReader(file)
	var entries []FileDetails

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		fileSize, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
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
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	temp, err := os.Create(TempCsvFileLocation)
	if err != nil {
		return err
	}
	defer func(temp *os.File) {
		err := temp.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(temp)

	reader := csv.NewReader(file)
	writer := csv.NewWriter(temp)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if record[0] == fileName {
			record[0] = newDetails.Filename
			record[1] = strconv.FormatInt(newDetails.FileSize, 10)
			record[2] = newDetails.FileHash
		}

		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	err = os.Remove(CsvFileLocation)
	if err != nil {
		return err
	}

	err = os.Rename(TempCsvFileLocation, CsvFileLocation)
	if err != nil {
		return err
	}

	return nil
}
func cleanCSV() error {
	file, err := os.Create(CsvFileLocation)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	return nil
}
