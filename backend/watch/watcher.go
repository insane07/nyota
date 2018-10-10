package watch

import (
	"goprizm/log"
	"time"

	redis "github.com/go-redis/redis"
)

type Watcher struct {
	redis    redis.UniversalClient
	redisSub *redis.PubSub
}

func New() *Watcher {
	rc := redisClient()
	return &Watcher{
		redis:    rc,
		redisSub: rc.PSubscribe(),
	}
}

func (watcher *Watcher) Notify(channel string, data interface{}) {
	notify := func() error {
		return watcher.redis.Publish(channel, data).Err()
	}
	err := retryOp("notify", notify, 3, time.Second)
	if nil != err {
		log.Errorf("failed to publish: %v", err)
	}
}

func (watcher *Watcher) SubscribeAndReceive(channel []string, msgC chan *redis.Message) {
	for _, channelVal := range channel {
		watcher.redisSub = watcher.redis.Subscribe(channelVal)
	}

	for {
		msg, err := watcher.redisSub.ReceiveMessage()
		if err != nil {
			log.Errorf("watch - redis pub/sub recv err:%v", err)
			continue
		} else {
			msgC <- msg
		}
	}
}

func retryOp(opName string, op func() error, n int, after time.Duration) (err error) {
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
