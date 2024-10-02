package file_reader

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"sync"
	"testing"
)

func getChunksFromReader(t *testing.T, fileReader *FileReader) [][]byte {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskChan := make(chan []byte)
	errorChan := make(chan error)

	var actualChunks [][]byte

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := fileReader.ReadChunks(ctx, taskChan); err != nil {
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

	fileReader := NewFileReader(file.Name(), 19)

	expectedChunks := [][]byte{
		[]byte("192.168.1.1\n"),
		[]byte("192.168.1.2\n192.168.1.3\n"),
	}

	actualChunks := getChunksFromReader(t, fileReader)
	if !reflect.DeepEqual(actualChunks, expectedChunks) {
		t.Errorf("Expected chunks: %v, but got: %v", expectedChunks, actualChunks)
	}
}

func createTempFile(t *testing.T, content string) *os.File {
	t.Helper()

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

func BenchmarkBytesLastIndex(b *testing.B) {
	data := []byte("192.168.1.1\n192.168.1.2\n192.168.1.3\n192.168.1.4\n192.168.1.5\n")
	for i := 0; i < b.N; i++ {
		bytes.LastIndex(data, []byte{'\n'})
	}
}

func BenchmarkFindLastNewLine(b *testing.B) {
	data := []byte("192.168.1.1\n192.168.1.2\n192.168.1.3\n192.168.1.4\n192.168.1.5\n")
	for i := 0; i < b.N; i++ {
		findLastNewLine(data)
	}
}
