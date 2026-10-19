// Package migrations встраивает SQL-миграции в бинарник,
// чтобы приложение могло применить их само при старте.
package migrations

import "embed"

// FS содержит все файлы миграций в формате golang-migrate.
//
//go:embed *.sql
var FS embed.FS
