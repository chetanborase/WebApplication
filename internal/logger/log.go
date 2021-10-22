package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
)

var (
	logger *zap.Logger
	Log    = log.New(os.Stderr, "line no ", log.Lshortfile)
)

func init() {
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.ErrorLevel
	})
	llInfo := &lumberjack.Logger{
		Filename:   infoFilePath(),
		MaxSize:    maxSize(),         //2 MB size
		MaxAge:     maxAgeOfLogFile(), //30 days
		MaxBackups: maxBackup(),
		LocalTime:  false,
		Compress:   false,
	}
	llErr := &lumberjack.Logger{
		Filename:   errFilePath(),
		MaxSize:    maxSize(),         //2 MB size
		MaxAge:     maxAgeOfLogFile(), //30 days
		MaxBackups: maxBackup(),
		LocalTime:  false,
		Compress:   false,
	}
	topicInfo := zapcore.AddSync(llInfo)
	topicErr := zapcore.AddSync(io.MultiWriter(llInfo, llErr))
	topicDebug := zapcore.AddSync(os.Stderr)
	//consoleInfo := zapcore.Lock(os.Stdout)
	consoleErr := zapcore.Lock(os.Stderr)
	logEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(logEncoder, topicInfo, lowPriority),
		zapcore.NewCore(logEncoder, topicErr, highPriority),
		//zapcore.NewCore(consoleEncoder, consoleInfo, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleErr, highPriority),
		zapcore.NewCore(consoleEncoder, topicDebug, lowPriority),
	)
	logger = zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
}
func Info(msg string, field ...zap.Field) {
	logger.Info(msg, field...)
}
func Error(err error, field ...zap.Field) {
	logger.Error(err.Error(), field...)
}
func Debug(msg string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(msg, args...))
}
func Println(a ...interface{}) {
	logger.Info(fmt.Sprint(a...))
}
func Print(a ...interface{}) {
	Println(a...)
}
