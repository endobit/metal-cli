// Package sql embeds SQL migrations.
package sql

import "embed"

//go:embed migrations/*.sql
var Migrations embed.FS
