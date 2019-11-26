// testcase_render.go
package business

import (
	"ApiRunner/business/template"
	// "ApiRunner/models"

	"bytes"
	// "fmt"
	// "log"
	// "strings"
	"sync"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type renderer struct {
	sync.Mutex
	buf *bytes.Buffer
}

func newRenderer() *renderer {
	return &renderer{buf: bytes.NewBufferString(``)}
}

func (r *renderer) fillData(src string, data interface{}) []byte {
	r.Lock()
	defer r.Unlock()
	t := template.New().Funcs(funcMap)
	t.Parse(src)
	if data == nil {
		data = make(map[string]interface{})
	}
	t.Execute(r.buf, data)
	val := r.buf.Bytes()
	r.buf.Reset()
	return val
}
