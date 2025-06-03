package common

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Warnf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type FakeLogger struct{}

func (fakeLogger *FakeLogger) Warnf(template string, args ...interface{})  {}
func (fakeLogger *FakeLogger) Infof(template string, args ...interface{})  {}
func (fakeLogger *FakeLogger) Fatalf(template string, args ...interface{}) {}
func (fakeLogger *FakeLogger) Errorf(template string, args ...interface{}) {}
func (fakeLogger *FakeLogger) Debugf(template string, args ...interface{}) {}

func GetComponentLogger(lvl string) Logger {
	var level zap.AtomicLevel
	switch lvl {
	case "debug", "Debug", "DEBUG":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info", "Info", "INFO":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn", "WARN":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	default:
		return &FakeLogger{}
	}

	var core zapcore.Core
	coreConsole := getComponentLoggerConsole(level)
	core = zapcore.NewTee(
		coreConsole,
	)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}
func getComponentLoggerConsole(level zap.AtomicLevel) zapcore.Core {
	configLogConsole := zap.NewDevelopmentEncoderConfig()
	// configLogConsole := zap.NewProductionEncoderConfig()
	configLogConsole.ConsoleSeparator = " | "
	configLogConsole.EncodeTime = zapcore.ISO8601TimeEncoder
	configLogConsole.EncodeLevel = zapcore.CapitalColorLevelEncoder
	configLogConsole.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConsole := zapcore.NewConsoleEncoder(configLogConsole)
	stdout := zapcore.AddSync(os.Stdout)
	return zapcore.NewCore(encoderConsole, stdout, level)
}
