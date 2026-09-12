package nebulastub

import (
	"embed"
	"io/fs"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func Asstest() fs.FS {
	distFS, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		panic(err)
	}

	return distFS
}
