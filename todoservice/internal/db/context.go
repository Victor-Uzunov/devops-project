package db

import (
	"context"
	"errors"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/jmoiron/sqlx"
)

const (
	dataBaseCtx = "dataBaseCtx"
)

func SaveToContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, dataBaseCtx, tx)
}

func FromContext(ctx context.Context) (*sqlx.Tx, error) {
	log.C(ctx).Info("getting database from context in method FromContext")
	tx, ok := ctx.Value(dataBaseCtx).(*sqlx.Tx)
	if !ok {
		log.C(ctx).Error("database not in context")
		return nil, errors.New("database not in context")
	}
	return tx, nil
}
