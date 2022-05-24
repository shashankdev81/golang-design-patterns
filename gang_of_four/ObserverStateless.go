package main

import (
	"fmt"
	"time"
)

type Observable interface {
	attach(c chan interface{})
	detach(c chan interface{})
	publish(val interface{})
}

type Observer interface {
	receive(chan interface{})
}

type Subject struct {
	//TODO makes this map thread safe
	privateChannels map[chan interface{}]bool
	done            chan interface{}
	publicChan      chan interface{}
}

type Monitor struct {
	privateChan chan interface{}
}

func NewSubject() *Subject {
	subject := &Subject{publicChan: make(chan interface{}), done: make(chan interface{}), privateChannels: make(map[chan interface{}]bool)}
	go func() {
		defer close(subject.publicChan)
		for {
			fmt.Println("Will read from public chan and publish to all private chans")
			select {
			case val := <-subject.publicChan:
				fmt.Println("Read from subjects chan", val)
				for privChan, _ := range subject.privateChannels {
					fmt.Println("Will publish to observers private chan")
					go func() {
						privChan <- val
					}()
				}
			case <-subject.done:
				fmt.Println("Cļosing all private channels")
				for c, _ := range subject.privateChannels {
					close(c)
				}
				fmt.Println("Subject closed so returning")
				return
			}
		}
	}()
	return subject
}

func NewMonitor() *Monitor {
	return &Monitor{privateChan: make(chan interface{})}
}

func (subject *Subject) attach(privateChan chan interface{}) {
	//create a new chan on which subject can publish events
	subject.privateChannels[privateChan] = true
	go func() {
		for {
			select {
			case val, ok := <-privateChan:
				if !ok {
					return
				}
				fmt.Println("Received val from subject", val, ok)
			}
		}
	}()

}

func (subject *Subject) detach(privateChan chan interface{}) {
	defer close(privateChan)
	delete(subject.privateChannels, privateChan)

}

func (subject *Subject) publish(val interface{}) {
	subject.publicChan <- val
}

func (monitor *Monitor) receive(receiveChan chan interface{}) {
	monitor.privateChan = receiveChan
}

func main() {
	type Stock struct {
		subject *Subject
	}
	type StockTicker struct {
		monitor *Monitor
	}

	appleStock := &Stock{subject: NewSubject()}
	defer close(appleStock.subject.done)

	moneyControl := &StockTicker{monitor: NewMonitor()}
	//client initiates attach by invoking attach api and passes observer that can be attached
	appleStock.subject.attach(moneyControl.monitor.privateChan)

	appleStock.subject.publish("Apple stock price is $230")
	time.Sleep(1 * time.Millisecond)

	economicTimes := &StockTicker{monitor: NewMonitor()}
	appleStock.subject.attach(economicTimes.monitor.privateChan)

	appleStock.subject.publish("Apple stock price is $240")
	time.Sleep(1 * time.Millisecond)

	appleStock.subject.detach(moneyControl.monitor.privateChan)
	time.Sleep(1 * time.Millisecond)

	appleStock.subject.publish("Apple stock price is $250")
	time.Sleep(1 * time.Millisecond)

	appleStock.subject.detach(economicTimes.monitor.privateChan)
	time.Sleep(1 * time.Millisecond)

	time.Sleep(1 * time.Second)

}
