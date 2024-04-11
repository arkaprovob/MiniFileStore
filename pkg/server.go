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

	// Get the file from the form
	file, _, err := r.FormFile("file") // retrieve the file from form data
	if err != nil {
		log.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer CloseMultipartFile(file)

	// Create a new file in the files directory
	// todo get teh filepath from the environment variable `os.Getenv("FILES_DIR")`
	dst, err := os.Create("files/" + r.FormValue("filename")) // create a new file with the same name
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

	// Write a success message to the response
	_, err = w.Write([]byte("File uploaded successfully"))
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	// Implement file updating logic here
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Implement file deleting logic here
}
