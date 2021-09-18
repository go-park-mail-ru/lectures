package main

import (
	"fmt"
	"sync"
	"time"
)

const idle = 1

type Connection struct {
	sync.Mutex

	timer *time.Timer
	count uint
}

func (c *Connection) Update() uint {
	c.Lock()
	defer c.Unlock()

	if c.timer != nil && !c.timer.Stop() {
		c.timer.Reset(idle * time.Second)
		return c.count
	}

	// if connection don't use more than idle timeout than release
	c.timer = time.AfterFunc(idle*time.Second, func() {
		c.Lock()
		defer c.Unlock()

		c.count += 1

		c.timer.Stop()
	})
	return c.count
}

func main() {
	c := Connection{}
	update := func() {
		c.Update()
		fmt.Printf("conn count -> %d\n", c.count)
	}

	// Opens new connection
	update()

	// Takes old connection. Do not increment count
	update()

	// Uncomment this line to imitate idle timeout and use new connection
	// time.Sleep((idle + 1) * time.Second)

	update()
}
