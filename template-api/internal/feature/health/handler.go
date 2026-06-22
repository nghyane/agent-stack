// Package health is a minimal feature slice: a single endpoint, no service or
// DB queries of its own. It shows the smallest shape a feature can take.
package health

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler reports service health and verifies DB connectivity.
type Handler struct {
	pool *pgxpool.Pool
}

// NewHandler builds the health handler.
func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

// RegisterRoutes mounts the health endpoint.
func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check",
		Tags:        []string{"system"},
	}, h.check)
}

type healthOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

func (h *Handler) check(ctx context.Context, _ *struct{}) (*healthOutput, error) {
	out := &healthOutput{}
	if h.pool == nil || h.pool.Ping(ctx) != nil {
		out.Body.Status = "degraded"
		return out, nil
	}
	out.Body.Status = "ok"
	return out, nil
}
