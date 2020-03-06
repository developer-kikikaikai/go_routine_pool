package main

import (
	"fmt"
	"github.com/developer-kikikaikai/grworker"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

//Define structure, and implement Do
type action struct {
	s string
}

func (a *action) Do() {
	say(a.s)
}

func ExampleGoAction() {

	numOfWorker := 4 //how many thread do you want to run?
	bufferSize := 4  //how many buffer size do you want to use?
	worker := grworker.RunWorker(numOfWorker, bufferSize)

	//instead of go say("world")
	worker.GoAction(&action{"world"})
	say("hello")
	worker.Stop()
	fmt.Println("vim-go")
}

func main() {
	ExampleGoAction()
}
