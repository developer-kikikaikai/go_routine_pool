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

func _say(s interface{}) {
	_s := s.(*string)
	say(*_s)
}

func ExampleGo() {
	numOfWorker := 4 //how many thread do you want to run?
	bufferSize := 4  //how many buffer size do you want to use?
	worker := grworker.RunWorker(numOfWorker, bufferSize)

	//instead of go say("world")
	msg := "world"
	worker.Go(_say, &msg)
	say("hello")
	worker.Stop()
	fmt.Println("vim-go")
}

func main() {
	ExampleGo()
}
