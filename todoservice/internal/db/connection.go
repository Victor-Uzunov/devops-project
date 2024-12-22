package db

import (
	"context"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type Config struct {
	Host       string        `envconfig:"APP_DB_HOST"`
	Port       string        `envconfig:"APP_DB_PORT"`
	User       string        `envconfig:"APP_DB_USER"`
	Password   string        `envconfig:"APP_DB_PASSWORD"`
	Database   string        `envconfig:"APP_DB_NAME"`
	SSLMode    string        `envconfig:"APP_DB_SSLMODE"`
	RetryCount int           `envconfig:"APP_DB_RETRYCOUNT"`
	Duration   time.Duration `envconfig:"APP_DB_DURATION"`
}

func Create(ctx context.Context, config Config) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", config.User, config.Password,
		config.Host, config.Port, config.Database, config.SSLMode)
	log.C(ctx).Info("connecting to postgresql", conStr)

	var persistence *sqlx.DB
	var err error
	for i := 0; i < config.RetryCount; i++ {
		persistence, err = sqlx.Open("postgres", conStr)
		if err != nil {
			log.C(ctx).Debugf("postgresql connection failed after %d retries for url: %s", config.RetryCount, conStr)
			time.Sleep(config.Duration)
			continue
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second)
		err = persistence.PingContext(ctx)
		cancel()
		if err != nil {
			fmt.Println(err)
			time.Sleep(config.Duration)
			continue
		}
		fmt.Println("Connected to database!")
		return persistence, err
	}
	return nil, err
}
