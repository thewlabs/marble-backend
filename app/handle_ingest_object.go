package app

import (
	"context"

	"golang.org/x/exp/slog"
)

func (app *App) IngestObject(ctx context.Context, payload Payload, table Table, logger *slog.Logger) (err error) {
	return app.repository.IngestObject(ctx, payload, table, logger)
}
