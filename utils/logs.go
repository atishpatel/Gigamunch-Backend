package utils

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
)

const (
	REDISLATENCY = "RedisLatency"
)

// Latencyf is used to log metrics for latency
func Latencyf(ctx context.Context, label string, format string, args ...interface{}) {
	aelog.Infof(ctx, label+" : "+format, args...)
}

// Infof logs info
func Infof(ctx context.Context, format string, args ...interface{}) {
	aelog.Infof(ctx, format, args...)
}

// Debugf only logs in development
func Debugf(ctx context.Context, format string, args ...interface{}) {
	if appengine.IsDevAppServer() {
		aelog.Debugf(ctx, format, args...)
	}
}

// Warningf logs warnings
func Warningf(ctx context.Context, format string, args ...interface{}) {
	aelog.Warningf(ctx, format, args...)
}

// Criticalf logs critical messages. Critical means the app is failing.
func Criticalf(ctx context.Context, format string, args ...interface{}) {
	aelog.Criticalf(ctx, format, args...)
}

// Errorf logs error messages.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	aelog.Errorf(ctx, format, args...)
}
