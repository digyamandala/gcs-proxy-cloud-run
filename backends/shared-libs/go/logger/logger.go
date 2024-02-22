package logger

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/commonutils"
	"github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
)

const (
	ReqIDField      = "reqId"
	StackTraceField = "stackTrace"
)

var maxDepth = 3

func SetupLogger(cfg Config) {
	if len(cfg.LogLevel) == 0 {
		cfg.LogLevel = "DEBUG"
	}
	if cfg.MaxDepth > 0 {
		maxDepth = cfg.MaxDepth
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(LogLevelStringToEnum(cfg.LogLevel))

}

func LogLevelStringToEnum(str string) logrus.Level {
	switch strings.ToUpper(str) {
	case "WARN":
		return logrus.WarnLevel
	case "INFO":
		return logrus.InfoLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "TRACE":
		return logrus.TraceLevel
	case "FATAL":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func Info(ctx context.Context, format string, values ...interface{}) {
	logrus.WithFields(getLogrusFields(ctx, values)).Info(fmt.Sprintf(format, values...))
}

func Debug(ctx context.Context, format string, values ...interface{}) {
	logrus.WithFields(getLogrusFields(ctx, values)).Debug(fmt.Sprintf(format, values...))
}

func Warn(ctx context.Context, format string, values ...interface{}) {
	logrus.WithFields(getLogrusFields(ctx, values)).Warn(fmt.Sprintf(format, values...))
}

func Error(ctx context.Context, format string, values ...interface{}) {
	logrus.WithFields(getLogrusFields(ctx, values)).Error(fmt.Sprintf(format, values...))
}

func Fatal(ctx context.Context, format string, values ...interface{}) {
	logrus.WithFields(getLogrusFields(ctx, values)).Fatal(fmt.Sprintf(format, values...))
}

func getLogrusFields(ctx context.Context, values []interface{}) logrus.Fields {
	f := logrus.Fields{
		ReqIDField: commonutils.ReqIDFromContext(ctx),
	}

	var err []error
	for _, v := range values {
		e, ok := v.(error)
		if ok {
			err = append(err, e)
		}
	}

	var stack []trace

	for _, e := range err {
		stack = append(stack, trace{
			ErrMsg: e.Error(),
			Trace:  stackTrace(e),
		})
	}

	if len(stack) > 0 {
		f[StackTraceField] = stack
	}

	return f
}

func stackTrace(err error) []tracerr.Frame {
	frames := tracerr.StackTrace(err)
	length := len(frames)

	if length > maxDepth {
		length = maxDepth
	}

	return frames[:length]
}
