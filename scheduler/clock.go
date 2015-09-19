package scheduler

import (
	"time"
)

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type realClock struct{}

func (c realClock) Now() time.Time {
	return time.Now()
}

func (c realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

type fakeClock struct {
	time.Time
}

func (c *fakeClock) rewind(t time.Time) {
	c.Time = t
}

func (c *fakeClock) Now() time.Time {
	return c.Time
}

func (c *fakeClock) After(d time.Duration) <-chan time.Time {
	c.rewind(c.Time.Add(d))
	ch := make(chan time.Time, 1)
	ch <- c.Time
	return ch
}
