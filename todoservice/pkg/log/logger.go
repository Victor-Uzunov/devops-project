package log

import (
	"context"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Level  string `envconfig:"APP_LOG_LEVEL" default:"debug"`
	Format string `envconfig:"APP_LOG_FORMAT" default:"text"`
}

type logKey struct{}

var (
	mutex = sync.RWMutex{}
	C     = LoggerFromContext
)

func init() {
	_, err := SetupLogger(context.Background(), Config{Level: "info", Format: "text"})
	if err != nil {
		panic(err)
	}
}

func SetupLogger(ctx context.Context, cfg Config) (context.Context, error) {
	mutex.Lock()
	defer mutex.Unlock()

	logEntry := logrus.NewEntry(logrus.StandardLogger())

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	logEntry.Logger.SetLevel(level)
	logEntry.Logger.SetOutput(os.Stdout)

	switch cfg.Format {
	case "json":
		logEntry.Logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logEntry.Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return ContextWithLogger(ctx, logEntry), nil
}

func ContextWithLogger(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, logKey{}, entry)
}

func LoggerFromContext(ctx context.Context) *logrus.Entry {
	mutex.RLock()
	defer mutex.RUnlock()
	entry := ctx.Value(logKey{})
	if entry == nil {
		entry = logrus.NewEntry(logrus.StandardLogger())
	}
	return entry.(*logrus.Entry)
}
