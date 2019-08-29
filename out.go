package main

import "fmt"

func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format + "\n", args...)
	fmt.Print("Error:", msg)
}

func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format + "\n", args...)
	fmt.Print("Warning:", msg)
}
