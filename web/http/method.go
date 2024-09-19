package http

type Method string

func (m Method) String() string {
	return string(m)
}

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	PATCH   Method = "PATCH"
	HEAD    Method = "HEAD"
	OPTIONS Method = "OPTIONS"
	TRACE   Method = "TRACE"
	CONNECT Method = "CONNECT"
)
