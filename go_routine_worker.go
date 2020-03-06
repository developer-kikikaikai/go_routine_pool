package grworker

import (
	"sync"
)

// Action is an interface instead of go func. Please write function which you want to use go func in Do() method. You can set input parameter in this instance
type Action interface {
	Do()
}

// GoRoutineWorker is an interface instead of go routine.
type GoRoutineWorker interface {
	GoAction(actor Action)
	//TODO: can we use variable parameters?
	Go(routine func(interface{}), input interface{})
	Stop()
}

// RunWorker is a function to create GoRoutineWorker. Please input a number of workers and buffer size of channel
func RunWorker(numOfWorker int, bufferSize int) GoRoutineWorker {
	worker := newWorker(bufferSize)
	worker.run(numOfWorker)
	return worker
}

/**************************************************************
 * internal definition
 **************************************************************/
// transportEvent is an interface to transport action event of GoAction/Go. Worker receive this interface from channel, and call do()
type transportEvent interface {
	// do is a function called by worker
	do()
}

// transportAction is an message to transport requested event of GoAction
type transportAction struct {
	actor Action
}

func (t *transportAction) do() {
	t.actor.Do()
}

// transportFunction is an message to transport requested event of Go
type transportFunction struct {
	routine func(interface{})
	input   interface{}
}

func (t *transportFunction) do() {
	t.routine(t.input)
}

// transportClose is an message to transport close event
type transportClose struct {
	ch chan transportEvent
}

func (t *transportClose) do() {
	close(t.ch)
}

// goRoutineWorker is an instance of GoRoutineWorker
type goRoutineWorker struct {
	ch   chan transportEvent  // ch is a channel to communicate worker thread
	wg   *sync.WaitGroup      // wg is a work group to stop worker
	send func(transportEvent) // send is a sender function with message transportEvent
	wait func()               //wait is a function to wait closing workers
}

// newWorker is a function to create goRoutineWorker instance
func newWorker(bufferSize int) *goRoutineWorker {
	worker := goRoutineWorker{}
	worker.ch = make(chan transportEvent, bufferSize)
	worker.wg = &sync.WaitGroup{}
	worker.setActualFunc()
	return &worker
}

func (worker *goRoutineWorker) run(numOfWorker int) {
	//run thread
	for i := 0; i < numOfWorker; i++ {
		go workerMain(worker)
	}
}

func workerMain(worker *goRoutineWorker) {
	//count work group to manage stopping
	worker.wg.Add(1)
	//run routine until stop
	worker.runWorker()
	worker.wg.Done()
}

func (worker *goRoutineWorker) runWorker() {
	for event := range worker.ch {
		event.do()
	}
}

// setActualFunc is a function to set actual send/wait function
func (worker *goRoutineWorker) setActualFunc() {
	worker.send = worker.sendActual
	worker.wait = worker.waitActual
}
func (worker *goRoutineWorker) sendActual(event transportEvent) {
	worker.ch <- event
}
func (worker *goRoutineWorker) waitActual() {
	worker.wg.Wait()
}

// setNoFunc is a function to set no-action function
func (worker *goRoutineWorker) setNoFunc() {
	worker.send = worker.notSend
	worker.wait = worker.notWait
}
func (worker *goRoutineWorker) notSend(event transportEvent) {
}
func (worker *goRoutineWorker) notWait() {
}

// implement interface
func (worker *goRoutineWorker) GoAction(actor Action) {
	event := &transportAction{
		actor: actor,
	}
	worker.send(event)
}

func (worker *goRoutineWorker) Go(routine func(interface{}), input interface{}) {
	event := &transportFunction{
		routine: routine,
		input:   input,
	}
	worker.send(event)
}

func (worker *goRoutineWorker) Stop() {
	worker.send(&transportClose{worker.ch})
	worker.wait()
	//after all, set function not to do anything
	worker.setNoFunc()
}
