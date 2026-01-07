package response

import (
	"strconv"

	"github.com/Quak1/learn-http-go/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
