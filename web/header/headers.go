package header

import (
	"mime"
	"net/http"
)

type Header http.Header

var (
	ContentTypeJson = mime.FormatMediaType("application/json", map[string]string{
		"charset": "utf-8",
	})
	ContentTypeForm = mime.FormatMediaType("application/x-www-form-urlencoded", nil)
)

const (
	ContentType        = "content-type"
	ContentLength      = "content-length"
	ContentDisposition = "content-disposition"
)
