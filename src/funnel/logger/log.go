package logger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"strings"
)

const (
  DebugLevel string = "debug"
  InfoLevel = "info"
  ErrorLevel = "error"
)

var formatter = &textFormatter{
	DisableTimestamp: true,
}

func Configure(conf config.Config) {
	logrus.SetFormatter(formatter)
	logger.SetLevel(conf.LogLevel)

	// TODO Good defaults, configuration, and reusable way to configure logging.
	//      Also, how do we get this to default to /var/log/tes/worker.log
	//      without having file permission problems? syslog?
	// Configure logging
	if conf.LogPath != "" {
		logFile, err := os.OpenFile(
			conf.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666,
		)
		if err != nil {
			logger.Error("Can't open log output file", "path", conf.LogPath)
		} else {
			logger.SetOutput(logFile)
		}
	}
}

// Logger is repsonsible for logging messages from code.
type Logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Error(string, ...interface{})
	WithFields(...interface{}) Logger
}

// SetLevel sets the level of logging
func SetLevel(l string) {
	switch strings.ToLower(l) {
	case DebugLevel:
		logrus.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		logrus.SetLevel(logrus.InfoLevel)
	case ErrorLevel:
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// ForceColors forces the log output formatter to use color. Useful during testing.
func ForceColors() {
	formatter.ForceColors = true
}

// New returns a new Logger instance.
func New(ns string, args ...interface{}) Logger {
	f := fields(args...)
	f["ns"] = ns
	l := logrus.WithFields(f)
	return &logger{l}
}

type logger struct {
	log *logrus.Entry
}

// Debug logs a debug message.
//
// After the first argument, arguments are key-value pairs which are written as structured logs.
//     log.Debug("Some message here", "key1", value1, "key2", value2)
func (l *logger) Debug(msg string, args ...interface{}) {
	defer recoverLogErr()
	f := fields(args...)
	l.log.WithFields(f).Debug(msg)
}

// Info logs an info message
//
// After the first argument, arguments are key-value pairs which are written as structured logs.
//     log.Info("Some message here", "key1", value1, "key2", value2)
func (l *logger) Info(msg string, args ...interface{}) {
	defer recoverLogErr()
	f := fields(args...)
	l.log.WithFields(f).Info(msg)
}

// Error logs an error message
//
// After the first argument, arguments are key-value pairs which are written as structured logs.
//     log.Error("Some message here", "key1", value1, "key2", value2)
//
// Error has a two-argument version that can be used as a shortcut.
//     err := startServer()
//     log.Error("Couldn't start server", err)
func (l *logger) Error(msg string, args ...interface{}) {
	defer recoverLogErr()
	var f map[string]interface{}
	if len(args) == 1 {
		f = fields("error", args[0])
	} else {
		f = fields(args...)
	}
	l.log.WithFields(f).Error(msg)
}

// WithFields returns a new Logger instance with the given fields added to all log messages.
func (l *logger) WithFields(args ...interface{}) Logger {
	defer recoverLogErr()
	f := fields(args...)
	n := l.log.WithFields(f)
	return &logger{n}
}

// SetOutput sets the output for all loggers.
func SetOutput(w io.Writer) {
	logrus.SetOutput(w)
}

var rootLogger = New("tes")

// Debug logs to the global logger at the Debug level
func Debug(msg string, args ...interface{}) {
	rootLogger.Debug(msg, args...)
}

// Info logs to the global logger at the Info level
func Info(msg string, args ...interface{}) {
	rootLogger.Info(msg, args...)
}

// Error logs to the global logger at the Error level
func Error(msg string, args ...interface{}) {
	rootLogger.Error(msg, args...)
}

// recoverLogErr is used to recover from any panics during logging.
// Panics aren't expected of course, but logging should never crash
// a program, so this failsafe tries to prevent those crashes.
func recoverLogErr() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from logging panic", r)
	}
}

func fields(args ...interface{}) map[string]interface{} {
	f := make(map[string]interface{}, len(args)/2)
	if len(args) == 1 {
		f["unknown"] = args[0]
		return f
	}
	for i := 0; i < len(args); i += 2 {
		k := args[i].(string)
		v := args[i+1]
		f[k] = v
	}
	if len(args)%2 != 0 {
		f["unknown"] = args[len(args)-1]
	}
	return f
}
