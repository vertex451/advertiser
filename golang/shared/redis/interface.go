package redis

type Interface interface {
	Push(key, val string) error
	Pop(key string) (string, error)
}
