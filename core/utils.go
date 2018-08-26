package core

import (
	"fmt"
	"os"
	"runtime"
)

// PathExists reports whether the named file or directory exists.
func PathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Check(err error, message string) {
	if err != nil {
		fmt.Println("Error:")
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", err)
		fmt.Println(message)
		os.Exit(1)
	}
}
