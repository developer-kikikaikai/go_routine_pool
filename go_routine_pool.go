// Copyright 2020 developer-kikikaikai.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gorp

import (
	"fmt"
	"sync"
)

//Please create struct which has input parameter of go routine
type Action interface {
	Do()
}

type GoRoutineWorker interface {
	GoAction(actor Action)
	//TODO: can we use variable parameters?
	Go(routine func(interface{}), input interface{})
	Stop()
}

func RunWorker(numOfWorker int) GoRoutineWorker {
	worker := newWorker()
	worker.run(numOfWorker)
	return worker
}

type transportEvent interface {
	do() (doesContinue bool)
}

type transportAction struct {
	action Action
}

func (t *transportAction) do() (doesContinue bool) {
	t.action.Do()
	doesContinue = true
	return
}

type transportFunction struct {
	routine func(interface{})
	input   interface{}
}

func (t *transportFunction) do() (doesContinue bool) {
	t.routine(t.input)
	doesContinue = true
	return
}

type transportClose struct {
}

func (t *transportClose) do() (doesContinue bool) {
	doesContinue = false
	return
}

type goRoutineWorker struct {
	ch chan transportEvent
	wg *sync.WaitGroup //work group to stop worker
}

func newWorker() *goRoutineWorker {
	worker := goRoutineWorker{}
	worker.ch = make(chan transportEvent, 0)
	worker.wg = &sync.WaitGroup{}
	return &worker
}

func (worker *goRoutineWorker) run(numOfWorker int) {
	//run thread
	for i := 0; i < numOfWorker; i++ {
		go func(worker *goRoutineWorker) {
			//count work group to manage stopping
			worker.wg.Add(1)
			//run routine until stop
			worker.runWorker()
			//worker is done!
			worker.wg.Done()
		}(worker)
	}
}

func (worker *goRoutineWorker) runWorker() {
	for event := range worker.ch {
		//do event action
		doesContinue := event.do()

		//if we don't continue, close channel
		if !doesContinue {
			//continue loop?
			close(worker.ch)
		}
	}
}

func (worker *goRoutineWorker) GoAction(actor Action) {
	//send routine + input to thread
	fmt.Printf("Not implement yet\n")
	actor.Do()
}

func (worker *goRoutineWorker) Go(routine func(interface{}), input interface{}) {
	fmt.Printf("Not implement yet\n")
	routine(input)
}

func (worker *goRoutineWorker) Stop() {
	//stop worker
}
