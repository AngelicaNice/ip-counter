package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/AngelicaNice/ip-counter/internal/file_reader"
	"github.com/AngelicaNice/ip-counter/internal/processor"
	"github.com/AngelicaNice/ip-counter/internal/profiler"
	"github.com/AngelicaNice/ip-counter/internal/utils"
)

const baseChunkSize int64 = 64 * 1024 * 1024

func main() {
	logger := utils.InitLogger("execution.log", logrus.InfoLevel)

	_, _ = maxprocs.Set(maxprocs.Logger(logger.Infof))

	prof := profiler.NewProfiler()
	defer prof.Stop()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	start := time.Now()

	fileName := "ip_addresses"

	taskChan := make(chan []byte, runtime.GOMAXPROCS(0))
	resultChan := make(chan processor.Result, runtime.GOMAXPROCS(0))

	fileReader := file_reader.NewFileReader(fileName, baseChunkSize)
	proc := processor.NewProcessor(logger, baseChunkSize)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := fileReader.ReadChunks(ctx, taskChan); err != nil {
			logger.Fatalf("Reader failed: %v", err)
		}
	}()

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for task := range taskChan {
				proc.Process(ctx, task, resultChan)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	finalHLL := hyperloglog.New()

	for result := range resultChan {
		if result.Err != nil {
			logger.Fatalf("Processing error: %v", result.Err)
		}

		if err := finalHLL.Merge(result.HLL); err != nil {
			logger.Fatalf("Error merging HyperLogLog results: %v", err)
		}
	}

	logger.Infof("Total unique IP addresses: %d", finalHLL.Estimate())
	logger.Infof("Total execution time: %v", time.Since(start))
}
