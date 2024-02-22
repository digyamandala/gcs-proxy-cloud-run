package logger

import "github.com/ztrue/tracerr"

type Config struct {
	LogLevel string
	MaxDepth int
}

type trace struct {
	ErrMsg string          `json:"errMsg"`
	Trace  []tracerr.Frame `json:"trace"`
}
