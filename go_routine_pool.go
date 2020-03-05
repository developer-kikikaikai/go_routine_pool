// Copyright 2020 developer-kikikaikai.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gorp

import "fmt"

type GoRoutineWorker interface {
	Start(numOfWorker int)
	//TODO: can we use variable parameters?
	Go(routine func(interface{}), input interface{})
	Stop()
	Delete()
}

func NewWorker() GoRoutineWorker {
	return newWorker()
}

type goRoutineWorker struct {
}

func newWorker() GoRoutineWorker {
	return &goRoutineWorker{}
}

func (worker *goRoutineWorker) Start() {
	//run thread
}

func (worker *goRoutineWorker) Go(routine func(interface{}), input interface{}) {
	//send routine + input to thread
}

func (worker *goRoutineWorker) Stop() {
	//stop worker
}

func (worker *goRoutineWorker) Delete() {
	//close resoruces
}
