package cache

import "time"

type janitor struct {
	interval time.Duration
	stop     chan bool
}

func newJanitor(interval time.Duration) *janitor {
	return &janitor{
		interval: interval,
		stop:     make(chan bool),
	}
}

func (j *janitor) Start(c *Cache) {
	ticker := time.NewTicker(j.interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func (j *janitor) Run(c *Cache) {
	go j.Start(c)
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}
