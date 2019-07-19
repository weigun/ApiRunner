// event.go
package business

import (
	"fmt"
)

const (
	EVT_STEP  = `STEP`
	EVT_STAGE = `STAGE`
)

type Event struct {
	Name string
	Type string
}

func (evt *Event) Fire() bool {
	topic := evt.Name
	eventBus.Publish(evt)
}

var eventBus EventBus.Bus

func init() {
	eventBus = EventBus.New()
}
