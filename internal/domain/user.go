package domain

type Username string

func (u Username) String() string { return string(u) }

const UnknownUsername Username = ""

type Password string

type User struct {
	Username Username
}

type Friend struct {
	Username Username
	ChatID   ChatID
}
