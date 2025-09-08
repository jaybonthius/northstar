package db

import "embed"

//go:embed migrations/*
var MigrationFiles embed.FS
