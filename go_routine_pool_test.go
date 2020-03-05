package gorp_test

import (
	"github.com/developer-kikikaikai/go_routine_pool"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	WaitSec     = 3
	NumOfWorker = 4
)

/* Structure to test running thread */
type Result struct {
	sync.Mutex
	result int
}

//add result
func (r *Result) Add(value int) {
	r.Lock()
	r.result += value
	r.Unlock()
}
func (r *Result) Result() int {
	return r.result
}

type InputFunc struct {
	value  int
	result *Result
	wg     *sync.WaitGroup
}

func runningWorkerTestFunc(input interface{}) {
	in := input.(InputFunc)
	time.Sleep(WaitSec)
	in.result.Add(in.value)
	in.wg.Done()
}

func runningWorkerTestFuncClose(input interface{}) {
	in := input.(*int)
	(*in)++
}

func TestRunningWorker(t *testing.T) {
	var worker gorp.GoRoutineWorker
	//Run worker
	Convey("Check Worker is running normally", t, func() {
		_worker, err := gorp.RunWorker(NumOfWorker)
		So(err, ShouldBeNil)
		worker = _worker
		result := Result{}
		expect := 0
		wg := &sync.WaitGroup{}
		for i := 0; i < NumOfWorker; i++ {
			value := 1 << i
			//Get time before
			begin := time.Now()
			wg.Add(1)
			worker.Go(runningWorkerTestFunc, &InputFunc{value, &result, wg})
			//check this function is not blocked
			end := time.Now()
			expectTimeLess := begin.Add(time.Second)
			So(end.Before(expectTimeLess), ShouldBeTrue)

			//update expect value
			expect += value
		}
		wg.Wait()
		//Check result
		So(result.Result(), ShouldEqual, expect)

		//Check stop
		worker.Stop()
		value := 0
		worker.Go(runningWorkerTestFuncClose, &value)
		time.Sleep(time.Millisecond * 200)
		//check result
		So(value, ShouldEqual, 0)
	})

	Convey("Check running worker num", t, func() {
	})
}
