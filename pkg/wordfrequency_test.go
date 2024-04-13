package pkg

import (
	"fmt"
	"testing"
)

func TestCountWordsParallel(t *testing.T) {
	directory := "test-resources"
	no := 10

	result := CountWordsFrequencyParallel(directory, no, true)

	if len(result) != no {
		t.Errorf("Expected %d results, got %d", no, len(result))
	}

	for _, wc := range result {
		fmt.Println(wc.Word, ":", wc.Count)
	}
}
