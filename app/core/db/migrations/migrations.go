package migrations

import (
	_ "embed"
)

//go:embed sql/schema.sql
var Schema string
