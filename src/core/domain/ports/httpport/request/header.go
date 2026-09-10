package request

type Header interface {
	Get(key string) string
	Set(key string, value string)
	Add(key, value string)
}
