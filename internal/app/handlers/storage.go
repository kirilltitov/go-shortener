package handlers

type Storage interface {
	Get(int) (string, bool)
	Set(int, string)
}
