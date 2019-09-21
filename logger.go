package logger

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var logFn func(format string, args ...interface{})

// contextKey is type for context key.
type contextKey string

// contextValue is type for context v`alue.
type contextValue string

const logIDKey contextKey = "__log_id__"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetContextWithLogID is used to setup context
// and set log ID into it.
func GetContextWithLogID(ctx context.Context, logID string) context.Context {
	token := randomString(8)
	logID = fmt.Sprintf("%s_%s", logID, token)
	return GetContextWithNoSubfixLogID(ctx, logID)

}

// GetContextWithNoSubfixLogID is used to setup context
// and set log ID into it without subfix added.
func GetContextWithNoSubfixLogID(ctx context.Context, logID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, logIDKey, contextValue(logID))
}

// GetLogID is used to get log ID
// from context.
func GetLogID(ctx context.Context) string {
	logID, ok := ctx.Value(logIDKey).(contextValue)
	if !ok {
		return ""
	}
	return string(logID)
}

// SetupLogger is used to setup logging function.
func SetupLogger(fn func(format string, args ...interface{})) {
	logFn = fn
}

// Infof is used to log information message.
func Infof(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, "INFO", format, args...)
}

// Warnf is used to log warning message.
func Warnf(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, "WARN", format, args...)
}

// Debugf is used to log warning message.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, "DEBUG", format, args...)
}

// Errorf is used to log error message.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, "ERROR", format, args...)
}

// Fatalf is used to log error message
// then call os.Exit(1).
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	printf(ctx, "FATAL", format, args...)
	os.Exit(1)
}

func getDefaultLogFn() func(format string, args ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	return func(format string, args ...interface{}) {
		log.Printf(fmt.Sprintf("\b\b%s", format), args...)
	}
}

func printf(ctx context.Context, mode, format string, args ...interface{}) {
	_, filePath, line, _ := runtime.Caller(2)
	path := strings.Split(filePath, "/")
	file := path[len(path)-1]
	if logFn == nil {
		logFn = getDefaultLogFn()
	}
	fileLine := fmt.Sprintf("%s:%d", file, line)
	logFormat := fmt.Sprintf("|%s|%s|%s", mode, fileLine, format)
	if ctx != nil {
		logID, ok := ctx.Value(logIDKey).(contextValue)
		if ok {
			logFormat = fmt.Sprintf("|%s|%s|log_id=%s|%s", mode, fileLine, logID, format)
		}
	}
	logFn(logFormat, args...)
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
