package reader

import (
	"context"
	"os"
	"reflect"
	"sync"
	"testing"
)

func getChunksFromReader(t *testing.T, reader *Reader) [][]byte {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskChan := make(chan []byte)
	errorChan := make(chan error)

	var actualChunks [][]byte

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := reader.ReadChunks(ctx, taskChan); err != nil {
			errorChan <- err
		}

		close(errorChan)
	}()

	for chunk := range taskChan {
		actualChunks = append(actualChunks, chunk)
	}

	if err := <-errorChan; err != nil {
		t.Fatalf("ReadChunks failed: %v", err)
	}

	wg.Wait()

	return actualChunks
}

func TestReadChunks_Normal(t *testing.T) {
	content := "192.168.1.1\n192.168.1.2\n192.168.1.3\n"

	file := createTempFile(t, content)
	defer os.Remove(file.Name())

	reader := NewReader(file.Name(), 19)

	expectedChunks := [][]byte{
		[]byte("192.168.1.1\n"),
		[]byte("192.168.1.2\n192.168.1.3\n"),
	}

	actualChunks := getChunksFromReader(t, reader)
	if !reflect.DeepEqual(actualChunks, expectedChunks) {
		t.Errorf("Expected chunks: %v, but got: %v", expectedChunks, actualChunks)
	}
}

func createTempFile(t *testing.T, content string) *os.File {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return file
}
