package settings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// notifyChannel is the Postgres LISTEN/NOTIFY channel for settings changes.
const notifyChannel = "app_settings_changed"

// Store reads and writes the singleton settings row.
type Store struct {
	pool *pgxpool.Pool
}

// NewStore builds a settings store over the pool.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// Load reads the settings row and decodes it onto Defaults, so any key absent
// from the stored JSON keeps its default value.
func (s *Store) Load(ctx context.Context) (Settings, error) {
	var raw []byte
	err := s.pool.QueryRow(ctx, `SELECT data FROM app_settings WHERE id = 1`).Scan(&raw)
	if err != nil {
		return Settings{}, fmt.Errorf("load settings: %w", err)
	}

	out := Defaults()
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &out); err != nil {
			return Settings{}, fmt.Errorf("decode settings: %w", err)
		}
	}
	return out, nil
}

// Save persists the settings and notifies all instances to reload. The write
// and the notify share one transaction so a notify never fires for an
// uncommitted change.
func (s *Store) Save(ctx context.Context, val Settings) error {
	raw, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("encode settings: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		`UPDATE app_settings SET data = $1, updated_at = now() WHERE id = 1`, raw,
	); err != nil {
		return fmt.Errorf("update settings: %w", err)
	}
	// NOTIFY payload is empty: listeners reload the full row, not the payload.
	if _, err := tx.Exec(ctx, fmt.Sprintf("NOTIFY %s", notifyChannel)); err != nil {
		return fmt.Errorf("notify: %w", err)
	}
	return tx.Commit(ctx)
}
