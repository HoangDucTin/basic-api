package logger

import (
	"fmt"
	"log"
)

// This function will print out
// your message with the prefix
// FAIL.
func Fail(format interface{}, args ...interface{}) {
	if f, ok := format.(string); ok {
		log.Println("FAIL", fmt.Sprintf(f, args...))
		return
	}
	if args == nil {
		log.Println("FAIL", format)
		return
	}
	var a string
	for _, _ = range args {
		a += "%v "
	}
	log.Println("FAIL", format, fmt.Sprintf(a, args...))
}

// This function will print out
// your message with the prefix
// HELP.
func Help(format interface{}, args ...interface{}) {
	if f, ok := format.(string); ok {
		log.Println("HELP", fmt.Sprintf(f, args...))
		return
	}
	if args == nil {
		log.Println("HELP", format)
		return
	}
	var a string
	for _, _ = range args {
		a += "%v "
	}
	log.Println("HELP", format, fmt.Sprintf(a, args...))
}

// This function will print out
// your message with the prefix
// INFO.
func Info(format interface{}, args ...interface{}) {
	if f, ok := format.(string); ok {
		log.Println("INFO", fmt.Sprintf(f, args...))
		return
	}
	if args == nil {
		log.Println("INFO", format)
		return
	}
	var a string
	for _, _ = range args {
		a += "%v "
	}
	log.Println("INFO", format, fmt.Sprintf(a, args...))
}

// This function will print out
// your message with the prefix
// WARN.
func Warn(format interface{}, args ...interface{}) {
	if f, ok := format.(string); ok {
		log.Println("WARN", fmt.Sprintf(f, args...))
		return
	}
	if args == nil {
		log.Println("WARN", format)
		return
	}
	var a string
	for _, _ = range args {
		a += "%v "
	}
	log.Println("WARN", format, fmt.Sprintf(a, args...))
}

// This function will print out
// your message with the prefix
// EXIT.
func Exit(format interface{}, args ...interface{}) {
	if f, ok := format.(string); ok {
		log.Println("EXIT", fmt.Sprintf(f, args...))
		return
	}
	if args == nil {
		log.Println("EXIT", format)
		return
	}
	var a string
	for _, _ = range args {
		a += "%v "
	}
	log.Println("EXIT", format, fmt.Sprintf(a, args...))
}

// This function will print out
// your message with the prefix
// ERROR.
func Error(format interface{}, args ...interface{}) {
	if f, ok := format.(string); ok {
		log.Println("ERROR", fmt.Sprintf(f, args...))
		return
	}
	if args == nil {
		log.Println("ERROR", format)
		return
	}
	var a string
	for _, _ = range args {
		a += "%v "
	}
	log.Println("ERROR", format, fmt.Sprintf(a, args...))
}

// End-of-file
