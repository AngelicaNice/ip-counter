package processor

import (
	"bytes"
	"context"
	"fmt"

	"bufio"
	"sync"

	"github.com/axiomhq/hyperloglog"
	"github.com/sirupsen/logrus"
)

type Processor struct {
	Logger    *logrus.Logger
	Pool      *sync.Pool
	ChunkSize int64
}

func NewProcessor(logger *logrus.Logger, chunkSize int64) *Processor {
	return &Processor{
		Logger:    logger,
		ChunkSize: chunkSize,
		Pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, chunkSize))
			},
		},
	}
}

type Result struct {
	HLL *hyperloglog.Sketch
	Err error
}

func (p *Processor) Process(ctx context.Context, data []byte, resultChan chan<- Result) {
	buf := p.Pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer p.Pool.Put(buf)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(buf.Bytes(), int(p.ChunkSize))

	hll := hyperloglog.New()

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			p.Logger.Warn("Processor terminated by context.")
			resultChan <- Result{nil, ctx.Err()}

			return
		default:
			line := scanner.Bytes()
			hll.Insert(line)
			//Compact(&line)
			//hll.Insert(line)
		}
	}

	if err := scanner.Err(); err != nil {
		resultChan <- Result{nil, fmt.Errorf("error reading data chunk: %w", err)}
		return
	}

	resultChan <- Result{HLL: hll, Err: nil}
}

//nolint:mnd
func Compact(ipBytes *[]byte) {
	partIndex := 0
	num := 0

	for i := 0; i < len(*ipBytes); i++ {
		char := (*ipBytes)[i]

		switch char {
		case '.':
			(*ipBytes)[partIndex] = byte(num)
			partIndex++
			num = 0
		default:
			num = num*10 + int(char-'0')
		}
	}

	(*ipBytes)[partIndex] = byte(num)
	*ipBytes = (*ipBytes)[:4]
}
