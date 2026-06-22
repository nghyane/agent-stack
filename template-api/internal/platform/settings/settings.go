// Package settings holds runtime-mutable operational config, backed by a single
// Postgres row and propagated across instances via LISTEN/NOTIFY.
//
// This is the DYNAMIC tier of configuration. It is for operational toggles that
// change WITHOUT a redeploy (e.g. which AI model is active, a feature flag).
// It is NOT for secrets or wiring — those are static config (env), read once at
// startup. Never put an API key or connection string here.
//
// Reads are lock-free from an in-memory snapshot (see Manager.Current). Writes
// persist to the row and notify every instance to reload.
package settings

// Settings is the typed shape of the operational config. Stored as JSONB, so
// adding a field needs no migration: add it here with a sensible zero/default
// in Defaults, and existing rows simply lack the key until next write.
type Settings struct {
	AIProvider string `json:"aiProvider"`
	AIModel    string `json:"aiModel"`
	AIEnabled  bool   `json:"aiEnabled"`
}

// Defaults returns the baseline used before anything is persisted and as the
// fallback when a stored value is missing or invalid.
func Defaults() Settings {
	return Settings{
		AIProvider: "openai",
		AIModel:    "gpt-4o-mini",
		AIEnabled:  false,
	}
}
