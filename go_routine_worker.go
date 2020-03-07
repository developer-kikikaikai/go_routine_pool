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
type messageID int

const (
	msgAction messageID = iota
	msgFunction
	msgClose
)

type transportMessage struct {
	msgType messageID

	//For action field
	actor Action

	//For function field
	routine func(interface{})
	input   interface{}
}

func doMsgAction(actor Action) {
	actor.Do()
}

func doMsgFunction(routine func(interface{}), input interface{}) {
	routine(input)
}

// goRoutineWorker is an instance of GoRoutineWorker
type goRoutineWorker struct {
	ch chan transportMessage // ch is a channel to communicate worker thread
	wg *sync.WaitGroup       // wg is a work group to stop worker
}

// newWorker is a function to create goRoutineWorker instance
func newWorker(bufferSize int) *goRoutineWorker {
	worker := goRoutineWorker{}
	worker.ch = make(chan transportMessage, bufferSize)
	worker.wg = &sync.WaitGroup{}
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
		switch event.msgType {
		case msgAction:
			doMsgAction(event.actor)
		case msgFunction:
			doMsgFunction(event.routine, event.input)
		case msgClose:
			close(worker.ch)
		default:
			//unknown
		}
	}
}

func (worker *goRoutineWorker) send(event transportMessage) {
	if worker.ch != nil {
		worker.ch <- event
	}
}
func (worker *goRoutineWorker) wait() {
	if worker.wg != nil {
		worker.wg.Wait()
	}
}

// implement interface
func (worker *goRoutineWorker) GoAction(actor Action) {
	event := transportMessage{
		msgType: msgAction,
		actor:   actor,
	}
	worker.send(event)
}

func (worker *goRoutineWorker) Go(routine func(interface{}), input interface{}) {
	event := transportMessage{
		msgType: msgFunction,
		routine: routine,
		input:   input,
	}
	worker.send(event)
}

func (worker *goRoutineWorker) Stop() {
	event := transportMessage{
		msgType: msgClose,
	}
	worker.send(event)
	worker.wait()
	worker.ch = nil
	worker.wg = nil
}
