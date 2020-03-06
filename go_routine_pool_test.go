package gorp_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	gorp "github.com/developer-kikikaikai/go_routine_pool"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	WaitMsec    = 1000
	NumOfWorker = 4
	BufferSize  = NumOfWorker
)

/**************************
 * Structure to test running thread
 **************************/
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

type testRunningActor struct {
	value  int
	result *Result
	wg     *sync.WaitGroup
}

/**************************
 * Definition Action of routine
 **************************/
func (actor *testRunningActor) Do() {
	time.Sleep(WaitMsec * time.Millisecond)
	actor.result.Add(actor.value)
	actor.wg.Done()
}

func runningWorkerTestFunc(input interface{}) {
	in := input.(*testRunningActor)
	in.Do()
}

type testFuncCloseActor struct {
	in int
}

func (actor *testFuncCloseActor) Do() {
	actor.in++
}

func runningWorkerTestFuncClose(input interface{}) {
	in := input.(*testFuncCloseActor)
	in.Do()
}

/**************************
 * functions for test code
 **************************/
func getBeforeAfterGoAction(worker gorp.GoRoutineWorker, action gorp.Action) (time.Time, time.Time) {
	before := time.Now()
	worker.GoAction(action)
	//check this function is not blocked
	after := time.Now()
	return before, after
}

func testNonBlockingGoAction(t *testing.T, worker gorp.GoRoutineWorker, action gorp.Action) {
	before, after := getBeforeAfterGoAction(worker, action)
	expectTimeLess := before.Add(time.Second)
	fmt.Printf("endtime:%v, expectTimeLess:%v\n", after, expectTimeLess)
	So(after.Before(expectTimeLess), ShouldBeTrue)
}
func getBeforeAfterGo(worker gorp.GoRoutineWorker, routine func(interface{}), input interface{}) (time.Time, time.Time) {
	before := time.Now()
	worker.Go(routine, input)
	//check this function is not blocked
	after := time.Now()
	return before, after
}

func testNonBlockingGo(t *testing.T, worker gorp.GoRoutineWorker, routine func(interface{}), input interface{}) {
	before, after := getBeforeAfterGo(worker, routine, input)
	expectTimeLess := before.Add(time.Second)
	fmt.Printf("endtime:%v, expectTimeLess:%v\n", after, expectTimeLess)
	So(after.Before(expectTimeLess), ShouldBeTrue)
}

/**************************
 * Test code main
 **************************/
func TestRunningWorkerWithGoAction(t *testing.T) {
	var worker gorp.GoRoutineWorker
	//Run worker
	Convey("Check Worker is running with GoAction normally", t, func() {
		_worker := gorp.RunWorker(NumOfWorker, BufferSize)
		So(_worker, ShouldNotBeNil)
		worker = _worker
		result := Result{}
		expect := 0
		wg := &sync.WaitGroup{}
		begin := time.Now()
		for i := 0; i < NumOfWorker; i++ {
			value := 1 << i
			//Get time before
			wg.Add(1)
			testNonBlockingGoAction(t, worker, &testRunningActor{value, &result, wg})
			//update expect value
			expect += value
		}
		exit := time.Now()
		wg.Wait()
		//Check result
		So(result.Result(), ShouldEqual, expect)
		end := time.Now()
		//result : WaitSec < do all action < WaitSec + 1 sec
		expectTimeGreater := begin.Add(WaitMsec * time.Millisecond)
		fmt.Printf("endtime:%v, expectTimeGreater:%v\n", end, expectTimeGreater)
		So(end.After(expectTimeGreater), ShouldBeTrue)
		expectTimeLess := exit.Add((WaitMsec + 1000) * time.Millisecond)
		fmt.Printf("endtime:%v, expectTimeLess:%v\n", end, expectTimeLess)
		So(end.Before(expectTimeLess), ShouldBeTrue)

		//Check stop
		worker.Stop()
		actor := testFuncCloseActor{}
		value := actor.in
		worker.GoAction(&actor)
		//actor won't update actor.in param because there is no worker routine
		time.Sleep(time.Millisecond * 200)
		//check result
		So(actor.in, ShouldEqual, value)
	})
}

func TestRunningWorkerWithGo(t *testing.T) {
	var worker gorp.GoRoutineWorker
	//Run worker
	Convey("Check Worker is running with Go normally", t, func() {
		_worker := gorp.RunWorker(NumOfWorker, BufferSize)
		So(_worker, ShouldNotBeNil)
		worker = _worker
		result := Result{}
		expect := 0
		wg := &sync.WaitGroup{}
		begin := time.Now()
		for i := 0; i < NumOfWorker; i++ {
			value := 1 << i
			//Get time before
			wg.Add(1)
			testNonBlockingGo(t, worker, runningWorkerTestFunc, &testRunningActor{value, &result, wg})

			//update expect value
			expect += value
		}
		exit := time.Now()
		wg.Wait()
		//Check result
		So(result.Result(), ShouldEqual, expect)
		end := time.Now()
		//result : WaitSec < do all action < WaitSec + 1 sec
		expectTimeGreater := begin.Add(WaitMsec * time.Millisecond)
		fmt.Printf("endtime:%v, expectTimeGreater:%v\n", end, expectTimeGreater)
		So(end.After(expectTimeGreater), ShouldBeTrue)
		expectTimeLess := exit.Add((WaitMsec + 1000) * time.Millisecond)
		fmt.Printf("endtime:%v, expectTimeLess:%v\n", end, expectTimeLess)
		So(end.Before(expectTimeLess), ShouldBeTrue)

		//Check stop
		worker.Stop()
		actor := testFuncCloseActor{}
		value := actor.in
		worker.Go(runningWorkerTestFuncClose, &actor)
		//actor won't update actor.in param because there is no worker routine
		time.Sleep(time.Millisecond * 200)
		//check result
		So(actor.in, ShouldEqual, value)
	})
}

func TestWorkerLimit(t *testing.T) {
	Convey("Check running worker num", t, func() {
		worker := gorp.RunWorker(NumOfWorker, NumOfWorker+1)
		So(worker, ShouldNotBeNil)
		result := Result{}
		expect := 0
		wg := &sync.WaitGroup{}
		begin := time.Now()
		for i := 0; i < NumOfWorker+1; i++ {
			value := 1 << i
			//Get time before
			wg.Add(1)
			testNonBlockingGoAction(t, worker, &testRunningActor{value, &result, wg})
			//update expect value
			expect += value
		}

		exit := time.Now()
		wg.Wait()
		//Check result
		So(result.Result(), ShouldEqual, expect)
		end := time.Now()
		//result : WaitSec * 2 < do all action < WaitSec * 2 + 1 sec
		// Because 1 action was blocked
		expectTimeGreater := begin.Add(WaitMsec * 2 * time.Millisecond)
		fmt.Printf("endtime:%v, expectTimeGreater:%v\n", end, expectTimeGreater)
		So(end.After(expectTimeGreater), ShouldBeTrue)
		expectTimeLess := exit.Add(((WaitMsec * 2) + 1000) * time.Millisecond)
		fmt.Printf("endtime:%v, expectTimeLess:%v\n", end, expectTimeLess)
		So(end.Before(expectTimeLess), ShouldBeTrue)
	})

	Convey("Check running worker buffer size", t, func() {
		worker := gorp.RunWorker(NumOfWorker, 0)
		So(worker, ShouldNotBeNil)
		result := Result{}
		expect := 0
		wg := &sync.WaitGroup{}
		begin := time.Now()
		for i := 0; i < NumOfWorker; i++ {
			value := 1 << i
			//Get time before
			wg.Add(1)
			testNonBlockingGoAction(t, worker, &testRunningActor{value, &result, wg})
			//update expect value
			expect += value
		}
		{
			value := 1 << NumOfWorker
			wg.Add(1)
			_, after := getBeforeAfterGoAction(worker, &testRunningActor{value, &result, wg})
			expectTimeGreater := begin.Add(time.Millisecond * WaitMsec)
			fmt.Printf("endtime:%v, expectTimeGreater:%v\n", after, expectTimeGreater)
			So(after.After(expectTimeGreater), ShouldBeTrue)
			expect += value
		}

		exit := time.Now()
		wg.Wait()
		//Check result
		So(result.Result(), ShouldEqual, expect)
		end := time.Now()
		//result : WaitSec * 2 < do all action < WaitSec * 2 + 1 sec
		// Because 1 action was blocked
		expectTimeGreater := begin.Add(WaitMsec * 2 * time.Millisecond)
		fmt.Printf("endtime:%v, expectTimeGreater:%v\n", end, expectTimeGreater)
		So(end.After(expectTimeGreater), ShouldBeTrue)
		expectTimeLess := exit.Add(((WaitMsec * 2) + 1000) * time.Millisecond)
		fmt.Printf("endtime:%v, expectTimeLess:%v\n", end, expectTimeLess)
		So(end.Before(expectTimeLess), ShouldBeTrue)
	})
}
