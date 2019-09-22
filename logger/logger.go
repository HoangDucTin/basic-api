package logger

import (
	"fmt"
)

// Fail prints out the message
// with prefix 'FAIL'.
func Fail(format string, args ...interface{}) {
	message := "[FAIL]: "
	if args != nil {
		message += fmt.Sprintf(format, args)
	} else {
		message += format
	}
	fmt.Println(message)
}

// Help prints out the message
// with prefix 'HELP'.
func Help(format string, args ...interface{}) {
	message := "[HELP]: "
	if args != nil {
		message += fmt.Sprintf(format, args)
	} else {
		message += format
	}
	fmt.Println(message)
}

// Info prints out the message
// with prefix 'INFO'.
func Info(format string, args ...interface{}) {
	message := "[INFO]: "
	if args != nil {
		message += fmt.Sprintf(format, args)
	} else {
		message += format
	}
	fmt.Println(message)
}

// Warn prints out the message
// with prefix 'WARN'.
func Warn(format string, args ...interface{}) {
	message := "[WARN]: "
	if args != nil {
		message += fmt.Sprintf(format, args)
	} else {
		message += format
	}
	fmt.Println(message)
}

// Exit prints out the message
// with prefix 'EXIT'.
func Exit(format string, args ...interface{}) {
	message := "[EXIT]: "
	if args != nil {
		message += fmt.Sprintf(format, args)
	} else {
		message += format
	}
	fmt.Println(message)
}

// Error prints out the message
// with prefix 'ERROR'.
func Error(format string, args ...interface{}) {
	message := "[ERROR]: "
	if args != nil {
		message += fmt.Sprintf(format, args)
	} else {
		message += format
	}
	fmt.Println(message)
}

// End-of-file
