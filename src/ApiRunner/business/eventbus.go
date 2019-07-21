package business

/*
每个testrunner都会携带一个eventbus实例
主要是用于用例各种事件的回调管理
*/

import (
	"fmt"
	"log"
	"sync"
)

const (
	MAX_TASKS = 8 //单个testrunner最大可同时执行的回调数量
)

type eventbus struct {
	queue   chan *event //事件队列
	running chan *event //维护运行中的事件
	wg      sync.WaitGroup
}

func NewEvevntBus() *eventbus {
	bus := &eventbus{
		queue:   make(chan *event, 2),
		running: make(chan *event, MAX_TASKS),
	}
	go bus.listen()
	return bus
}

func (bus *eventbus) Fire(evt event) {
	if !evt.asyncFlag {
		err := evt.handler(evt.params)
		if err != nil {
			log.Println(`event `, evt.name, ` callback error:`, err.Error())
		}
		return
	}
	bus.queue <- &evt
	log.Printf(`%s add to queue`, evt.name)
}

func (bus *eventbus) listen() {
	for evt := range bus.queue {
		if len(bus.running) >= MAX_TASKS {
			//如果已经达到上限了,先等待所有回调完成
			log.Println(`reach max tasks,join`)
			bus.join()
		}
		bus.running <- evt //将即将要执行的回调放到运行队列中
		bus.wg.Add(1)
		go func() {
			err := evt.handler(evt.params)
			if err != nil {
				log.Println(`event `, evt.name, ` async callback error:`, err.Error())
			}
			evt.status = EVT_FINISHED
			bus.wg.Done()
		}()

	}
}

func (bus *eventbus) join() {
	bus.wg.Wait()
	for i := 0; i < MAX_TASKS; i++ {
		e := <-bus.running
		if e.status != EVT_FINISHED {
			log.Println(`omg!!event status not finished,`, e.name)
		}
	}
}

func (bus *eventbus) Shutdown() {
	close(bus.queue)
	close(bus.running)
}
