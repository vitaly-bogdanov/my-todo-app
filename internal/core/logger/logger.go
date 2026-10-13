package core_logger

import (
	"fmt"
	"os"
	"path/filepath"

	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerContextKey struct{}

var LoggerContextKey = loggerContextKey{}

// Logger пишет сообщения в файл и в стандартный вывод.
// Методы zap доступны напрямую через встроенный *zap.Logger.
type Logger struct {
	*zap.Logger
	file *os.File
}

// FromContext достаёт Logger из контекста.
// Если логера в контексте нет или тип значения другой, FromContext паникует.
func FromContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value(LoggerContextKey).(*Logger)
	if !ok {
		panic("logger not found in context")
	}
	return logger
}

// NewLogger создаёт Logger по config.
// Уровень берётся из config.Level. Файл создаётся в config.Folder, его имя — текущее время.
// Каждая запись уходит и в этот файл, и в стандартный вывод.
func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log folder: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(config.Folder, fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{Logger: zapLogger, file: logFile}, nil
}

// With возвращает новый Logger с дополнительными полями.
// Файл логов остаётся тем же.
func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(field...), file: l.file}
}

// Close закрывает файл логов.
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close logger:", err)
	}
}
