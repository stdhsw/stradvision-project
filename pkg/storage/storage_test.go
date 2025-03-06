package storage

import (
	"fmt"
	"os"
	"testing"
)

func TestStorage(t *testing.T) {
	path, _ := os.Getwd()

	sHandler, err := NewHandler("testfile", path,
		WithMaxFileSize(50),
		WithMaxFileCount(6),
	)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 20; i++ {
		sHandler.WriteData([]byte("0123456789\n"))
	}

	fmt.Println(sHandler.GetSortFileList())
	fmt.Printf("")
}
