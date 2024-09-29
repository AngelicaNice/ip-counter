package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"

	//"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"

	//"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/sirupsen/logrus"
)

const (
	MB            uint64 = 1024 * 1024
	GB            uint64 = 1024 * MB
	baseChunkSize int64  = 2 * 1024 * 1024 * 1024
)

var (
	bufferPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	logger = logrus.New()
)

type Result struct {
	hll *hyperloglog.Sketch
	err error
}

func initLogger() {
	file, err := os.OpenFile("execution.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	logger.Out = file
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	initLogger()

	/*
		// Uncomment this section to enable CPU profiling
		cpuFile, err := os.Create("cpu_profile.prof")
		if err != nil {
			logger.Fatalf("Failed to create CPU profile file: %v", err)
		}
		defer cpuFile.Close()

		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			logger.Fatalf("Failed to start CPU profiling: %v", err)
		}
		defer pprof.StopCPUProfile()

		// Uncomment this section to start HTTP server for live profiling
		go func() {
			logger.Info("Starting profiling server at http://localhost:6060/debug/pprof/")
			logger.Info(http.ListenAndServe("localhost:6060", nil))
		}()
	*/

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		logger.Warn("Received termination signal. Waiting for tasks to complete...")
		cancel()
	}()

	fileName := "ip_addresses"

	maxWorkers := runtime.GOMAXPROCS(0)
	logger.Infof("Maximum possible workers: %d", maxWorkers)

	totalStart := time.Now()
	printMemoryUsage("Initial memory usage")

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		logger.Fatalf("Failed to get file info: %v", err)
	}

	fileSize := fileInfo.Size()
	logger.Infof("File size: %.2f GB", float64(fileSize)/float64(GB))

	workerCount := 1
	if fileSize > baseChunkSize {
		workerCount = int(fileSize / baseChunkSize)
		if workerCount > maxWorkers {
			workerCount = maxWorkers
		}
		logger.Infof("Adjusted number of workers: %d", workerCount)
	}

	chunkSize := fileSize / int64(workerCount)
	logger.Infof("Using chunk size: %.2f GB", float64(chunkSize)/float64(GB))

	taskChan := make(chan int64, workerCount)
	resultChan := make(chan Result)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		go worker(ctx, fileName, chunkSize, taskChan, resultChan, &wg)
	}

	for start := int64(0); start < fileSize; start += chunkSize {
		wg.Add(1)
		taskChan <- start
	}

	close(taskChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	finalHLL := hyperloglog.New()
	for result := range resultChan {
		if result.err != nil {
			logger.Fatalf("Error processing chunk: %v", result.err)
		}
		if err := finalHLL.Merge(result.hll); err != nil {
			logger.Fatalf("Error merging HyperLogLog results: %v", err)
		}
	}

	logger.Infof("Total unique IP addresses: %d", finalHLL.Estimate())
	logger.Infof("Total execution time: %v", time.Since(totalStart))

	printMemoryUsage("Memory usage after completion")

	/*
		// Uncomment this section to save heap profile
		memFile, err := os.Create("heap_profile.prof")
		if err != nil {
			logger.Fatalf("Failed to create heap profile file: %v", err)
		}
		defer memFile.Close()

		if err := pprof.WriteHeapProfile(memFile); err != nil {
			logger.Fatalf("Failed to write heap profile: %v", err)
		}
		logger.Info("Profiles saved to cpu_profile.prof and heap_profile.prof.")
	*/
}

func worker(ctx context.Context, fileName string, chunkSize int64, taskChan <-chan int64, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for startPos := range taskChan {
		select {
		case <-ctx.Done():
			logger.Warnf("Worker terminated: skipping block at position %d", startPos)
			return
		default:
			results <- processChunk(ctx, fileName, startPos, chunkSize)
		}
	}
}

func processChunk(ctx context.Context, fileName string, startPos int64, size int64) Result {
	chunkStart := time.Now()
	file, err := os.Open(fileName)
	if err != nil {
		return Result{nil, fmt.Errorf("failed to open file: %v", err)}
	}
	defer file.Close()

	_, err = file.Seek(startPos, 0)
	if err != nil {
		return Result{nil, fmt.Errorf("failed to set file read position: %v", err)}
	}

	scanner := bufio.NewScanner(file)

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	if startPos != 0 {
		scanner.Scan() // Skip the first partial line
	}

	hll := hyperloglog.New()
	bytesRead := int64(0)

readLoop:
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			logger.Warnf("ProcessChunk terminated: stopping at position %d", startPos+bytesRead)
			break readLoop
		default:
			line := scanner.Text()
			hll.Insert([]byte(line))

			bytesRead += int64(len(line)) + 1
			if bytesRead > size { // Read the boundary line to avoid losing unique IP addresses
				break readLoop
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Result{nil, fmt.Errorf("error reading chunk: %v", err)}
	}

	logger.Infof("Processed chunk (starting at %d) in: %v", startPos, time.Since(chunkStart))
	return Result{hll, nil}
}

func printMemoryUsage(when string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logger.Infof("%s:", when)
	logger.Infof("Memory usage: %v MB", bToMb(m.Alloc))
	logger.Infof("Total allocations: %v MB", bToMb(m.TotalAlloc))
	logger.Infof("System memory: %v MB", bToMb(m.Sys))
	logger.Infof("Garbage collections: %v", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / MB
}
