package pkg

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Frequency struct {
	Word  string
	Count int
}

type Frequencies []Frequency

func (w Frequencies) Len() int           { return len(w) }
func (w Frequencies) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w Frequencies) Less(i, j int) bool { return w[i].Count > w[j].Count } // Sort in descending order

func CountWordsFrequencyParallel(directory string, no int, mostFrequent bool) Frequencies {

	wordCounts := make(map[string]int)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Walk the directory and process files concurrently
	err := fs.WalkDir(os.DirFS(directory), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			// Join the directory path with the file path
			filePath := filepath.Join(directory, path)

			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer CloseFile(file)

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanWords)

			var lineWg sync.WaitGroup

			for scanner.Scan() {
				lineWg.Add(1)
				go func(word string) {
					defer lineWg.Done()
					word = strings.ToLower(word)
					mutex.Lock()
					wordCounts[word]++
					mutex.Unlock()

				}(scanner.Text())
			}
			lineWg.Wait() // Wait for all line-processing goroutines
		}()

		return nil
	})

	if err != nil {
		fmt.Println("Error walking directory:", err)
	}

	wg.Wait() // Wait for all file-processing goroutines

	// Convert the map into a slice for sorting
	result := make(Frequencies, 0, len(wordCounts))
	for word, count := range wordCounts {
		result = append(result, Frequency{word, count})
	}

	sort.Sort(result)
	if len(result) < 10 {
		return result
	}
	if mostFrequent {
		return result[:no] // Return the top 10
	}
	return result[len(result)-no:] // Return the top 10
}
