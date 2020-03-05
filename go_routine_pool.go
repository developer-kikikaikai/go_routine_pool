// Copyright 2020 developer-kikikaikai.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gorp

import (
	"errors"
	"fmt"
)

type Action interface {
	Do()
}

type GoRoutineWorker interface {
	//TODO: can we use variable parameters?
	Go(actor Action)
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

func (worker *goRoutineWorker) Go(actor Action) {
	//send routine + input to thread
	fmt.Printf("Not implement yet\n")
	actor.Do()
}

func (worker *goRoutineWorker) Stop() {
	//stop worker
}
