package gbwatch

import "time"

func NewTimer() *Timer {
	timer := &Timer{}
	timer.Reset()
	return timer
}

type Timer struct {
	initialTime time.Time
	stampedTime time.Time
}

func (timer *Timer) Reset() *Timer {
	timer.initialTime = time.Now()
	timer.stampedTime = timer.initialTime
	return timer
}

func (timer *Timer) Stamp() *Timer {
	timer.stampedTime = time.Now()
	return timer
}

func (timer *Timer) InitialTime() time.Time {
	return timer.initialTime
}

func (timer *Timer) StampedTime() time.Time {
	return timer.stampedTime
}

func (timer *Timer) ElapsedTime() time.Duration {
	return timer.stampedTime.Sub(timer.initialTime)
}

func (timer *Timer) ElapsedTimeS() int64 {
	return int64(timer.ElapsedTime().Seconds())
}

func (timer *Timer) ElapsedTimeMs() int64 {
	return timer.ElapsedTime().Milliseconds()
}

func (timer *Timer) ElapsedTimeNs() int64 {
	return timer.ElapsedTime().Nanoseconds()
}

func (timer *Timer) HasElapsedTime(value time.Duration) bool {
	return timer.ElapsedTime() >= value
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
