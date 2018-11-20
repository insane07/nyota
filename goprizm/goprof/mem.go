package goprof

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"runtime"
	"runtime/debug"
	"time"
)

// Subset of memory stats provided by runtime.
type MemStats struct {
	// General statistics.
	AllocKB uint64 // KBs allocated and still in use
	SysKB   uint64 // KBs obtained from system (should be sum of XxxSys below)

	// Main allocation heap statistics.
	HeapAllocKB    uint64 // KB allocated and still in use
	HeapSysKB      uint64 // KB obtained from system
	HeapIdleKB     uint64 // KB in idle spans
	HeapInuseKB    uint64 // KB in non-idle span
	HeapReleasedKB uint64 // KB released to the OS
	HeapObjects    uint64 // KB number of allocated objects

	// GC stats
	NextGCKB uint64    // next run in HeapAlloc time (KB)
	LastGC   time.Time // last run in absolute time (msec)
}

// Convert from bytes to KB
func bytesToKB(b uint64) uint64 {
	return b / (1024)
}

// Provide runtime memory stats.
func ReadMemStats() *MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &MemStats{
		AllocKB: bytesToKB(m.Alloc),
		SysKB:   bytesToKB(m.Sys),

		HeapAllocKB:    bytesToKB(m.HeapAlloc),
		HeapSysKB:      bytesToKB(m.HeapSys),
		HeapIdleKB:     bytesToKB(m.HeapIdle),
		HeapInuseKB:    bytesToKB(m.HeapInuse),
		HeapReleasedKB: bytesToKB(m.HeapReleased),
		HeapObjects:    m.HeapObjects,

		NextGCKB: bytesToKB(m.NextGC),
		LastGC:   time.Unix(int64(m.LastGC/(1000000000)), 0),
	}
}

// Force GC and release memory back to OS.
// Warning:
// This should be used only for debugging memory leaks in running services.
// Services can expose a debug url and invoke this method in load test setups.
func FreeMem() {
	debug.FreeOSMemory()
}

var procMemFields = map[string]bool{
	"MemTotal": true,
	"MemFree":  true,
	"Buffers":  true,
	"Cached":   true,
}

// SysMem return total and free memory in bytes.
func SysMem() (total uint64, free uint64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return
	}
	defer f.Close()

	pending := len(procMemFields)
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		fields := strings.Split(scan.Text(), ":")
		if len(fields) != 2 {
			continue
		}

		name := strings.TrimSpace(fields[0])
		memKB := strings.Fields(strings.TrimSpace(fields[1]))[0]
		if _, ok := procMemFields[name]; !ok {
			continue
		}

		var mem uint64
		mem, err = strconv.ParseUint(memKB, 10, 64)
		if err != nil {
			err = fmt.Errorf("invalid mem field:%s val:%s (%s)", name, memKB, err)
			return
		}

		mem = mem * 1024
		switch name {
		case "MemTotal":
			total = mem
		case "MemFree", "Buffers", "Cached":
			free += mem
		}

		pending--
		if pending == 0 {
			return
		}
	}

	if err = scan.Err(); err != nil {
		return
	}

	err = fmt.Errorf("unable to get required mem fields")
	return
}

//ThisProcMem returns total physical memory consumed by this process.
func ProcMemThis() (uint64, error) {
	return ProcMem(os.Getpid())
}

//ProcMem returns total physical memory consumed by process given its pid.
func ProcMem(pid int) (uint64, error) {
	procFile := fmt.Sprintf("/proc/%d/status", pid)
	f, err := os.Open(procFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		if !strings.HasPrefix(line, "VmRSS:") {
			continue
		}

		fields := strings.Fields(line)
		if l := len(fields); l != 3 {
			return 0, fmt.Errorf("incorrect num of field:%d in VmRSS", l)
		}

		memStr := strings.TrimSpace(fields[1])
		memKB, err := strconv.ParseUint(memStr, 10, 64)
		if err != nil {
			return 0, err
		}

		return memKB * 1024, nil
	}

	if err := scan.Err(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("VmRSS missing")
}
