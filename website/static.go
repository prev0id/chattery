package static

import (
	"embed"
)

const (
	RootPath = "/"
	SrcPath  = "/src/"
)

//go:embed index.html
var IndexHTML []byte

//go:embed src/*
var Src embed.FS
