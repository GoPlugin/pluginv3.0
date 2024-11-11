package config

import (
	"go.uber.org/zap/zapcore"

	"github.com/goplugin/pluginv3.0/v2/core/utils"
)

type File interface {
	Dir() string
	MaxSize() utils.FileSize
	MaxAgeDays() int64
	MaxBackups() int64
}

type Log interface {
	DefaultLevel() zapcore.Level
	JSONConsole() bool
	Level() zapcore.Level
	UnixTimestamps() bool

	File() File
}
