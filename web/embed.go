package web

import (
	"embed"
	_ "embed"
)

//go:embed dist/assets/*
var Assets embed.FS

//go:embed dist/login.html
var LoginPage []byte

//go:embed dist/signup.html
var SignupPage []byte

//go:embed dist/app.html
var AppPage []byte
