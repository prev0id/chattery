package domain

type User struct {
	ID       UserID
	Username Username
	AvatarID ImageID
	Login    Login
	Password Password
}

type ImageID string

func (i ImageID) String() string { return string(i) }

type Username string

func (u Username) String() string { return string(u) }
