package main

import "fmt"

func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Print("Error:", msg)
}

func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Print("Warning:", msg)
}
