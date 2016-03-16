package log

import (
	"os"
	"strings"
	"sync"

	"github.com/zimmski/logrus"
)

var log = logrus.New()
var logIndentation int
var logIndentationLock sync.RWMutex

func init() {
	log.Formatter = &TextFormatter{}
	log.Out = os.Stderr
}

// setting functions

// Level sets the current log level
func Level(level logrus.Level) {
	log.Level = level
}

// LevelDebug sets the current log level to Debug
func LevelDebug() {
	Level(logrus.DebugLevel)
}

// LevelInfo sets the current log level to Info
func LevelInfo() {
	Level(logrus.InfoLevel)
}

// LevelWarn sets the current log level to Warn
func LevelWarn() {
	Level(logrus.WarnLevel)
}

// LevelError sets the current log level to Error
func LevelError() {
	Level(logrus.ErrorLevel)
}

// indentation functions

func indentation() string {
	logIndentationLock.RLock()
	defer logIndentationLock.RUnlock()

	return strings.Repeat("  ", logIndentation)
}

// IncreaseIndentation increase the log indentation.
func IncreaseIndentation() {
	logIndentationLock.Lock()
	defer logIndentationLock.Unlock()

	logIndentation++

}

// DecreaseIndentation decrease the log indentation.
func DecreaseIndentation() {
	logIndentationLock.Lock()
	defer logIndentationLock.Unlock()

	logIndentation--

	if logIndentation < 0 {
		panic("Log indentation is negative")
	}
}

// logging functions

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	log.Debugf(indentation()+format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	log.Infof(indentation()+format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	log.Printf(indentation()+format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	log.Warnf(indentation()+format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	log.Warnf(indentation()+format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	log.Errorf(indentation()+format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(indentation()+format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	log.Panicf(indentation()+format, args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	log.Debug(append([]interface{}{indentation()}, args...)...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	log.Info(append([]interface{}{indentation()}, args...)...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	log.Info(append([]interface{}{indentation()}, args...)...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	log.Warn(append([]interface{}{indentation()}, args...)...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	log.Warn(append([]interface{}{indentation()}, args...)...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	log.Error(append([]interface{}{indentation()}, args...)...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	log.Fatal(append([]interface{}{indentation()}, args...)...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	log.Panic(append([]interface{}{indentation()}, args...)...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	log.Debugln(append([]interface{}{indentation()}, args...)...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	log.Infoln(append([]interface{}{indentation()}, args...)...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	log.Println(append([]interface{}{indentation()}, args...)...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	log.Warnln(append([]interface{}{indentation()}, args...)...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	log.Warnln(append([]interface{}{indentation()}, args...)...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	log.Errorln(append([]interface{}{indentation()}, args...)...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	log.Fatalln(append([]interface{}{indentation()}, args...)...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	log.Panicln(append([]interface{}{indentation()}, args...)...)
}
