// Package log is a simple wrapper over golang log logger for minimal log level support.
//
// - Debug logs will be logged only if debug flag is true.
// - Error and Warn logs will have ERROR and WARN preprended.
//
package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	debug = false
)

func init() {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		SetLevel("DEBUG")
	}
}

// SetLevel - use "DEBUG" to enable debug logging
func SetLevel(l string) {
	debug = (l == "DEBUG")
	log.Printf("Set log level debug = %t", debug)
}

// IsDebug returns true if log level is DEBUG.
func IsDebug() bool {
	return debug
}

// Printf - These logs cannot be suppresed.
func Printf(format string, l ...interface{}) {
	log.Printf("INFO %s", fmt.Sprintf(format, l...))
}

// Debugf - These logs will be suppresed unless log level=DEBUG.
func Debugf(format string, l ...interface{}) {
	if debug {
		log.Printf("DEBUG %s", fmt.Sprintf(format, l...))
	}
}

// Errorf logs with ERROR preprended. These logs cannot be suppressed.
func Errorf(format string, l ...interface{}) {
	log.Printf("ERROR %s", fmt.Sprintf(format, l...))
}

// Warnf logs with WARN preprended. These logs cannot be suppressed.
func Warnf(format string, l ...interface{}) {
	log.Printf("WARN %s", fmt.Sprintf(format, l...))
}

// Fatalf logs with FATAL preprended and exit.
func Fatalf(format string, l ...interface{}) {
	log.Fatalf("FATAL %s", fmt.Sprintf(format, l...))
}

// T returns a logger which prefixes tenantID and optional fields to log messages.
// Example:
//     log.T("23", "req-id", "p244").Printf("Updated endpoint profile")
//                    prints
//     2017/12/14 11:41:38 INFO [t=23 req-id=p244] Updated endpoint profile
func T(tenantID string, fields ...string) ContextLogger {
	return ContextLogger{context{append([]string{"t", tenantID}, fields...)}}
}

// With is used to get a ContextLogger with context fields set. It can be used to
// perform logging in different levels.
// Example:
//  	log := log.With("12", "req_id", "p244", "thread_id": 10)
//  	log.Errorf("Failed to process request")
//             prints
//      2017/12/14 14:15:37 ERROR [req_id=p244 thread_id=10] Failed to process request
func With(fields ...string) ContextLogger {
	return ContextLogger{context{fields}}
}

// Context interface can be implemented to add prefix string to log messages.
type Context interface {
	// Prefix - string build form context used as prefix of log msg.
	Prefix() string
}

type ContextLogger struct {
	Context
}

func (ctxLog ContextLogger) Printf(format string, l ...interface{}) {
	Printf(ctxLog.Prefix()+format, l...)
}

func (ctxLog ContextLogger) Debugf(format string, l ...interface{}) {
	Debugf(ctxLog.Prefix()+format, l...)
}

func (ctxLog ContextLogger) Errorf(format string, l ...interface{}) {
	Errorf(ctxLog.Prefix()+format, l...)
}

func (ctxLog ContextLogger) Warnf(format string, l ...interface{}) {
	Warnf(ctxLog.Prefix()+format, l...)
}

func (ctxLog ContextLogger) Fatalf(format string, l ...interface{}) {
	Fatalf(ctxLog.Prefix()+format, l...)
}

// context is a implementation of Context which creates a prefix by combining fields
// as [fields1=value1 field2=value2...]
type context struct {
	fields []string
}

func (ctx context) Prefix() string {
	if len(ctx.fields) == 0 || len(ctx.fields)%2 != 0 {
		return ""
	}

	var fields []string
	for i := 0; i < len(ctx.fields); i += 2 {
		fields = append(fields, fmt.Sprintf("%s=%s", ctx.fields[i], ctx.fields[i+1]))
	}

	return fmt.Sprintf("[%s] ", strings.Join(fields, " "))
}
