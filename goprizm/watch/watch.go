// TODO - If event(pub/sub msg) is send when subscriber is not connected, redis will drop those msgs.
// and subscriber will loose them. This need to be fixed.
//
// Redis blocking queue with a single element can be used for level triggered notification.
package watch

import (
	"encoding/json"
	"fmt"
	"goprizm/log"
	"time"

	redis "github.com/go-redis/redis"
)

type Event struct {
	Channel string
	Data    interface{}
}

type Watcher struct {
	// redis - client on which namespace:objects are stored and watched.
	redis redis.UniversalClient

	// redis subscribe handle.
	redisSub *redis.PubSub

	// Events is the chan to which notifications are send to watcher.
	Events chan Event
}

func New(r redis.UniversalClient) *Watcher {
	return &Watcher{
		redis:    r,
		redisSub: r.PSubscribe(),
		Events:   make(chan Event, 100),
	}
}

func (w *Watcher) Watch(channels ...string) {
	sub := func() error {
		return w.redisSub.Subscribe(channels...)
	}

	w.retryOp(fmt.Sprintf("subscribe channels:%v", channels), sub, -1, time.Second)
	log.Printf("watch - subscribe channel:%v done", channels)
}

func (w *Watcher) Notify(ev Event) error {
	var data string
	switch ev.Data.(type) {
	// If ev.Data is nil send empty string
	case nil:
		data = ""

	// If ev.Data is string, send without further modification.
	case string:
		data = ev.Data.(string)

	// If ev.Data is arbitary object, try to json serialize it.
	default:
		js, err := json.Marshal(ev.Data)
		if err != nil {
			return err
		}
		data = string(js)
	}

	notify := func() error {
		return w.redis.Publish(ev.Channel, data).Err()
	}
	return w.retryOp("notify", notify, 3, time.Second)
}

func (w *Watcher) Start() {
	run := func() {
		defer w.redisSub.Close()

		// Spawn goroutine to handle msgs. Separate goroutine/chan is used for buffering
		// since redis pubsub does not enqueue msgs
		msgC := make(chan *redis.Message, 5000)
		go func() {
			for msg := range msgC {
				w.Events <- Event{msg.Channel, msg.Payload}
			}
		}()

		for {
			msg, err := w.redisSub.ReceiveMessage()
			if err != nil {
				log.Errorf("watch - redis pub/sub recv err:%v", err)
				continue
			}
			msgC <- msg
		}
	}
	go run()
}

// retryOp retries given operation atmost n times if it returns error. Retry is
// performed after given duration. If n is -1, retry is performed till op succeeds.
func (w *Watcher) retryOp(opName string, op func() error, n int, after time.Duration) (err error) {
	i := 0
	for {
		// If n is +ve and num of retries reached n.
		if n > 0 && i == n {
			return err
		}

		// Execute op and on error sleep for given duration.
		if err = op(); err != nil {
			i += 1
			log.Errorf("watch - %s (retry:%d) failed err:%v", opName, i, err)
			time.Sleep(after)
			continue
		}
		return nil
	}
}
