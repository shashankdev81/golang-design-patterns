package main

import (
	"fmt"
	"time"
)

type Responsibility interface {
	AddNext(nextChan chan Payload)
	Ingress(payload Payload)
	Egress() chan Payload
}
type Step struct {
	name        string
	done        chan interface{}
	privateChan chan Payload
	ouputChan   chan Payload
	stepFunc    func(interface{}) interface{}
}
type Payload struct {
	request interface{}
	err     error
}

func (step *Step) AddNext(nextChan chan Payload) {
	var destChan chan Payload
	if nextChan == nil {
		destChan = step.ouputChan
	} else {
		destChan = nextChan
	}
	go func() {
		for {
			select {
			case payload := <-step.privateChan:
				fmt.Println("Reading from step ", step.name)
				if payload.err == nil {
					payload.request, payload.err = step.execute(payload.request)
				}
				destChan <- payload
			case <-step.done:
				close(step.privateChan)
				return
			}
		}
	}()
}

func (step *Step) execute(val interface{}) (interface{}, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %+v\n", r)
			err = fmt.Errorf("%v", r)

		}
	}()
	fmt.Println("Execute call on step  ", step.name, " ", val)
	output := step.stepFunc(val)
	return output, err

}

func (step *Step) Ingress(payload Payload) {
	step.privateChan <- payload
}

func (step *Step) Egress() chan Payload {
	return step.ouputChan
}

func NewChainBuilder(steps []*Step) (*Step, *Step) {
	var start, prev, step *Step
	for _, step = range steps {
		if prev == nil {
			start = step
			prev = step
		} else {
			prev.AddNext(step.privateChan)
			prev = step
		}
	}
	step.AddNext(nil)
	return start, step
}

func NewStep(nameStr string, done chan interface{}, f func(req interface{}) interface{}) *Step {
	return &Step{name: nameStr, done: done, privateChan: make(chan Payload), ouputChan: make(chan Payload), stepFunc: f}
}

func main() {
	go doSomething(true)
	time.Sleep(3 * time.Second)
}

func doSomething(isPanic bool) {
	done := make(chan interface{})
	step1 := NewStep("step1", done, func(req interface{}) interface{} {
		return fmt.Sprintf("%v", req) + " in step 1, "
	})
	step2 := NewStep("step2", done, func(req interface{}) interface{} {
		if isPanic {
			panic("Something went wrong")
		}

		return fmt.Sprintf("%v", req) + " in step 2, "
	})
	step3 := NewStep("step3", done, func(req interface{}) interface{} {
		return fmt.Sprintf("%v", req) + " in step 3, "
	})

	request := Payload{request: fmt.Sprintf("%v", "Test event")}
	start, end := NewChainBuilder([]*Step{step1, step2, step3})
	start.Ingress(request)
	stream := end.Egress()
	go func() {
		for {
			select {
			case v := <-stream:
				fmt.Println("Received value from step ", end.name, " ", v)
			}
		}
	}()
	time.Sleep(1 * time.Second)
	close(done)
}
