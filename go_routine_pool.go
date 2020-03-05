// Copyright 2020 developer-kikikaikai.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gorp

import (
	"fmt"
	"sync"
	"time"
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
	actor Action
}

func (t *transportAction) do() (doesContinue bool) {
	t.actor.Do()
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

type workerAction interface {
	send(event transportEvent)
	wait()
}

type workerNormalAction struct {
	ch chan transportEvent
	wg *sync.WaitGroup //work group to stop worker
}

func (w *workerNormalAction) send(event transportEvent) {
	w.ch <- event
}
func (w *workerNormalAction) wait() {
	w.wg.Wait()
}

type workerNoAction struct {
}

func (w *workerNoAction) send(event transportEvent) {
}
func (w *workerNoAction) wait() {
}

type goRoutineWorker struct {
	ch     chan transportEvent
	wg     *sync.WaitGroup //work group to stop worker
	action workerAction
}

func newWorker() *goRoutineWorker {
	worker := goRoutineWorker{}
	worker.ch = make(chan transportEvent)
	worker.wg = &sync.WaitGroup{}
	worker.action = &workerNormalAction{worker.ch, worker.wg}
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
	//send transportAction event to worker
	fmt.Printf("begin:%v\n", time.Now())
	event := &transportAction{
		actor: actor,
	}
	worker.action.send(event)
	fmt.Printf("end:%v\n", time.Now())
}

func (worker *goRoutineWorker) Go(routine func(interface{}), input interface{}) {
	//send transportAction event to worker
	event := &transportFunction{
		routine: routine,
		input:   input,
	}
	worker.action.send(event)
}

func (worker *goRoutineWorker) Stop() {
	//send transportClose event to worker
	worker.action.send(&transportClose{})
	worker.action.wait()
	worker.action = &workerNoAction{}
}
