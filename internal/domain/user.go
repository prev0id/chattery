package domain

type User struct {
	Username Username
	ImageID  ImageID
	Login    string
	Password Password
}

type Username string

func (u Username) String() string { return string(u) }

const UserUnknown Username = ""

type ImageID string

type Login string
