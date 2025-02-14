package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var TestFileName = "testfile.txt"
var TestFileLocation = filepath.Join("test-resources", TestFileName)
var Test2FileLocation = filepath.Join("test-resources", "test2file.txt")

func removeAllFilesInDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer func(d *os.File) {
		err := d.Close()
		if err != nil {
			log.Fatalf("Failed to close directory: %s", err)
		}
	}(d)

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		if name == ".ignore" {
			continue
		}
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return nil
}

func teardown() {
	err := removeAllFilesInDirectory("test-resources/file-store")
	if err != nil {
		log.Fatalf("Failed to remove all files in directory: %s", err)
	}

	err = removeAllFilesInDirectory("test-resources/record-store")
	if err != nil {
		log.Fatalf("Failed to remove all files in directory: %s", err)
	}
}

func fileStoreSetup(t *testing.T) func() {
	TestCleanCSV(t)
	TestStoreHandler(t)
	// return the teardown function
	return func() {
		teardown()
	}
}

func csvSetup(t *testing.T) func() {
	TestCleanCSV(t)
	TestStoreMultipleEntriesInCSV(t)
	// return the teardown function
	return func() {
		teardown()
	}
}

// todo need to check if in good idea to refer to the a test function from another test file

func TestStoreHandler(t *testing.T) {
	TestCleanCSV(t)

	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Open the file we want to add to the form
	file, err := os.Open(TestFileLocation)
	if err != nil {
		t.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(file)

	// Add the file to the form
	part, err := writer.CreateFormFile("file", TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return
	}

	err = writer.WriteField("filename", TestFileName)
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	// Create a new HTTP request with the form data
	req, err := http.NewRequest("POST", "/api/v1/store", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the storeHandler function
	storeHandler(rr, req)

	// Check the HTTP response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the HTTP response body
	expected := "File uploaded successfully"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateHandlerCaseNewFile(t *testing.T) {

	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the previous file name to the form
	err := writer.WriteField("prevFilename", "oldfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Add the new file name to the form
	err = writer.WriteField("filename", "newfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Add the duplicate flag to the form
	err = writer.WriteField("duplicate", "false")
	if err != nil {
		t.Fatal(err)
	}

	// Open the file we want to add to the form
	file, err := os.Open(TestFileLocation)
	if err != nil {
		t.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing the file:", err)
		}
	}(file)

	// Add the file to the form
	part, err := writer.CreateFormFile("file", TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP request with the form data
	req, err := http.NewRequest("POST", "/api/v1/update", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the updateHandler function
	updateHandler(rr, req)

	// Check the HTTP response code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	// Check the HTTP response body
	expected := "Record does not exists"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
func TestUpdateHandlerCaseUpdateName(t *testing.T) {
	teardown := fileStoreSetup(t)
	defer teardown()

	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the previous file name to the form
	err := writer.WriteField("prevFilename", TestFileName)
	if err != nil {
		t.Fatal(err)
	}

	// Add the new file name to the form
	err = writer.WriteField("filename", "newTestFile.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Add the duplicate flag to the form
	err = writer.WriteField("duplicate", "false")
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP request with the form data
	req, err := http.NewRequest("POST", "/api/v1/update", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the updateHandler function
	updateHandler(rr, req)

	// Check the HTTP response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the HTTP response body
	expected := "successfully updated the file"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestUpdateHandlerCaseDuplicate(t *testing.T) {
	teardown := fileStoreSetup(t)
	defer teardown()

	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the previous file name to the form
	err := writer.WriteField("prevFilename", TestFileName)
	if err != nil {
		t.Fatal(err)
	}

	// Add the new file name to the form
	err = writer.WriteField("filename", "newTestFile.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Add the duplicate flag to the form
	err = writer.WriteField("duplicate", "true")
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP request with the form data
	req, err := http.NewRequest("POST", "/api/v1/update", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the updateHandler function
	updateHandler(rr, req)

	// Check the HTTP response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the HTTP response body
	expected := "successfully updated the file"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestUpdateHandlerCaseUpdateFileContent(t *testing.T) {
	teardown := fileStoreSetup(t)
	defer teardown()

	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the previous file name to the form
	err := writer.WriteField("prevFilename", TestFileName)
	if err != nil {
		t.Fatal(err)
	}

	// Add the new file name to the form
	err = writer.WriteField("filename", TestFileName)
	if err != nil {
		t.Fatal(err)
	}

	// Add the duplicate flag to the form
	err = writer.WriteField("duplicate", "false")
	if err != nil {
		t.Fatal(err)
	}

	// Open the file we want to add to the form
	file, err := os.Open(Test2FileLocation)
	if err != nil {
		t.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing the file:", err)
		}
	}(file)

	// Add the file to the form
	part, err := writer.CreateFormFile("file", TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP request with the form data
	req, err := http.NewRequest("POST", "/api/v1/update", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the updateHandler function
	updateHandler(rr, req)

	// Check the HTTP response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the HTTP response body
	expected := "successfully updated the file"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestExistenceCheckHandler(t *testing.T) {
	teardown := csvSetup(t)
	defer teardown()

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/api/v1/exists", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the hash and name to the request's form data
	q := req.URL.Query()
	q.Add("hash", "abcd1234")
	q.Add("name", "testfile1.txt")
	req.URL.RawQuery = q.Encode()

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call existenceCheckHandler
	existenceCheckHandler(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect
	// todo use FileDetails object instead of string comparison
	expected := `{"Filename":"testfile1.txt","FileSize":1234,"FileHash":"abcd1234","WordCount":10}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestNonExistentRecordInExistenceCheckHandler(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/api/v1/exists", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the hash and name to the request's form data
	q := req.URL.Query()
	q.Add("hash", "invalid")
	q.Add("name", "invalid")
	req.URL.RawQuery = q.Encode()

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call existenceCheckHandler
	existenceCheckHandler(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// Check the response body is what we expect
	expected := `record does not exist`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestListHandler(t *testing.T) {
	teardown := csvSetup(t)
	defer teardown()

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/api/v1/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call listHandler
	listHandler(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Unmarshal the response body into a slice of FileDetails
	var entries []FileDetails
	err = json.Unmarshal(rr.Body.Bytes(), &entries)
	if err != nil {
		t.Fatal(err)
	}

	// Check the length of the entries slice is greater than 0
	if len(entries) <= 0 {
		t.Errorf("handler returned no entries, expected at least one entry")
	}
}

func TestDeleteHandler(t *testing.T) {
	teardown := fileStoreSetup(t)
	defer teardown()

	// Create form data
	data := url.Values{}
	data.Set("filename", TestFileName)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/api/v1/delete", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call deleteHandler
	deleteHandler(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := `File deleted successfully`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestWordFrequencyHandler(t *testing.T) {
	teardown := fileStoreSetup(t)
	defer teardown()

	// Create a new HTTP request with the form data
	form := url.Values{}
	form.Add("noOfWords", "10")
	form.Add("mostFrequent", "true")
	req, err := http.NewRequest("POST", "/api/v1/frequency", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the wordFrequencyHandler function
	wordFrequencyHandler(rr, req)

	// Check the HTTP response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response body into a slice of WordCount
	var result Frequencies
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	// Check the length of the result slice is 10
	if len(result) != 10 {
		t.Errorf("Expected 10 results, got %d", len(result))
	}

	// Check the count for each word is greater than 0
	for _, frequencyCount := range result {
		if frequencyCount.Count <= 0 {
			t.Errorf("Expected count to be greater than 0 for word %s, got %d",
				frequencyCount.Word, frequencyCount.Count)
		}
	}
}
