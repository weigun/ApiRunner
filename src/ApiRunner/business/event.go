// event.go
package business

const (
	EVT_RUNNING = iota
	EVT_FINISHED
)

type handler func(params interface{}) error

type event struct {
	name      string
	asyncFlag bool
	timestamp int64
	status    int
	handler   handler
	params    interface{}
}
