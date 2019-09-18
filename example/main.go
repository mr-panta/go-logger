package main

import (
	"context"
	"errors"

	"github.com/mr-panta/go-logger"
)

func main() {
	ctx := logger.GetContextWithLogID(context.Background(), "example")
	logger.Infof(ctx, "this is information message, %d", 1234)
	logger.Warnf(ctx, "this is warning message, %t", false)
	logger.Debugf(ctx, "this is debug message, %v", ctx)
	logger.Errorf(ctx, "this is error message, %v", errors.New("bye guys"))
	logger.Fatalf(ctx, "this is fatal message, %f", 3.14)
}
