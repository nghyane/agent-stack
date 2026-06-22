package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nghiahoang/template-api/internal/feature/health"
	settingsfeat "github.com/nghiahoang/template-api/internal/feature/settings"
	"github.com/nghiahoang/template-api/internal/feature/user"
	"github.com/nghiahoang/template-api/internal/platform/config"
	"github.com/nghiahoang/template-api/internal/platform/session"
	"github.com/nghiahoang/template-api/internal/platform/settings"
)

// API metadata shown in the OpenAPI document. Rename per project.
const (
	apiTitle   = "Template API"
	apiVersion = "1.0.0"
)

// feature is anything that can mount its own routes onto the API. Every feature
// slice implements this; the server treats them uniformly.
type feature interface {
	RegisterRoutes(api huma.API)
}

// New builds the HTTP handler: Chi router + Huma API with all features mounted.
func New(cfg *config.Config, pool *pgxpool.Pool, sessions *session.Manager, settingsMgr *settings.Manager) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors(cfg.CORSOrigins))
	// Session middleware must wrap requests so session context is available.
	r.Use(sessions.Wrap)

	api := humachi.New(r, huma.DefaultConfig(apiTitle, apiVersion))
	mountFeatures(api, pool, sessions, settingsMgr)
	return r
}

// OpenAPISpec returns the OpenAPI 3.1 document as YAML, without needing a DB.
// Handlers are constructed with nil deps: spec generation only reads types, it
// never invokes the handler functions.
func OpenAPISpec() ([]byte, error) {
	api := humachi.New(chi.NewMux(), huma.DefaultConfig(apiTitle, apiVersion))
	mountFeatures(api, nil, nil, nil)
	return api.OpenAPI().YAML()
}

// mountFeatures is the single place every feature is wired in.
// Add a feature: construct its handler here and append to the list.
func mountFeatures(api huma.API, pool *pgxpool.Pool, sessions *session.Manager, settingsMgr *settings.Manager) {
	features := []feature{
		health.NewHandler(pool),
		user.NewHandler(pool, sessions),
		settingsfeat.NewHandler(settingsMgr, sessions),
	}
	for _, f := range features {
		f.RegisterRoutes(api)
	}
}
