package collections
import (
"sync"
"time"
)

type TimeoutCond struct {
	L      sync.Locker
	signal chan int
}

func NewTimeoutCond(l sync.Locker) *TimeoutCond {
	cond  := TimeoutCond{L:l, signal:make(chan int,0)}
	return &cond
}

/**
return remain wait time, and is interrupt
 */
func (this *TimeoutCond) WaitWithTimeout(timeout time.Duration) (time.Duration,bool) {
	//wait should unlock mutex,  if not will cause deadlock
	this.L.Unlock()
	defer this.L.Lock()
	begin := time.Now().Nanosecond()
	select {
	case  _,ok := <- this.signal:
		end := time.Now().Nanosecond()
		return time.Duration(end - begin), !ok
	case  <- time.After(timeout):
		return 0,false
	}
}

/**
return is interrupt
 */
func (this *TimeoutCond) Wait() (bool) {
	this.L.Unlock()
	defer this.L.Lock()
	_, ok := <- this.signal
	return !ok
}

func (this *TimeoutCond) Signal()  {
	select {
	case this.signal <- 1:
	default:
	}
}

func (this *TimeoutCond) Interrupt(){
	close(this.signal)
	this.signal = make(chan int, 0)
}