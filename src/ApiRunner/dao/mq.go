package dao

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

var (
	Prefix = "ApiRunnerMQ"
)

type MQ interface {
	Subscribe(channel string) //error //*redis.PubSub
	Publish(channel string, msg interface{}) error
	Fetch() (string, interface{})
	Close()
}

type RedisMQ struct {
	core   *redisCache
	pubSub *redis.PubSub
}

func (this *RedisMQ) Subscribe(channel string) {
	channel = fmt.Sprintf(`%s:%s`, Prefix, channel)
	this.pubSub = this.core.Client.Subscribe(channel)
}

func (this *RedisMQ) Publish(channel string, msg interface{}) error {
	channel = fmt.Sprintf(`%s:%s`, Prefix, channel)
	if err := this.core.Client.Publish(channel, msg).Err(); err != nil {
		return errors.New(fmt.Sprintf(`MQ Publish error,reason %s`, err.Error()))
	}
	return nil
}

func (this *RedisMQ) Fetch() (string, interface{}) {
	channel := this.pubSub.Channel()
	msg := <-channel
	fmt.Println(`get message:`, msg.String())
	return msg.Channel, msg.Payload
}

func (this *RedisMQ) Close() {
	this.pubSub.Close()
}

func NewMQ(mqType string) MQ {
	var mqPtr MQ
	switch mqType {
	case `redis`:
		mqPtr = &RedisMQ{core: GetCache()}
	default:
		mqPtr = &RedisMQ{core: GetCache()}
	}
	return mqPtr
}
