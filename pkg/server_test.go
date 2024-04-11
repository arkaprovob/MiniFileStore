package pkg

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var TestFileName = "testfile.txt"
var TestFileLocation = filepath.Join("test-resources", TestFileName)

func TestStoreHandler(t *testing.T) {
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
	req, err := http.NewRequest("POST", "/store", body)
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
