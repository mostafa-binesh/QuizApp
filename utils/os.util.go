package utils

import (
	"os"
)

func FileExistenceCheck(fileName string, dir string) bool {
	if _, err := os.Stat(dir + "/" + fileName); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
