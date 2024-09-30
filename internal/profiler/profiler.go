package profiler

import (
	"ipaddresses/internal/utils"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/sirupsen/logrus"
)

type Profiler struct {
	cpuFile  *os.File
	heapFile *os.File
	logger   *logrus.Logger
}

func NewProfiler() *Profiler {
	profiler := &Profiler{
		logger: utils.InitLogger("profiler.log", logrus.InfoLevel),
	}

	if strings.ToLower(os.Getenv("ENABLE_CPU_PROFILING")) == "true" {
		profiler.startCPUProfiling()
	}

	if strings.ToLower(os.Getenv("ENABLE_HEAP_PROFILING")) == "true" {
		profiler.startHeapProfiling()
	}

	return profiler
}

func (p *Profiler) startCPUProfiling() {
	cpuFileName := os.Getenv("CPU_PROFILE_FILE")
	if cpuFileName == "" {
		cpuFileName = "cpu_profile.prof"
	}

	cpuFile, err := os.Create(cpuFileName)
	if err != nil {
		p.logger.Fatalf("Could not create CPU profile file: %v", err)
	}
	p.cpuFile = cpuFile

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		p.logger.Fatalf("Could not start CPU profiling: %v", err)
	}
	p.logger.Infof("CPU profiling started, results will be saved to: %s", cpuFileName)
}

func (p *Profiler) startHeapProfiling() {
	heapFileName := os.Getenv("HEAP_PROFILE_FILE")
	if heapFileName == "" {
		heapFileName = "heap_profile.prof"
	}

	heapFile, err := os.Create(heapFileName)
	if err != nil {
		p.logger.Fatalf("Could not create heap profile file: %v", err)
	}
	p.heapFile = heapFile

	p.logger.Infof("Heap profiling enabled, results will be saved to: %s", heapFileName)
}

func (p *Profiler) Stop() {
	if p.cpuFile != nil {
		pprof.StopCPUProfile()
		p.cpuFile.Close()
		p.logger.Info("CPU profiling stopped.")
	}

	if p.heapFile != nil {
		if err := pprof.WriteHeapProfile(p.heapFile); err != nil {
			p.logger.Fatalf("Could not write heap profile: %v", err)
		}
		p.heapFile.Close()
		p.logger.Info("Heap profile written successfully.")
	}
}
