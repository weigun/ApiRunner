// hook.go
package business

// "net/http"

type hookFunc func(interface{}) interface{}

type hooks struct {
	beforeRequest hookFunc
	afterResponse hookFunc
}
