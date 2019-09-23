package tlog

import (
	// Native packages
	"context"
	"io/ioutil"
	"log"
	"os"

	// Third parties
	"gopkg.in/yaml.v2"
)

// LoggerT is used for getting
// the logger out of the
// context of the request.
const LoggerT = "TLogger"

type (
	// Logger is the API wrapper for underlying logging libraries.
	Logger interface {
		WithFields(fields map[string]interface{}) Logger
		WithPrefix(prefix string) Logger
		WithError(err error) Logger

		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Printf(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Warningf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Panicf(format string, args ...interface{})
		Fatalf(format string, args ...interface{})

		Debug(args ...interface{})
		Info(args ...interface{})
		Print(args ...interface{})
		Warn(args ...interface{})
		Warning(args ...interface{})
		Error(args ...interface{})
		Panic(args ...interface{})
		Fatal(args ...interface{})

		Debugln(args ...interface{})
		Infoln(args ...interface{})
		Println(args ...interface{})
		Warnln(args ...interface{})
		Warningln(args ...interface{})
		Errorln(args ...interface{})
		Panicln(args ...interface{})
		Fatalln(args ...interface{})

		TDebugf(ctx context.Context, format string, args ...interface{})
		TInfof(ctx context.Context, format string, args ...interface{})
		TPrintf(ctx context.Context, format string, args ...interface{})
		TWarnf(ctx context.Context, format string, args ...interface{})
		TWarningf(ctx context.Context, format string, args ...interface{})
		TErrorf(ctx context.Context, format string, args ...interface{})
		TPanicf(ctx context.Context, format string, args ...interface{})
		TFatalf(ctx context.Context, format string, args ...interface{})

		TDebug(ctx context.Context, args ...interface{})
		TInfo(ctx context.Context, args ...interface{})
		TPrint(ctx context.Context, args ...interface{})
		TWarn(ctx context.Context, args ...interface{})
		TWarning(ctx context.Context, args ...interface{})
		TError(ctx context.Context, args ...interface{})
		TPanic(ctx context.Context, args ...interface{})
		TFatal(ctx context.Context, args ...interface{})

		TDebugln(ctx context.Context, args ...interface{})
		TInfoln(ctx context.Context, args ...interface{})
		TPrintln(ctx context.Context, args ...interface{})
		TWarnln(ctx context.Context, args ...interface{})
		TWarningln(ctx context.Context, args ...interface{})
		TErrorln(ctx context.Context, args ...interface{})
		TPanicln(ctx context.Context, args ...interface{})
		TFatalln(ctx context.Context, args ...interface{})
	}

	// Config is the configuration of logger.
	Config struct {
		// Log level.
		// Can be one of: debug, info, warn, error, panic
		Level string `yaml:"log_level" json:"level"`
		// Where log will be written to.
		// Can be one of: stdout, stderr, discard, file://path/to/log/file
		Output string `yaml:"log_output" json:"output"`
		// Log output format.
		// Can be one of: text, json
		Format string `yaml:"log_format" json:"format"`
	}
)

// Standard logger that will be used as the default logger for this package.
var std Logger

func init() {
	cfg := getLoggerConfig()

	var err error
	std, err = New(cfg)
	if err != nil {
		panic(err)
	}
}

// readLogConfig tries to parses TLog configurations from config files and environment variables.
func getLoggerConfig() *Config {
	conf := Config{
		Level:  "debug",
		Format: "json",
		Output: "file",
	}

	if cfgFile := os.Getenv("KLOG_CONFIG_FILE"); cfgFile != "" {
		log.Printf("TLog: Read logger config file: %s", cfgFile)
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			panic(err)
		}
		log.Printf("config file: %s", string(b))
		if err := yaml.Unmarshal(b, &conf); err != nil {
			panic(err)
		}
		log.Printf("config file: %+v", conf)
		return &conf
	}

	if env := os.Getenv("KLOG_LEVEL"); env != "" {
		conf.Level = env
	}
	if env := os.Getenv("KLOG_FORMAT"); env != "" {
		conf.Format = env
	}
	if env := os.Getenv("KLOG_OUTPUT"); env != "" {
		conf.Output = env
	}
	return &conf
}

// New creates new logger by provided configurations.
// Currently creating logrus logger by default.
func New(cfg *Config) (Logger, error) {
	return newLogrusLogger(cfg)
}

// WithFields allows to add additional fields to every log record written by the returning logger.
func WithFields(fields map[string]interface{}) Logger {
	return std.WithFields(fields)
}

// WithPrefix allows to add prefix to every log record written by the returning logger.
func WithPrefix(prefix string) Logger {
	return std.WithPrefix(prefix)
}

// WithError adds the error
// to every log record written
// by the running logger.
func WithError(err error) Logger {
	return std.WithError(err)
}

// Debugf prints out the log
// with level of DEBUG.
// It receives format string
// and a list of arguments.
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// Infof prints out the log
// with level of INFO.
// It receives format string
// and a list of arguments.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Printf prints out the log
// with no level.
// It receives format string
// and a list of arguments.
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}

// Warnf prints out the log
// with level of WARN.
// It receives format string
// and a list of arguments.
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Warningf prints out the log
// with level of WARN.
// It receives format string
// and a list of arguments.
func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

// Errorf prints out the log
// with level of ERROR.
// It receives format string
// and a list of arguments.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Panicf prints out the log
// with level of PANIC and
// makes a panic to application.
// It receives format string
// and a list of arguments.
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

// Fatalf prints out the log
// with level of FATAL and
// makes a panic to application.
// It receives format string
// and a list of arguments.
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// Debug prints out the log
// with level of DEBUG.
// It receives a list of arguments.
func Debug(args ...interface{}) {
	std.Debug(args...)
}

// Info prints out the log
// with level of INFO.
// It receives a list of arguments.
func Info(args ...interface{}) {
	std.Info(args...)
}

// Print prints out the log
// with no level.
// It receives a list of arguments.
func Print(args ...interface{}) {
	std.Print(args...)
}

// Warn prints out the log
// with level of WARN.
// It receives a list of arguments.
func Warn(args ...interface{}) {
	std.Warn(args...)
}

// Warning prints out the log
// with level of WARN.
// It receives a list of arguments.
func Warning(args ...interface{}) {
	std.Warning(args...)
}

// Error prints out the log
// with level of ERROR.
// It receives a list of arguments.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Panic prints out the log
// with level of PANIC and
// makes a panic to the application.
// It receives a list of arguments.
func Panic(args ...interface{}) {
	std.Panic(args...)
}

// Fatal prints out the log
// with level of FATAL and
// makes a panic to the application.
// It receives a list of arguments.
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

// Debugln prints out the log
// with level of DEBUG.
// It receives a list of arguments.
func Debugln(args ...interface{}) {
	std.Debugln(args...)
}

// Infoln prints out the log
// with level of INFO.
// It receives a list of arguments.
func Infoln(args ...interface{}) {
	std.Infoln(args...)
}

// Println prints out the log
// with no level.
// It receives a list of arguments.
func Println(args ...interface{}) {
	std.Println(args...)
}

// Warnln prints out the log
// with level of WARN.
// It receives a list of arguments.
func Warnln(args ...interface{}) {
	std.Warnln(args...)
}

// Warningln prints out the log
// with level of WARN.
// It receives a list of arguments.
func Warningln(args ...interface{}) {
	std.Warningln(args...)
}

// Errorln prints out the log
// with level of ERROR.
// It receives a list of arguments.
func Errorln(args ...interface{}) {
	std.Errorln(args...)
}

// Panicln prints out the log
// with level of PANIC and
// makes a panic to the application.
// It receives a list of arguments.
func Panicln(args ...interface{}) {
	std.Panicln(args...)
}

// Fatalln prints out the log
// with level of FATAL and
// makes a panic to the application.
// It receives a list of arguments.
func Fatalln(args ...interface{}) {
	std.Fatalln(args...)
}

// TDebugf prints out the log
// with level of DEBUG.
// It receives the context for
// getting request ID and logger,
// and format string, and a list
// of arguments.
func TDebugf(ctx context.Context, format string, args ...interface{}) {
	std.TDebugf(ctx, format, args...)
}

// TInfof prints out the log
// with level of DEBUG.
// It receives the context for
// getting request ID and logger,
// and format string, and a list
// of arguments.
func TInfof(ctx context.Context, format string, args ...interface{}) {
	std.TInfof(ctx, format, args...)
}

// TPrintf prints out the log
// with no level.
// It receives the context for
// getting request ID and logger,
// and format string, and a list
// of arguments.
func TPrintf(ctx context.Context, format string, args ...interface{}) {
	std.TPrintf(ctx, format, args...)
}

// TWarnf prints out the log
// with level of WARN.
// It receives the context for
// getting request ID and logger,
// and format string, and a list
// of arguments.
func TWarnf(ctx context.Context, format string, args ...interface{}) {
	std.TWarnf(ctx, format, args...)
}

// TWarningf prints out the log
// with level of WARN.
// It receives the context for
// getting request ID and logger,
// and format string, and a list
// of arguments.
func TWarningf(ctx context.Context, format string, args ...interface{}) {
	std.TWarningf(ctx, format, args...)
}

// TErrorf prints out the log
// with level of WARN.
// It receives the context for
// getting request ID and logger,
// and format string, and a list
// of arguments.
func TErrorf(ctx context.Context, format string, args ...interface{}) {
	std.TErrorf(ctx, format, args...)
}

// TPanicf prints out the log
// with level of PANIC and makes
// a panic to the application.
// It receives the context for
// getting request ID and logger,
// and format string, and a list
// of arguments.
func TPanicf(ctx context.Context, format string, args ...interface{}) {
	std.TPanicf(ctx, format, args...)
}

// TFatalf prints out the log
// with level of FATAL and makes
// a panic to the application.
// It receives the context for
// getting request ID and logger,
// and format string, and a list
// of arguments.
func TFatalf(ctx context.Context, format string, args ...interface{}) {
	std.TFatalf(ctx, format, args...)
}

// TDebug prints out the log
// with level of DEBUG.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TDebug(ctx context.Context, args ...interface{}) {
	std.TDebug(ctx, args...)
}

// TInfo prints out the log
// with level of INFO.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TInfo(ctx context.Context, args ...interface{}) {
	std.TInfo(ctx, args...)
}

// TPrint prints out the log
// with no level.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TPrint(ctx context.Context, args ...interface{}) {
	std.TPrint(ctx, args...)
}

// TWarn prints out the log
// with level of WARN.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TWarn(ctx context.Context, args ...interface{}) {
	std.TWarn(ctx, args...)
}

// TWarning prints out the log
// with level of WARN.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TWarning(ctx context.Context, args ...interface{}) {
	std.TWarning(ctx, args...)
}

// TError prints out the log
// with level of ERROR.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TError(ctx context.Context, args ...interface{}) {
	std.TError(ctx, args...)
}

// TPanic prints out the log
// with level of PANIC and
// make a panic to the application.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TPanic(ctx context.Context, args ...interface{}) {
	std.TPanic(ctx, args...)
}

// TFatal prints out the log
// with level of FATAL and
// make a panic to the application.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TFatal(ctx context.Context, args ...interface{}) {
	std.TFatal(ctx, args...)
}

// TDebugln prints out the log
// with level of DEBUG.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TDebugln(ctx context.Context, args ...interface{}) {
	std.TDebugln(ctx, args...)
}

// TInfoln prints out the log
// with level of INFO.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TInfoln(ctx context.Context, args ...interface{}) {
	std.TInfoln(ctx, args...)
}

// TPrintln prints out the log
// with no level.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TPrintln(ctx context.Context, args ...interface{}) {
	std.TPrintln(ctx, args...)
}

// TWarnln prints out the log
// with level of WARN.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TWarnln(ctx context.Context, args ...interface{}) {
	std.TWarnln(ctx, args...)
}

// TWarningln prints out the log
// with level of WARN.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TWarningln(ctx context.Context, args ...interface{}) {
	std.TWarningln(ctx, args...)
}

// TErrorln prints out the log
// with level of ERROR.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TErrorln(ctx context.Context, args ...interface{}) {
	std.TErrorln(ctx, args...)
}

// TPanicln prints out the log
// with level of PANIC and
// makes a panic to the application.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TPanicln(ctx context.Context, args ...interface{}) {
	std.TPanicln(ctx, args...)
}

// TFatalln prints out the log
// with level of FATAL and
// makes a panic to the application.
// It receives the context for
// getting request ID and logger,
// and a list of arguments.
func TFatalln(ctx context.Context, args ...interface{}) {
	std.TFatalln(ctx, args...)
}
