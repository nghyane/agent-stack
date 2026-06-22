package settings

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Manager keeps an in-memory snapshot of settings and keeps it fresh by
// listening for Postgres NOTIFY. Reads (Current) are lock-free; writes (Update)
// persist and let the notify fan out to every instance, including this one.
type Manager struct {
	store   *Store
	pool    *pgxpool.Pool
	current atomic.Pointer[Settings]
}

// NewManager loads the initial snapshot and starts the listener goroutine.
// The goroutine runs until ctx is cancelled.
func NewManager(ctx context.Context, pool *pgxpool.Pool) (*Manager, error) {
	m := &Manager{store: NewStore(pool), pool: pool}

	if err := m.reload(ctx); err != nil {
		return nil, err
	}

	go m.listen(ctx)
	return m, nil
}

// Current returns the latest in-memory snapshot. Lock-free; safe for hot paths.
func (m *Manager) Current() Settings {
	return *m.current.Load()
}

// Update applies fn to a copy of the current settings, persists it, and
// notifies all instances. The local snapshot refreshes via the notify it emits.
func (m *Manager) Update(ctx context.Context, fn func(*Settings)) error {
	next := m.Current()
	fn(&next)
	return m.store.Save(ctx, next)
}

// reload reads the row and swaps the in-memory snapshot.
func (m *Manager) reload(ctx context.Context) error {
	s, err := m.store.Load(ctx)
	if err != nil {
		return err
	}
	m.current.Store(&s)
	return nil
}

// listen holds a dedicated connection in LISTEN mode and reloads on every
// notification. It reloads once on each (re)connect too, so changes missed
// while disconnected are still picked up. On connection loss it reconnects
// with backoff until ctx is cancelled.
func (m *Manager) listen(ctx context.Context) {
	const (
		minBackoff = 1 * time.Second
		maxBackoff = 30 * time.Second
	)
	backoff := minBackoff

	for ctx.Err() == nil {
		if err := m.listenOnce(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("settings listener dropped; reconnecting",
				"error", err, "retry_in", backoff)
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
			backoff = min(backoff*2, maxBackoff)
			continue
		}
		backoff = minBackoff
	}
}

// listenOnce acquires a connection, LISTENs, reloads once to catch up, then
// blocks on notifications until the connection drops or ctx is cancelled.
func (m *Manager) listenOnce(ctx context.Context) error {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, "LISTEN "+notifyChannel); err != nil {
		return err
	}

	// Catch up on anything changed while we were not listening.
	if err := m.reload(ctx); err != nil {
		slog.Warn("settings reload on connect failed", "error", err)
	}

	for {
		if _, err := conn.Conn().WaitForNotification(ctx); err != nil {
			return err
		}
		if err := m.reload(ctx); err != nil {
			slog.Warn("settings reload on notify failed", "error", err)
		}
	}
}
