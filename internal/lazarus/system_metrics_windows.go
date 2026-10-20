//go:build windows

package lazarus

import (
	"errors"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32Metrics      = syscall.NewLazyDLL("kernel32.dll")
	getSystemTimes       = kernel32Metrics.NewProc("GetSystemTimes")
	globalMemoryStatusEx = kernel32Metrics.NewProc("GlobalMemoryStatusEx")

	ErrSystemMetricsUnavailable = errors.New("Windows system metrics unavailable")
	ErrCPUIntervalTooShort      = errors.New("CPU sampling interval too short")
)

type windowsFiletime struct {
	LowDateTime  uint32
	HighDateTime uint32
}

func (f windowsFiletime) value() uint64 {
	return uint64(f.HighDateTime)<<32 | uint64(f.LowDateTime)
}

type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

type SystemMetricsCollector struct {
	mu            sync.Mutex
	previousTotal uint64
	previousIdle  uint64
	primed        bool
}

func NewSystemMetricsCollector() *SystemMetricsCollector {
	return &SystemMetricsCollector{}
}

func readWindowsCPUTimes() (total uint64, idle uint64, err error) {
	var idleTime, kernelTime, userTime windowsFiletime

	result, _, _ := getSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)
	runtime.KeepAlive(&idleTime)
	runtime.KeepAlive(&kernelTime)
	runtime.KeepAlive(&userTime)

	if result == 0 {
		return 0, 0, ErrSystemMetricsUnavailable
	}

	idle = idleTime.value()
	total = kernelTime.value() + userTime.value()
	return total, idle, nil
}

func readWindowsMemoryUsage() (float64, error) {
	status := memoryStatusEx{
		Length: uint32(unsafe.Sizeof(memoryStatusEx{})),
	}

	result, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&status)))
	runtime.KeepAlive(&status)

	if result == 0 || status.TotalPhys == 0 || status.AvailPhys > status.TotalPhys {
		return 0, ErrSystemMetricsUnavailable
	}

	used := status.TotalPhys - status.AvailPhys
	return float64(used) / float64(status.TotalPhys) * 100, nil
}

// Collect membaca CPU dan RAM sistem. Pengambilan CPU pertama menjadi baseline,
// sehingga collector menunggu interval pendek untuk menghitung penggunaan CPU.
func (s *SystemMetricsCollector) Collect() (float64, float64, error) {
	if s == nil {
		return 0, 0, ErrSystemMetricsUnavailable
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	total, idle, err := readWindowsCPUTimes()
	if err != nil {
		return 0, 0, err
	}

	if !s.primed {
		s.previousTotal = total
		s.previousIdle = idle
		s.primed = true
		time.Sleep(250 * time.Millisecond)

		total, idle, err = readWindowsCPUTimes()
		if err != nil {
			return 0, 0, err
		}
	}

	var cpuUsage float64
	deltaTotal := total - s.previousTotal
	deltaIdle := idle - s.previousIdle

	if deltaTotal == 0 || deltaIdle > deltaTotal {
		s.previousTotal = total
		s.previousIdle = idle
		return 0, 0, ErrCPUIntervalTooShort
	}

	cpuUsage = float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100
	s.previousTotal = total
	s.previousIdle = idle

	memoryUsage, err := readWindowsMemoryUsage()
	if err != nil {
		return 0, 0, err
	}

	return cpuUsage, memoryUsage, nil
}
