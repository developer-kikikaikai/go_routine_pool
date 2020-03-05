// Copyright 2020 developer-kikikaikai.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gorp

import "errors"

type GoRoutineWorker interface {
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

func (worker *goRoutineWorker) Go(routine func(interface{}), input interface{}) {
	//send routine + input to thread
	routine(input)
}

func (worker *goRoutineWorker) Stop() {
	//stop worker
}
