package assets

import "embed"

var (
	//go:embed dist/*
	Static embed.FS
)
