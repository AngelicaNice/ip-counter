package main

import (
	"bytes"
	"testing"
)

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
