package utils

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
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
	formatedMessage := fmt.Sprintf(format, args...)
	msgC := message.New(ctx)
	_ = msgC.SendSMS("9316445311", formatedMessage)
	_ = msgC.SendSMS("6155454989", formatedMessage)
	aelog.Criticalf(ctx, formatedMessage)
}

// Errorf logs error messages.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	aelog.Errorf(ctx, format, args...)
}
