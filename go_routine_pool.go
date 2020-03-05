// Copyright 2020 developer-kikikaikai.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gorp

import (
	"errors"
	"fmt"
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

func RunWorker(numOfWorker int) (GoRoutineWorker, error) {
	worker := newWorker()
	if err := worker.run(numOfWorker); err != nil {
		return nil, err
	}
	return worker, nil
}

type goRoutineWorker struct {
}

func newWorker() *goRoutineWorker {
	//run thread
	return &goRoutineWorker{}
}

func (worker *goRoutineWorker) run(numOfWorker int) error {
	return errors.New("Not implement yet")
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
