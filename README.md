# go_routine_pool (gorp)
This is a package to save number of go routine.

Go routine is a wonderful functions to run featres parallel.  
But sometimes our go code reaches the limit of go routines number. In this case, I want go routines to run such as threadpool.  
I can't find those libraries, so I've tried to create it.

# Installation

Please import this package:

```
import "github.com/developer-kikikaikai/go_routine_pool"
```

If you don't use go modules (go mod), please get source code:

```
go get github.com/developer-kikikaikai/go_routine_pool
```

# Usage

I explain them by using base code in [A Tour of Go: Goroutines](https://tour.golang.org/concurrency/1):

```
package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func main() {
	go say("world")
	say("hello")
}
```

## 1. Initialize worker

At first, please run worker with 1. number of workers and 2. buffer size.

```
	numOfWorker := 4 //how many thread do you want to run?
	bufferSize := 4  //how many buffer size do you want to use?
	worker := grworker.RunWorker(numOfWorker, bufferSize)
```


## 2-1. Run function by using Go

Please implement `func(interface{})` instead of your function which want to run by go func.

```
func _say(s interface{}) {
	_s := s.(*string)
	say(*_s)
}
```

And call `worker.Go` with the function and a parameter (pointer)

```
	//instead of go say("world")
	msg := "world"
	worker.Go(_say, &msg)
```

Please see [example/example_go.go](https://github.com/developer-kikikaikai/grworker/blob/master/example/example_go.go)

This usage is nearly like go func. But we need to cast interface in function.

## 2-2. Run function by using GoAction

Please implement struct related to the interface

```
type Action interface {
	Do()
}
```

such as:

```
type action struct {
	s string
}

func (a *action) Do() {
	say(a.s)
}
```

* Instead of func/input param, it uses Action.Do()/members of struct.

And call `worker.GoAction` with the instance

```
	//instead of go say("world")
	worker.GoAction(&action{"world"})
```

Please see [example/example_go_action.go](https://github.com/developer-kikikaikai/grworker/blob/master/example/example_go_action.go)

## 3. Stop workers

Finally, please call `worker.Stop()` to stop workers.

```
	worker.Stop()
```
