package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Serve(port string) {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/v1/store", storeHandler)
	http.HandleFunc("/api/v1/update", updateHandler)
	http.HandleFunc("/api/v1/exists", existenceCheckHandler)
	http.HandleFunc("/api/v1/list", listHandler)
	http.HandleFunc("/api/v1/delete", deleteHandler)
	// Add more handlers for other operations

	log.Println(fmt.Sprintf("Server is starting on port %s...", port))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("resources/readme.md")
	if err != nil {
		http.Error(w, "File reading error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form in the request
	err := r.ParseMultipartForm(50 << 20) // limit your maxMemory here
	if err != nil {
		log.Println("Error parsing the form:", err)
		http.Error(w, "Error parsing the file", http.StatusInternalServerError)
		return
	}

	fileName := r.FormValue("filename")
	err = validateRequiredField("filename", fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the form
	// todo use the fileHeader for accurate file size information `fileHeader.Size`
	file, _, err := r.FormFile("file") // retrieve the file from form data
	if err != nil {
		log.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	if file == nil {
		log.Println("No file was uploaded")
		http.Error(w, "No file was uploaded", http.StatusBadRequest)
		return
	}
	defer CloseMultipartFile(file)

	// Create a new file in the files directory
	dst, err := os.Create(getFilePath(fileName)) // create a new file with the same name
	if err != nil {
		log.Println("Error creating the file:", err)
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer CloseFile(dst)

	// Write the contents of the uploaded file to the new file
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println("Error copying the file:", err)
		http.Error(w, "Error writing to the file", http.StatusInternalServerError)
		return
	}

	// Compute the MD5 hash of the uploaded file
	// todo use filepath.Join function to generate the file-path
	md5Hash, err := ComputeMD5Hash(getFilePath(fileName))
	if err != nil {
		log.Println("Error computing the MD5 hash:", err)
		http.Error(w, "Error computing the MD5 hash", http.StatusInternalServerError)
		return
	}

	//check if the file already exists
	entry, err := findByHash(md5Hash)
	if err != nil {
		log.Println("Error finding file hash:", err)
		// todo read this message from a config file
		http.Error(w, "There was a problem verifying existing hashes. "+
			"Please try the following:\n\n    Refresh the page and try again.\n\nIf the error persists, "+
			"contact an administrator for assistance.", http.StatusInternalServerError)
		return
	}

	if entry != nil {
		log.Println("File already exists")
		http.Error(w, "File already exists", http.StatusConflict)
		return
	}

	// store the file details in the csv file
	// todo
	fileDetails := FileDetails{Filename: fileName, FileSize: r.ContentLength, FileHash: md5Hash}
	err = storeInCSV(fileDetails)
	if err != nil {
		log.Println("Error storing file details:", err)
		http.Error(w, "Error storing file details", http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("File uploaded successfully"))
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form in the request
	err := r.ParseMultipartForm(50 << 20) // limit your maxMemory here
	if err != nil {
		log.Println("Error parsing the form:", err)
		http.Error(w, "Error parsing the file", http.StatusInternalServerError)
		return
	}

	// Get the previous file name from the form
	prevFilename := r.FormValue("prevFilename")
	if prevFilename == "" {
		log.Println("Previous file name is required")
		http.Error(w, "Previous file name is required", http.StatusBadRequest)
		return
	}

	newFileName := r.FormValue("filename")

	// Get the duplicate flag from the form
	duplicate, err := strconv.ParseBool(r.FormValue("duplicate"))
	if err != nil {
		log.Println("Error parsing duplicate value:", err)
		http.Error(w, "Invalid duplicate value", http.StatusBadRequest)
		return
	}

	// Get the new file from the form
	file, _, err := r.FormFile("file") // retrieve the file from form data
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		log.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	if file != nil {
		defer CloseMultipartFile(file)
	}

	record, err := findByName(prevFilename)
	// todo use the helper-function to reduce the code duplication of error handling
	if err != nil {
		log.Println("Error finding file name:", err)
		http.Error(w, "Error finding file name", http.StatusInternalServerError)
		return
	}
	if record == nil {
		log.Println("File does not exist")
		http.Error(w, "Record does not exists",
			http.StatusNotFound)
		return
	}

	if file == nil {
		// If no new file is provided, either update the existing record entry or create
		// a duplicate of the existing file with new record.
		// todo case to handle when duplicate is true and file name is also changed but content is not changed
		err := ManageFileUpdate(duplicate, newFileName, *record)
		if err != nil {
			http.Error(w, "Error updating the file", http.StatusInternalServerError)
			return
		}
	} else {
		// If a new file is provided, validate the new file name
		err := validateRequiredField("filename", newFileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create a new file in the files directory
		dst, err := os.Create(getFilePath(newFileName))
		if err != nil {
			log.Println("Error creating the file:", err)
			http.Error(w, "Error creating the file", http.StatusInternalServerError)
			return
		}
		defer CloseFile(dst)

		// Write the contents of the uploaded file to the new file
		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println("Error copying the file:", err)
			http.Error(w, "Error writing to the file", http.StatusInternalServerError)
			return
		}
		// Compute the MD5 hash of the new file
		md5Hash, err := ComputeMD5Hash(getFilePath(newFileName))
		newRecord := FileDetails{Filename: newFileName, FileSize: r.ContentLength, FileHash: md5Hash}
		// Update the old record with the new record and delete the old file
		err = modifyRecordAndFile(*record, newRecord)
		if err != nil {
			http.Error(w, "Error updating the old record and deleting the old file: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
	}
	_, err = w.Write([]byte("successfully updated the file"))
	if err != nil {
		log.Println("Error writing success response:", err)
	}
}

func existenceCheckHandler(w http.ResponseWriter, r *http.Request) {

	// Parse the form data to get the hash value
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing the form:", err)
		http.Error(w, "Error parsing the form", http.StatusInternalServerError)
		return
	}

	hash := r.FormValue("hash")
	name := r.FormValue("name")

	// Create a closure that checks if either the hash or the name is provided
	if func(hash string, name string) error {
		if hash == "" && name == "" {
			return errors.New("either hash or name is required")
		}
		return nil
	}(hash, name) != nil {
		http.Error(w, "either hash or name is required", http.StatusBadRequest)
		return
	}

	// Use the findByHash function to check if a file with the given hash exists
	record, err := findByHashOrName(hash, name)
	if err != nil {
		log.Println("Error executing findByHashOrName:", err)
		http.Error(w, "Error in finding record by hash or name", http.StatusInternalServerError)
		return
	}

	// If the file does not exist, respond with an appropriate message
	if record == nil {
		http.Error(w, "record does not exist", http.StatusNotFound)
		return // return here to prevent further execution
	}

	// Marshal the record into JSON
	recordJson, err := json.Marshal(record)
	if err != nil {
		log.Println("Error marshalling the record:", err)
		http.Error(w, "Error marshalling the record", http.StatusInternalServerError)
		return
	}
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON record to the response
	_, err = w.Write(recordJson)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {

	entries, err := getAllEntries()
	if err != nil {
		log.Println("Error getting all entries:", err)
		http.Error(w, "Error getting all entries", http.StatusInternalServerError)
		return
	}
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Use json.NewEncoder to write entries as a JSON array to writer
	err = json.NewEncoder(w).Encode(entries)
	if err != nil {
		log.Println("Error encoding entries to JSON:", err)
		http.Error(w, "Error encoding entries to JSON", http.StatusInternalServerError)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	// Parse the form data to get the file name
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing the form:", err)
		http.Error(w, "Error parsing the form", http.StatusInternalServerError)
		return
	}

	filename := r.FormValue("filename")

	// Use the findByName function to check if a file with the given name exists
	record, err := findByName(filename)
	if err != nil {
		log.Println("Error executing findByName:", err)
		http.Error(w, "Error in finding record by name", http.StatusInternalServerError)
		return
	}

	// If the file does not exist, respond with an appropriate message
	if record == nil {
		http.Error(w, "record does not exist", http.StatusNotFound)
		return
	}

	// Delete the record from the CSV file
	err = deleteFromCSV(filename)
	if err != nil {
		log.Println("Error deleting record from CSV:", err)
		http.Error(w, "Error deleting record from CSV", http.StatusInternalServerError)
		return
	}

	// Delete the file from the file system
	err = os.Remove(getFilePath(filename))
	if err != nil {
		log.Println("Error deleting the file:", err)
		http.Error(w, "Error deleting the file", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	_, err = w.Write([]byte("File deleted successfully"))
	if err != nil {
		log.Println("Error writing response:", err)
	}
}
