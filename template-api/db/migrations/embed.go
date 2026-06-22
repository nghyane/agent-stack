// Package migrations embeds SQL migration files for goose.
package migrations

import "embed"

// FS holds the SQL migration files, embedded for goose.
//
//go:embed *.sql
var FS embed.FS
