package static

import (
	"embed"
)

const (
	RootPath     = "/"
	AuthPath     = "/sign"
	NotFoundPath = "/404"
	SettingsPath = "/settings"
	SrcPath      = "/src/"
)

//go:embed index.html
var IndexHTML []byte

//go:embed auth.html
var Auth []byte

//go:embed settings.html
var Settings []byte

//go:embed src/*
var Src embed.FS
