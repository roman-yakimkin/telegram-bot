package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
)

type DBConnect interface {
	Connect(ctx context.Context) (*pgxpool.Pool, error)
	Disconnect(ctx context.Context)
}

type dBConnect struct {
	service *config.Service
	pool    *pgxpool.Pool
}

func NewDBConnect(service *config.Service) DBConnect {
	return &dBConnect{
		service: service,
	}
}

func (c *dBConnect) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	c.pool, err = pgxpool.Connect(ctx, c.service.GetConfig().DBConnect)
	if err != nil {
		return nil, err
	}
	return c.pool, nil
}

func (c *dBConnect) Disconnect(_ context.Context) {
	c.pool.Close()
}
