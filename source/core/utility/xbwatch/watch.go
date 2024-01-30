package xbwatch

import "time"

func NewWatch() *Watch {
	watch := &Watch{}
	watch.Reset()
	return watch
}

type Watch struct {
	initialTime time.Time
	stampedTime time.Time
}

func (watch *Watch) Reset() *Watch {
	watch.initialTime = time.Now()
	watch.stampedTime = watch.initialTime
	return watch
}

func (watch *Watch) Stamp() *Watch {
	watch.stampedTime = time.Now()
	return watch
}

func (watch *Watch) InitialTime() time.Time {
	return watch.initialTime
}

func (watch *Watch) StampedTime() time.Time {
	return watch.stampedTime
}

func (watch *Watch) ElapsedTime() time.Duration {
	return watch.stampedTime.Sub(watch.initialTime)
}

func (watch *Watch) ElapsedTimeS() int64 {
	return int64(watch.ElapsedTime().Seconds())
}

func (watch *Watch) ElapsedTimeMs() int64 {
	return watch.ElapsedTime().Milliseconds()
}

func (watch *Watch) ElapsedTimeNs() int64 {
	return watch.ElapsedTime().Nanoseconds()
}

func (watch *Watch) HasElapsedTime(value time.Duration) bool {
	return watch.ElapsedTime() >= value
}

func NewCounter(initial, interval int) *Counter {
	counter := &Counter{initial: initial, interval: interval}
	counter.Reset()
	return counter
}

func NewDefaultCounter() *Counter {
	counter := NewCounter(0, 1)
	return counter
}

type Counter struct {
	count    int
	initial  int
	interval int
}

func (counter *Counter) Reset() *Counter {
	counter.count = counter.initial
	return counter
}

func (counter *Counter) Count() int {
	return counter.count
}

func (counter *Counter) Up() *Counter {
	counter.count += counter.interval
	return counter
}

func (counter *Counter) Down() *Counter {
	counter.count -= counter.interval
	return counter
}

func (counter *Counter) Plus(value int) *Counter {
	counter.count += value
	return counter
}

func (counter *Counter) Minus(value int) *Counter {
	counter.count -= value
	return counter
}

func (counter *Counter) HasCountedOver(value int) bool {
	return counter.count > value
}

func (counter *Counter) HasCountedBelow(value int) bool {
	return counter.count < value
}

func (counter *Counter) HasCountedUpTo(value int) bool {
	return counter.count >= value
}

func (counter *Counter) HasCountedDownTo(value int) bool {
	return counter.count <= value
}
