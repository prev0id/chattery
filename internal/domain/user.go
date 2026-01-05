package domain

type User struct {
	Username Username
	ImageID  ImageID
	Login    Login
	Password Password
}

type ImageID string
