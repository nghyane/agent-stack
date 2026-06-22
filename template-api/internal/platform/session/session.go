// Package session manages the authenticated user's session. It is platform
// infrastructure (cross-cutting): any feature can read "who is logged in" via
// the helpers here. Sessions are stored server-side in Postgres (scs); the
// cookie holds only a random token.
package session

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// userKey is the scs session key storing the authenticated user ID.
const userKey = "user_id"

// Manager wraps scs with the project's cookie policy and typed helpers.
type Manager struct {
	scs *scs.SessionManager
}

// NewManager builds a session manager backed by Postgres.
func NewManager(pool *pgxpool.Pool, secure bool) *Manager {
	sm := scs.New()
	sm.Store = pgxstore.New(pool)
	sm.Lifetime = 7 * 24 * time.Hour
	sm.Cookie.Name = "session"
	sm.Cookie.HttpOnly = true
	sm.Cookie.Secure = secure
	// SameSite=None is required for cross-site SPA (Cloudflare Pages -> API).
	// Only valid alongside Secure=true, so fall back to Lax in dev.
	if secure {
		sm.Cookie.SameSite = http.SameSiteNoneMode
	} else {
		sm.Cookie.SameSite = http.SameSiteLaxMode
	}
	return &Manager{scs: sm}
}

// Wrap returns the middleware that loads/saves the session per request.
// It must wrap every request so session context is available to handlers.
func (m *Manager) Wrap(next http.Handler) http.Handler {
	return m.scs.LoadAndSave(next)
}

// Login stores the user ID in the session and renews the token.
func (m *Manager) Login(ctx context.Context, userID string) error {
	if err := m.scs.RenewToken(ctx); err != nil {
		return err
	}
	m.scs.Put(ctx, userKey, userID)
	return nil
}

// Logout clears the current session.
func (m *Manager) Logout(ctx context.Context) error {
	return m.scs.Destroy(ctx)
}

// UserID returns the authenticated user ID, or "" if not logged in.
func (m *Manager) UserID(ctx context.Context) string {
	return m.scs.GetString(ctx, userKey)
}
