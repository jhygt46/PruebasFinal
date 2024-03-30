package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Verificar si se proporciona un argumento
	if len(os.Args) < 2 {
		fmt.Println("Uso: /main <disk|cpu|memory>")
		return
	}

	command := os.Args[1]
	if command == "disk" {
		n, err := DirSize("C:")
		if err == nil {
			fmt.Println(n)
		}
	} else if command == "cpu" {
		fmt.Println(GetMonitoringsCpu())
	} else if command == "memory" {
		fmt.Println(GetMemUsage())
	} else {
		fmt.Println("Error: Commando inexistente")
	}
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}
func GetMonitoringsCpu() float64 {
	idle0, total0 := getCPUSample()
	//fmt.Println(idle0, total0)
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()
	//fmt.Println(idle1, total1)
	IdleTicks := float64(idle1 - idle0)
	TotalTicks := float64(total1 - total0)
	return 100 * (TotalTicks - IdleTicks) / TotalTicks
}
func DirSize(path string) (float64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	sizeMB := float64(size) / 1024.0 / 1024.0
	return sizeMB, err
}
func GetMemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("%v %v %v %v", m.Alloc, m.TotalAlloc, m.Sys, m.NumGC)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func Uint32ToBytes(n uint32) []byte {
	b := make([]byte, 4)
	b[0] = uint8(n / 16777216)
	b[1] = uint8(n / 65536)
	b[2] = uint8(n / 256)
	b[3] = uint8(n % 256)
	return b
}
func Uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	b[0] = uint8(n / 72057594037927936)
	b[1] = uint8(n / 281474976710656)
	b[2] = uint8(n / 1099511627776)
	b[3] = uint8(n / 4294967296)
	b[4] = uint8(n / 16777216)
	b[5] = uint8(n / 65536)
	b[6] = uint8(n / 256)
	b[7] = uint8(n % 256)
	return b
}
