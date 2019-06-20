package resolve

type Redis interface {
	LoadRedis(host, pass string, port float64) bool
}
