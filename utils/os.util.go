package utils

import (
	"fmt"
	"os"
	"runtime"
)

func FileExistenceCheck(fileName string, dir string) bool {
	if _, err := os.Stat(dir + "/" + fileName); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc / 1024 / 1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc / 1024 / 1024)
	fmt.Printf("\tSys = %v MiB", m.Sys / 1024 / 1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}