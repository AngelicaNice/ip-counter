package reader

import (
	"context"
	"fmt"
	"io"
	"os"
)

const maxTailLen = 15

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

	defer close(taskChan)

	leftOver := make([]byte, 0, maxTailLen)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			buf := make([]byte, r.ChunkSize, r.ChunkSize+maxTailLen)
			readBytes, err := file.Read(buf)

			if err != nil && err != io.EOF {
				return fmt.Errorf("failed to read file: %w", err)
			}

			if readBytes == 0 {
				return nil
			}

			buf = buf[:readBytes]

			toSend := make([]byte, readBytes)
			copy(toSend, buf)

			//lastNewLineIndex := bytes.LastIndex(buf, []byte{'\n'})
			lastNewLineIndex := findLastNewLine(buf)

			toSend = append(leftOver, buf[:lastNewLineIndex+1]...)
			leftOver = make([]byte, len(buf[lastNewLineIndex+1:]))
			copy(leftOver, buf[lastNewLineIndex+1:])

			taskChan <- toSend
		}
	}
}

func findLastNewLine(buf []byte) int {
	searchRange := 16
	if len(buf) < searchRange {
		searchRange = len(buf)
	}

	for i := len(buf) - 1; i >= len(buf)-searchRange; i-- {
		if buf[i] == '\n' {
			return i
		}
	}

	return -1
}
