package domain

type User struct {
	Username Username
	AvatarID ImageID
	Login    Email
	Password Password
	ID       UserID
}

type ImageID string

func (i ImageID) String() string { return string(i) }

type Username string

func (u Username) String() string { return string(u) }
