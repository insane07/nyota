package watch

import (
	"fmt"
	"goprizm/services"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	w := New(services.Redis())
	w.Start()

	n := 5
	notify := func(channel string) {
		time.Sleep(500 * time.Millisecond)
		for i := 0; i < n; i++ {
			data := fmt.Sprintf("%s:data:%d", channel, i)
			if err := w.Notify(Event{channel, data}); err != nil {
				t.Fatalf("Failed to notify err:%v", err)
			}
		}
	}

	verify := func(channels ...string) {
		i := 0
		for {
			select {
			case <-w.Events:
				if i++; i == n*len(channels) {
					return
				}
			case <-time.After(5 * time.Second):
				t.Fatalf("Watch timeout")
			}
		}

	}

	var channels []string
	for i := 0; i < 10; i++ {
		channels = append(channels, fmt.Sprintf("test.chan:%d", i))
	}
	for _, channel := range channels {
		c := channel

		// Spawn each watch/notify in separate goroutine to check race conditions.
		go func() {
			w.Watch(c)
			notify(c)
		}()
	}

	verify(channels...)
}
