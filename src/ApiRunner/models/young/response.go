// response.go
package young

import (
	"net/http"
)

type Response struct {
	Code    int
	Header  http.Header
	Cookies []*http.Cookie
	Content string
	ErrMsg  string
}
