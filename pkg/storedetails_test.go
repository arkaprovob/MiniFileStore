package pkg

import (
	"testing"
)

func setup(t *testing.T) func() {
	TestStoreMultipleEntriesInCSV(t)
	return func() {
		TestCleanCSV(t)
		teardown()
	}
}

func TestStoreMultipleEntriesInCSV(t *testing.T) {

	entries := []FileDetails{
		{
			Filename:  "testfile1.txt",
			FileSize:  1234,
			FileHash:  "abcd1234",
			WordCount: 10,
		},
		{
			Filename:  "testfile.txt",
			FileSize:  4567,
			FileHash:  "aeyd4567",
			WordCount: 11,
		},
		{
			Filename:  "testfile2.txt",
			FileSize:  5678,
			FileHash:  "efgh5678",
			WordCount: 12,
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

	teardown := fileStoreSetup(t)
	defer teardown()

	fileName := "testfile.txt"

	err := deleteFromCSV(fileName)
	if err != nil {
		t.Errorf("deleteFromCSV failed with error: %v", err)
	}
}

func TestGetAllEntries(t *testing.T) {
	teardown := setup(t)
	defer teardown()
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

	teardown := setup(t)
	defer teardown()

	fileHash := "abcd1234"
	entry, err := findByHash(fileHash)
	if err != nil {
		t.Errorf("findByHash failed with error: %v", err)
	}
	if entry == nil {
		t.Errorf("findByHash failed to find the entry")
	}
}

func TestFindByName(t *testing.T) {

	teardown := setup(t)
	defer teardown()

	fileName := "testfile2.txt"
	entry, err := findByName(fileName)
	if err != nil {
		t.Errorf("fileName failed with error: %v", err)
	}
	if entry == nil {
		t.Errorf("fileName failed to find the entry")
	}
}

func TestFindByHashOrNameWithCorrectHash(t *testing.T) {

	teardown := setup(t)
	defer teardown()

	fileName := "invalid.txt"
	fileHash := "efgh5678"
	entry, err := findByHashOrName(fileHash, fileName)
	if err != nil {
		t.Errorf("fileName failed with error: %v", err)
	}
	if entry == nil {
		t.Errorf("fileHash failed to find the entry")
	}
}

func TestFindByHashOrNameWithCorrectName(t *testing.T) {

	teardown := setup(t)
	defer teardown()

	fileName := "testfile2.txt"
	fileHash := "xxxx"
	entry, err := findByHashOrName(fileHash, fileName)
	if err != nil {
		t.Errorf("fileName failed with error: %v", err)
	}
	if entry == nil {
		t.Errorf("fileName failed to find the entry")
	}
}

func TestFindByHashOrNameWithInvalidDetails(t *testing.T) {

	fileName := "invalid.txt"
	fileHash := "xxxx"
	entry, err := findByHashOrName(fileHash, fileName)
	if err != nil {
		t.Errorf("findByHashOrName failed with error: %v", err)
	}
	if entry != nil {
		t.Errorf("Entry should be empty as per test conditions")
	}
}

func TestCleanCSV(t *testing.T) {
	err := cleanCSV()
	if err != nil {
		t.Errorf("cleanCSV failed with error: %v", err)
	}
}
