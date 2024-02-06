package logger

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	defaultLogger Logger
	mu            sync.RWMutex
)

type log struct{}

func init() {
	defaultLogger = &logger{fl: logrus.New()}
}

var C = From

func SetDefault(logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	defaultLogger = logger
}

func Set(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, log{}, logger)
}

func From(ctx context.Context) Logger {
	log, ok := ctx.Value(log{}).(Logger)
	if ok {
		return log
	}
	if defaultLogger != nil {
		return defaultLogger
	}
	logr := New()
	mu.Lock()
	defer mu.Unlock()
	defaultLogger = logr
	return defaultLogger
}
