package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, data []byte) *os.File {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_, err = file.Write(data)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	return file
}

func TestProcessChunk_Normal(t *testing.T) {
	testData := []byte("192.168.1.1\n192.168.1.2\n192.168.1.3\n")
	file := createTempFile(t, testData)
	defer os.Remove(file.Name())

	ctx := context.Background()

	result := processChunk(ctx, file.Name(), 0, int64(len(testData)))

	assert.NoError(t, result.err, "Expected no error in normal processing")
	assert.NotNil(t, result.hll, "Expected a valid HLL result")
	assert.Equal(t, uint64(3), result.hll.Estimate(), "Expected 3 unique IP addresses")
}

func TestProcessChunk_SmallBlock(t *testing.T) {
	testData := []byte("192.168.1.1\n192.168.1.2\n")
	file := createTempFile(t, testData)
	defer os.Remove(file.Name())

	ctx := context.Background()

	result := processChunk(ctx, file.Name(), 0, int64(len(testData)))

	assert.NoError(t, result.err, "Expected no error in small block processing")
	assert.NotNil(t, result.hll, "Expected a valid HLL result")
	assert.Equal(t, uint64(2), result.hll.Estimate(), "Expected 2 unique IP addresses")
}

func TestProcessChunk_BoundaryLines(t *testing.T) {
	testData := []byte("192.168.1.1\n192.168.1.2\n192.168.1.3\n192.168.1.4\n192.168.1.5\n")
	file := createTempFile(t, testData)
	defer os.Remove(file.Name())

	ctx := context.Background()

	blockSize := int64(len("192.168.1.1\n192.168.1.2\n192.168.1.3\n192."))

	result := processChunk(ctx, file.Name(), 0, blockSize)

	assert.NoError(t, result.err, "Expected no error in boundary line processing")
	assert.NotNil(t, result.hll, "Expected a valid HLL result")
	assert.LessOrEqual(t, uint64(3), result.hll.Estimate(), "Expected 3 or more unique IP addresses")
}

func TestProcessChunk_WithContextCancellation(t *testing.T) {
	testData := []byte("192.168.1.1\n192.168.1.2\n192.168.1.3\n192.168.1.4\n")
	file := createTempFile(t, testData)
	defer os.Remove(file.Name())

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	result := processChunk(ctx, file.Name(), 0, int64(len(testData)))

	if result.hll.Estimate() != 0 {
		t.Errorf("Expected hll to be empty due to context cancellation, but got a full object")
	}
}
