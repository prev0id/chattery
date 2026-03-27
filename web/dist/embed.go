package web

import (
	"embed"
	_ "embed"
)

//go:embed assets/*
var Assets embed.FS

//go:embed login.html
var LoginPage []byte

//go:embed app.html
var AppPage []byte
