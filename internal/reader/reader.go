package reader

import (
	"context"
	"fmt"
	"io"
	"os"
)

type Reader struct {
	FileName  string
	ChunkSize int64
}

func NewReader(fileName string, chunkSize int64) *Reader {
	return &Reader{
		FileName:  fileName,
		ChunkSize: chunkSize,
	}
}

func (r *Reader) ReadChunks(ctx context.Context, taskChan chan<- []byte) error {
	file, err := os.Open(r.FileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	buf := make([]byte, r.ChunkSize)
	defer close(taskChan)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
				return fmt.Errorf("failed to read file: %w", err)
			}
			if n == 0 {
				return nil
			}

			taskChan <- buf[:n]
		}
	}
}
