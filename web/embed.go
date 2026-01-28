package web

import (
	"embed"
	_ "embed"
)

//go:embed src/*
var Src embed.FS

//go:embed login.html
var LoginPage []byte

//go:embed signup.html
var SignupPage []byte

//go:embed app.html
var AppPage []byte
