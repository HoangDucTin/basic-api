package logger

import (
	"fmt"
)

// Fail prints out the message
// with prefix 'FAIL'.
func Fail(format string, args ...interface{}) {
	fmt.Println("FAIL", fmt.Sprintf(format, args...))
}

// Help prints out the message
// with prefix 'HELP'.
func Help(format string, args ...interface{}) {
	fmt.Println("HELP", fmt.Sprintf(format, args...))
}

// Info prints out the message
// with prefix 'INFO'.
func Info(format string, args ...interface{}) {
	fmt.Println("INFO", fmt.Sprintf(format, args))
}

// Warn prints out the message
// with prefix 'WARN'.
func Warn(format string, args ...interface{}) {
	fmt.Println("WARN", fmt.Sprintf(format, args))
}

// Exit prints out the message
// with prefix 'EXIT'.
func Exit(format string, args ...interface{}) {
	fmt.Println("EXIT", fmt.Sprintf(format, args))
}

// Error prints out the message
// with prefix 'ERROR'.
func Error(format string, args ...interface{}) {
	fmt.Println("ERROR", fmt.Sprintf(format, args))
}

// End-of-file
