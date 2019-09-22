package tlog

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// logrusLogger is the wrapper of logrus logger to satisfy the TLog.Logger interface.
type logrusLogger struct {
	prefix string // Holds the string prefix of logger
	*logrus.Entry
}

// newLogrusLogger creates a new logrus logger based on provided configurations.
func newLogrusLogger(cfg *Config) (Logger, error) {
	// Log level
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(lvl)

	// Log format
	switch strings.ToLower(cfg.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano,
		})
	}

	// Log output
	output := strings.ToLower(cfg.Output)
	if output == "stdout" {
		logger.SetOutput(os.Stdout)
	} else if output == "discard" {
		logger.SetOutput(ioutil.Discard)
	} else if output == "file" {
		path := "logs/go-facebook-payment"
		if _, err := os.Stat(path); err == nil {
			f, err := os.OpenFile(path+"/"+time.Now().Format("20060102")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				return nil, err
			}
			logger.SetOutput(io.MultiWriter(os.Stdout, f))
			// Start Go Cron for new date in logfile
			c := cron.New()
			_, _ = c.AddFunc("0 0 0 * * *", func() {
				thisdate := time.Now().Format("20060102")
				f, err := os.OpenFile(path+"/"+thisdate+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
				if err != nil {
					return
				}
				logger.SetOutput(io.MultiWriter(os.Stdout, f))
			})
			c.Start()
			trackOpeningLogWriter(output, f)
		} else {
			logger.SetOutput(os.Stdout)
		}
	} else {
		logger.SetOutput(os.Stderr)
	}
	return &logrusLogger{Entry: logrus.NewEntry(logger)}, nil
}

// WithFields allows to add additional fields to every log record written by returning logger.
func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &logrusLogger{
		prefix: l.prefix,
		Entry:  l.Entry.WithFields(fields),
	}
}

// WithPrefix allows to add prefix to every log record written by the returning logger.
func (l *logrusLogger) WithPrefix(prefix string) Logger {
	if l.prefix != "" {
		prefix = l.prefix + "/" + prefix
	}
	return &logrusLogger{
		prefix: prefix,
		Entry:  l.Entry,
	}
}

func (l *logrusLogger) WithError(err error) Logger {
	return &logrusLogger{
		prefix: l.prefix,
		Entry:  l.Entry.WithError(err),
	}
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Debugf(format, args...)
}
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Infof(format, args...)
}
func (l *logrusLogger) Printf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Printf(format, args...)
}
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Warnf(format, args...)
}
func (l *logrusLogger) Warningf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Warningf(format, args...)
}
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Errorf(format, args...)
}
func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Panicf(format, args...)
}
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Fatalf(format, args...)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Debug(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Debug(args...)
}
func (l *logrusLogger) Info(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Info(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Info(args...)
}
func (l *logrusLogger) Print(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Print(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Print(args...)
}
func (l *logrusLogger) Warn(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Warn(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Warn(args...)
}
func (l *logrusLogger) Warning(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Warning(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Warning(args...)
}
func (l *logrusLogger) Error(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Error(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Error(args...)
}
func (l *logrusLogger) Panic(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Panic(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Panic(args...)
}
func (l *logrusLogger) Fatal(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Fatal(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Fatal(args...)
}

func (l *logrusLogger) Debugln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Debugln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Debugln(args...)
}
func (l *logrusLogger) Infoln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Infoln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Infoln(args...)

}
func (l *logrusLogger) Println(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Println(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Println(args...)
}
func (l *logrusLogger) Warnln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Warnln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Warnln(args...)
}
func (l *logrusLogger) Warningln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Warningln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Warningln(args...)
}
func (l *logrusLogger) Errorln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Errorln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Errorln(args...)
}
func (l *logrusLogger) Panicln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Panicln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Panicln(args...)
}
func (l *logrusLogger) Fatalln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Fatalln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Fatalln(args...)
}

// NOTE: If the `l` logger and the logger from context both has custom field(s) (set by WithField)
// then the field(s) of `l` logger will be ignored, only fields of context's logger will be
// written out.
// We can merge the map fields of two loggers before writing, but it requires more works
// in log function and I don't find it really necessary,
// so I just ignore it to gain performance.
func (l *logrusLogger) TDebugf(ctx context.Context, format string, args ...interface{}) {
	if id := middleware.GetReqID(ctx); id != "" {
		l.Entry = l.Entry.WithFields(logrus.Fields{"RequestID": id})
	}
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		//lg.Debugf(l.getPrefixedFormat(format), args...)
		lg.Debugf(format, args...)
		return
	}
	//l.Debugf(l.getPrefixedFormat(format), args...)
	l.Debugf(format, args...)
}
func (l *logrusLogger) TInfof(ctx context.Context, format string, args ...interface{}) {
	if id := middleware.GetReqID(ctx); id != "" {
		l.Entry = l.Entry.WithFields(logrus.Fields{"RequestID": id})
	}
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		//lg.Infof(l.getPrefixedFormat(format), args...)
		lg.Infof(format, args...)
		return
	}
	//l.Infof(l.getPrefixedFormat(format), args...)
	l.Infof(format, args...)
}
func (l *logrusLogger) TPrintf(ctx context.Context, format string, args ...interface{}) {
	if id := middleware.GetReqID(ctx); id != "" {
		l.Entry = l.Entry.WithFields(logrus.Fields{"RequestID": id})
	}
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		//lg.Printf(l.getPrefixedFormat(format), args...)
		lg.Printf(format, args...)
		return
	}
	//l.Printf(l.getPrefixedFormat(format), args...)
	l.Printf(format, args...)
}
func (l *logrusLogger) TWarnf(ctx context.Context, format string, args ...interface{}) {
	if id := middleware.GetReqID(ctx); id != "" {
		l.Entry = l.Entry.WithFields(logrus.Fields{"RequestID": id})
	}
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		//lg.Warnf(l.getPrefixedFormat(format), args...)
		lg.Warnf(format, args...)
		return
	}
	//l.Warnf(l.getPrefixedFormat(format), args...)
	l.Warnf(format, args...)
}
func (l *logrusLogger) TWarningf(ctx context.Context, format string, args ...interface{}) {
	if id := middleware.GetReqID(ctx); id != "" {
		l.Entry = l.Entry.WithFields(logrus.Fields{"RequestID": id})
	}
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		//lg.Warningf(l.getPrefixedFormat(format), args...)
		lg.Warningf(format, args...)
		return
	}
	//l.Warningf(l.getPrefixedFormat(format), args...)
	l.Warningf(format, args...)
}
func (l *logrusLogger) TErrorf(ctx context.Context, format string, args ...interface{}) {
	if id := middleware.GetReqID(ctx); id != "" {
		l.Entry = l.Entry.WithFields(logrus.Fields{"RequestID": id})
	}
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		//lg.Errorf(l.getPrefixedFormat(format), args...)
		lg.Errorf(format, args...)
		return
	}
	//l.Errorf(l.getPrefixedFormat(format), args...)
	l.Errorf(format, args...)
}
func (l *logrusLogger) TPanicf(ctx context.Context, format string, args ...interface{}) {
	if id := middleware.GetReqID(ctx); id != "" {
		l.Entry = l.Entry.WithFields(logrus.Fields{"RequestID": id})
	}
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		//lg.Panicf(l.getPrefixedFormat(format), args...)
		lg.Panicf(format, args...)
		return
	}
	//l.Panicf(l.getPrefixedFormat(format), args...)
	l.Panicf(format, args...)
}
func (l *logrusLogger) TFatalf(ctx context.Context, format string, args ...interface{}) {
	if id := middleware.GetReqID(ctx); id != "" {
		l.Entry = l.Entry.WithFields(logrus.Fields{"RequestID": id})
	}
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		//lg.Fatalf(l.getPrefixedFormat(format), args...)
		lg.Fatalf(format, args...)
		return
	}
	//l.Fatalf(l.getPrefixedFormat(format), args...)
	l.Fatalf(format, ctx)
}

func (l *logrusLogger) TDebug(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Debug(prefixHelper(l.prefix, args)...)
		return
	}
	l.Debug(args...)
}
func (l *logrusLogger) TInfo(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Info(prefixHelper(l.prefix, args)...)
		return
	}
	l.Info(args...)
}
func (l *logrusLogger) TPrint(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Print(prefixHelper(l.prefix, args)...)
		return
	}
	l.Print(args...)
}
func (l *logrusLogger) TWarn(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Warn(prefixHelper(l.prefix, args)...)
		return
	}
	l.Warn(args...)
}
func (l *logrusLogger) TWarning(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Warning(prefixHelper(l.prefix, args)...)
		return
	}
	l.Warning(args...)
}
func (l *logrusLogger) TError(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Error(prefixHelper(l.prefix, args)...)
		return
	}
	l.Error(args...)
}
func (l *logrusLogger) TPanic(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Panic(prefixHelper(l.prefix, args)...)
		return
	}
	l.Panic(args...)
}
func (l *logrusLogger) TFatal(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Fatal(prefixHelper(l.prefix, args)...)
		return
	}
	l.Fatal(args...)
}

func (l *logrusLogger) TDebugln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Debugln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Debugln(args...)
}
func (l *logrusLogger) TInfoln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Infoln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Infoln(args...)

}
func (l *logrusLogger) TPrintln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Println(prefixHelper(l.prefix, args)...)
		return
	}
	l.Println(args...)

}
func (l *logrusLogger) TWarnln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Warnln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Warnln(args...)
}
func (l *logrusLogger) TWarningln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Warningln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Warningln(args...)
}
func (l *logrusLogger) TErrorln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Errorln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Errorln(args...)
}
func (l *logrusLogger) TPanicln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Panicln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Panicln(args...)
}
func (l *logrusLogger) TFatalln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Fatalln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Fatalln(args...)
}

func (l *logrusLogger) getLoggerFromContext(ctx context.Context) *logrusLogger {
	if ctx == nil {
		return nil
	}

	lg := ctx.Value(LoggerT)
	if lg == nil {
		return nil
	}
	if _, ok := lg.(*logrusLogger); !ok {
		return nil
	}
	return lg.(*logrusLogger)
}

func (l *logrusLogger) getPrefixedFormat(format string) string {
	if l.prefix != "" {
		return l.prefix + ": " + format
	}
	return format
}

func prefixHelper(prefix interface{}, s []interface{}) []interface{} {
	if len(s) == 0 {
		return []interface{}{prefix}
	}
	s[0] = fmt.Sprintf("%s: %v", prefix, s[0])
	return s
}
