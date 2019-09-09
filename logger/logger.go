package logger

import (
	"fmt"
	"log"
)

// This function will print out
// your message with the prefix
// FAIL.
func Fail(format string, args ...interface{}) {
	log.Println(
		"FAIL",
		fmt.Sprintf(format, args...))
}

// This function will print out
// your message with the prefix
// HELP.
func Help(format string, args ...interface{}) {
	log.Println(
		"HELP",
		fmt.Sprintf(format, args...))
}

// This function will print out
// your message with the prefix
// INFO.
func Info(format string, args ...interface{}) {
	log.Println(
		"INFO",
		fmt.Sprintf(format, args...))
}

// This function will print out
// your message with the prefix
// WARN.
func Warn(format string, args ...interface{}) {
	log.Println(
		"WARN",
		fmt.Sprintf(format, args...))
}

// This function will print out
// your message with the prefix
// EXIT.
func Exit(format string, args ...interface{}) {
	log.Fatalln(
		"EXIT",
		fmt.Sprintf(format, args...))
}

// This function will print out
// your message with the prefix
// ERROR.
func Error(format string, args ...interface{}) {
	log.Println(
		"ERROR",
		fmt.Sprintf(format, args...))
}

// End-of-file
