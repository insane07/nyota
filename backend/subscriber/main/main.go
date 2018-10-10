package main

import (
	"nyota/backend/model"
	str "nyota/backend/store"
	"nyota/backend/watch"
	"encoding/json"
	"goprizm/log"

	redis "github.com/go-redis/redis"
)

func main() {
	msgC := make(chan *redis.Message, 5000)
	go func() {
		log.Debugf("waiting for msg")
		for msg := range msgC {
			log.Debugf("channel:" + msg.Channel)
			var event model.Event
			unmarshallErr := json.Unmarshal([]byte(msg.Payload), &event)
			if unmarshallErr != nil {
				log.Errorf("Unmarshal error: %v", unmarshallErr)
			} else {
				log.Debugf("payload Event UUID:%s", event.UUID)
				log.Debugf("payload Event Entity Name:%s", event.Data.EntityName)
			}
		}
	}()
	subStore := &str.Store{Watcher: watch.New()}
	// defer subStore.watcher.redisSub.Close()
	// defer subStore.watcher.redis.Close()
	subStore.Watcher.SubscribeAndReceive([]string{"event"}, msgC)
}
