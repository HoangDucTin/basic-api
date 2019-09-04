package logger

import (
	"fmt"
	"log"
)

func Fail(format string, args ...interface{}) {
	log.Println(
		"FAIL",
		fmt.Sprintf(format, args...))
}

func Help(format string, args ...interface{}) {
	log.Println(
		"HELP",
		fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	log.Println(
		"INFO",
		fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	log.Println(
		"WARN",
		fmt.Sprintf(format, args...))
}

func Exit(format string, args ...interface{}) {
	log.Fatalln(
		"EXIT",
		fmt.Sprintf(format, args...))
}

// End-of-file
