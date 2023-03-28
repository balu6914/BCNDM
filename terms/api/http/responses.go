package http

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}
