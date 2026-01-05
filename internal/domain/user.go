package domain

type User struct {
	Username Username
	AvatarID ImageID
	Login    Login
	Password Password
}

type ImageID string

func (i ImageID) String() string { return string(i) }
