package logger

import (
	"fmt"
	"os"
	"strings"

	"wallpaper-system/internal/infrastructure/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger создает новый структурированный логгер
func NewZapLogger(cfg config.LoggerConfig) (*zap.Logger, error) {
	// Настраиваем уровень логирования
	level, err := parseLogLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("неверный уровень логирования: %w", err)
	}

	// Настраиваем энкодер
	encoderConfig := getEncoderConfig(cfg.Format)
	
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Настраиваем вывод
	writeSyncer, err := getWriteSyncer(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка настройки вывода логгера: %w", err)
	}

	// Создаем core логгера
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Создаем логгер с дополнительными опциями
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

// parseLogLevel парсит строковый уровень логирования
func parseLogLevel(levelStr string) (zapcore.Level, error) {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("неизвестный уровень логирования: %s", levelStr)
	}
}

// getEncoderConfig возвращает конфигурацию энкодера
func getEncoderConfig(format string) zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	
	// Настраиваем временные метки
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	
	// Настраиваем уровни
	config.LevelKey = "level"
	if format == "json" {
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	} else {
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	// Настраиваем caller
	config.CallerKey = "caller"
	config.EncodeCaller = zapcore.ShortCallerEncoder
	
	// Настраиваем имя логгера
	config.NameKey = "logger"
	
	// Настраиваем сообщение
	config.MessageKey = "message"
	
	// Настраиваем stack trace
	config.StacktraceKey = "stacktrace"

	return config
}

// getWriteSyncer настраивает куда направлять логи
func getWriteSyncer(cfg config.LoggerConfig) (zapcore.WriteSyncer, error) {
	var writers []zapcore.WriteSyncer

	// Основной вывод (stdout, stderr или файл)
	switch cfg.Output {
	case "stdout":
		writers = append(writers, zapcore.AddSync(os.Stdout))
	case "stderr":
		writers = append(writers, zapcore.AddSync(os.Stderr))
	default:
		// Если указан путь к файлу
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("не удается открыть файл логов %s: %w", cfg.Output, err)
		}
		writers = append(writers, zapcore.AddSync(file))
	}

	// Дополнительный файловый вывод если включен
	if cfg.EnableFile && cfg.FilePath != "" && cfg.FilePath != cfg.Output {
		// Создаем директорию если она не существует
		if err := os.MkdirAll(getDir(cfg.FilePath), 0755); err != nil {
			return nil, fmt.Errorf("не удается создать директорию для файла логов: %w", err)
		}

		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("не удается открыть файл логов %s: %w", cfg.FilePath, err)
		}
		writers = append(writers, zapcore.AddSync(file))
	}

	if len(writers) == 1 {
		return writers[0], nil
	}

	return zapcore.NewMultiWriteSyncer(writers...), nil
}

// getDir извлекает директорию из пути к файлу
func getDir(filePath string) string {
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			return filePath[:i]
		}
	}
	return "."
}

// NewDevelopmentLogger создает логгер для разработки
func NewDevelopmentLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config.Build()
}

// NewProductionLogger создает логгер для продакшена
func NewProductionLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	return config.Build()
}

// LoggerMiddleware возвращает middleware для логирования HTTP запросов
func LoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Создаем wrapper для response writer чтобы получить статус код
			ww := &responseWriter{ResponseWriter: w, statusCode: 200}
			
			next.ServeHTTP(ww, r)
			
			duration := time.Since(start)
			
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.Int("status", ww.statusCode),
				zap.Duration("duration", duration),
				zap.String("user_agent", r.UserAgent()),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

// responseWriter обертка для http.ResponseWriter чтобы получить статус код
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
} 