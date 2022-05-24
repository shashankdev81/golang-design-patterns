package main
//
//import (
//	"fmt"
//	"time"
//)
//
//type Responsibility interface {
//	AddNext(chan interface{})
//	Ingress(request interface{})
//	Egress() chan interface{}
//	Close()
//}
//type Step struct {
//	name        string
//	done        chan interface{}
//	privateChan chan interface{}
//	ouputChan   chan interface{}
//	stepFunc    func(interface{}) interface{}
//}
//type Payload struct {
//	request interface{}
//	err     error
//}
//
//func (step *Step) AddNext(nextChan chan interface{}) {
//	var destChan chan interface{}
//	if nextChan == nil {
//		destChan = step.ouputChan
//	} else {
//		destChan = nextChan
//	}
//	go func() {
//		for {
//			select {
//			case val := <-step.privateChan:
//				fmt.Println("Reading from step ", step.name)
//				destChan <- step.execute(val)
//			case <-step.done:
//				close(step.privateChan)
//				return
//			}
//		}
//	}()
//}
//
//func (step *Step) execute(val interface{}) interface{} {
//	fmt.Println("Execute call on step  ", step.name, " ", val)
//	return step.stepFunc(val)
//
//}
//
////func (step *Step) execute(val interface{}) (interface{}, error) {
////	var err error
////	defer func() {
////		if r := recover(); r != nil {
////			fmt.Printf("Panic: %+v\n", r)
////			err = r.(error)
////		}
////	}()
////	fmt.Println("Execute call on step  ", step.name, " ", val)
////	output := step.stepFunc(val)
////	return output, err
////
////}
//
//func (step *Step) Ingress(req interface{}) {
//	step.privateChan <- req
//}
//func (step *Step) Egress() chan interface{} {
//	return step.ouputChan
//}
//func NewChainBuilder(steps []*Step) (*Step, *Step) {
//	var start, prev, step *Step
//	for _, step = range steps {
//		if prev == nil {
//			start = step
//			prev = step
//		} else {
//			prev.AddNext(step.privateChan)
//			prev = step
//		}
//	}
//	step.AddNext(nil)
//	return start, step
//}
//
//func NewStep(nameStr string, done chan interface{}, f func(req interface{}) interface{}) *Step {
//	return &Step{name: nameStr, done: done, privateChan: make(chan interface{}), ouputChan: make(chan interface{}), stepFunc: f}
//}
//
//func main() {
//	go doSomething(false)
//
//	time.Sleep(3 * time.Second)
//}
//
//func doSomething(isPanic bool) {
//	done := make(chan interface{})
//	step1 := NewStep("step1", done, func(req interface{}) interface{} {
//		if isPanic {
//			panic("Something went wrong")
//		}
//		return fmt.Sprintf("%v", req) + " in step 1, "
//	})
//	step2 := NewStep("step2", done, func(req interface{}) interface{} {
//		return fmt.Sprintf("%v", req) + " in step 2, "
//	})
//	step3 := NewStep("step3", done, func(req interface{}) interface{} {
//		return fmt.Sprintf("%v", req) + " in step 3, "
//	})
//
//	request := fmt.Sprintf("%v", "Test event")
//	start, end := NewChainBuilder([]*Step{step1, step2, step3})
//	start.Ingress(request)
//	stream := end.Egress()
//	go func() {
//		for {
//			select {
//			case v := <-stream:
//				fmt.Println("Received value from step ", end.name, " ", v)
//			}
//		}
//	}()
//	time.Sleep(1 * time.Second)
//	close(done)
//}
