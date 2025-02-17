package client

type Server interface {
	InsertNickname(string) error
	UpdateNickname(string, string) error
	RemoveNickname(string) error

	InsertUsername(string) error
	RemoveUsername(string) error
}
