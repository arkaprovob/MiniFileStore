package pkg

import (
	"testing"
)

func TestComputeMD5Hash(t *testing.T) {
	// Compute the MD5 hash of the file
	hash, err := ComputeMD5Hash(TestFileLocation)
	if err != nil {
		t.Fatalf("ComputeMD5Hash failed: %v", err)
	}

	// The expected MD5 hash of  the test file is "2d79685c999a6f7f77756e9948bd975e"
	expectedHash := "2d79685c999a6f7f77756e9948bd975e"
	if hash != expectedHash {
		t.Errorf("Unexpected hash. Got %s, expected %s", hash, expectedHash)
	}
}

//func TestHandleFileDuplicationOrUpdate(t *testing.T) {
//	// Prepare the previous file details
//	previousFileDetails := FileDetails{
//		Filename: "testfile.txt",
//		FileSize: 12008,
//		FileHash: "2d79685c999a6f7f77756e9948bd975e",
//	}
//
//	err := HandleFileDuplicationOrUpdate(true, "newFile", previousFileDetails)
//	if err != nil {
//		t.Fatalf("HandleFileDuplicationOrUpdate failed when duplicate is true: %v", err)
//	}
//
//	err = HandleFileDuplicationOrUpdate(false, "updatedFile", previousFileDetails)
//	if err != nil {
//		t.Fatalf("HandleFileDuplicationOrUpdate failed when duplicate is false: %v", err)
//	}
//
//}
