package pkg

import (
	"testing"
)

func TestStoreMultipleEntriesInCSV(t *testing.T) {
	entries := []FileDetails{
		{
			Filename: "testfile1.txt",
			FileSize: 1234,
			FileHash: "abcd1234",
		},
		{
			Filename: "testfile.txt",
			FileSize: 4567,
			FileHash: "aeyd4567",
		},
		{
			Filename: "testfile2.txt",
			FileSize: 5678,
			FileHash: "efgh5678",
		},
		// Add more entries as needed
	}

	for _, details := range entries {
		err := storeInCSV(details)
		if err != nil {
			t.Errorf("storeInCSV failed with error: %v", err)
		}
	}
}

func TestDeleteFromCSV(t *testing.T) {
	fileName := "testfile.txt"

	err := deleteFromCSV(fileName)
	if err != nil {
		t.Errorf("deleteFromCSV failed with error: %v", err)
	}
}

func TestGetAllEntries(t *testing.T) {
	_, err := getAllEntries()
	if err != nil {
		t.Errorf("getAllEntries failed with error: %v", err)
	}
}

func TestUpdateInCSV(t *testing.T) {
	fileName := "testfile1.txt"
	newDetails := FileDetails{
		Filename: "newfile.txt",
		FileSize: 1212,
		FileHash: "updated",
	}

	err := updateInCSV(fileName, newDetails)
	if err != nil {
		t.Errorf("updateInCSV failed with error: %v", err)
	}
}

func TestFindByHash(t *testing.T) {
	fileHash := "abcd1234"
	entry, err := findByHash(fileHash)
	if err != nil {
		t.Errorf("findByHash failed with error: %v", err)
	}
	if entry == nil {
		t.Errorf("findByHash failed to find the entry")
	}
}

func TestCleanCSV(t *testing.T) {
	err := cleanCSV()
	if err != nil {
		t.Errorf("cleanCSV failed with error: %v", err)
	}
}
