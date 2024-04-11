package pkg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Serve(port string) {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/store", storeHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)
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
	duplicate := r.FormValue("duplicate") == "false"

	// Get the new file from the form
	file, _, err := r.FormFile("file") // retrieve the file from form data
	if err != nil && duplicate {
		log.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer CloseMultipartFile(file)

	record, err := findByName(prevFilename)
	// todo use the helper-function to reduce the code duplication of error handling
	if err != nil {
		log.Println("Error finding file name:", err)
		http.Error(w, "Error finding file name", http.StatusInternalServerError)
		return
	}
	if record == nil {
		log.Println("File does not exist")
		http.Error(w, "The requested file could not be found. Please verify the file name and try again.",
			http.StatusNotFound)
		return
	}

	if file == nil {
		// If no new file is provided, either update the existing record entry or create
		// a duplicate of the existing file with new record.
		err := ManageFileUpdate(duplicate, prevFilename, *record)
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Implement file deleting logic here
}
