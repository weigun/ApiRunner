// hook.go
package business

import (
	"log"
	"net/http"
)

type hooks struct {
	beforeRequest func(*http.Request) *http.Request
	afterResponse func(*http.Response) *http.Response
}
