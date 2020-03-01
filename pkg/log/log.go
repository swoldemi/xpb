// Package log provides support for logging to stdout and stderr.
// Original source: https://github.com/kelseyhightower/confd/blob/master/log/log.go
package log // import "github.com/swoldemi/xpb/pkg/log"

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// XPBFormatter encapsulates methods to satify the logrus.Formatter interface.
type XPBFormatter struct{}

// Format returns the byte string of how each logged message is formatted.
func (x *XPBFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Format(time.RFC3339)
	hostname, _ := os.Hostname()
	return []byte(fmt.Sprintf("%s %s[%d]: %s %s\n", timestamp, hostname, os.Getpid(), strings.ToUpper(entry.Level.String()), entry.Message)), nil
}

func init() {
	log.SetFormatter(&XPBFormatter{})
}

// SetLevel sets the log level.
func SetLevel(level string) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		Fatal("not a valid log level: '%s'", level)
	}
	log.SetLevel(lvl)
}

// Trace logs a message with severity TRACE.
func Trace(format string, v ...interface{}) {
	log.Trace(fmt.Sprintf(format, v...))
}

// Debug logs a message with severity DEBUG.
func Debug(format string, v ...interface{}) {
	log.Debug(fmt.Sprintf(format, v...))
}

// Error logs a message with severity ERROR.
func Error(format string, v ...interface{}) {
	log.Error(fmt.Sprintf(format, v...))
}

// Fatal logs a message with severity ERROR followed by a call to os.Exit().
func Fatal(format string, v ...interface{}) {
	log.Fatal(fmt.Sprintf(format, v...))
}

// Info logs a message with severity INFO.
func Info(format string, v ...interface{}) {
	log.Info(fmt.Sprintf(format, v...))
}

// Warning logs a message with severity WARNING.
func Warning(format string, v ...interface{}) {
	log.Warning(fmt.Sprintf(format, v...))
}
