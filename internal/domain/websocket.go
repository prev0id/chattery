package domain

type Connection interface {
	ReadPump()
	WritePump()
}
