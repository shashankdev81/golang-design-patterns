package main
//
//import (
//	"fmt"
//	"time"
//)
//
//type Observable interface {
//	attach(o Observer)
//	detach(o Observer)
//	publish(val interface{})
//}
//
//type Observer interface {
//	receive(chan interface{}, chan interface{})
//}
//
//type Subject struct {
//	observerToChanMap map[Observer]chan interface{}
//	done              chan interface{}
//	publicChan        chan interface{}
//}
//
//type Monitor struct {
//	privateChan chan interface{}
//}
//
//func NewSubject() *Subject {
//	subject := &Subject{publicChan: make(chan interface{}), done: make(chan interface{}), observerToChanMap: make(map[Observer]chan interface{})}
//	go func() {
//		defer close(subject.publicChan)
//		for {
//			fmt.Println("Will read from public chan and publish to all private chans")
//			select {
//			case val := <-subject.publicChan:
//				fmt.Println("Read from subjects chan", val)
//				for _, obsChan := range subject.observerToChanMap {
//					fmt.Println("Will publish to observers private chan")
//					go func() {
//						obsChan <- val
//					}()
//				}
//			case <-subject.done:
//				fmt.Println("Subject closed so returning")
//				return
//			}
//		}
//	}()
//	return subject
//}
//
//func NewMonitor() *Monitor {
//	return &Monitor{privateChan: make(chan interface{})}
//}
//
//func (subject *Subject) attach(observer Observer) {
//	//create a new chan on which subject can publish events
//	privateChan := make(chan interface{})
//	subject.observerToChanMap[observer] = privateChan
//	observer.receive(subject.done, privateChan)
//}
//
//func (subject *Subject) detach(observer Observer) {
//	privateChan := subject.observerToChanMap[observer]
//	defer close(privateChan)
//
//	delete(subject.observerToChanMap, observer)
//
//}
//
//func (subject *Subject) publish(val interface{}) {
//	subject.publicChan <- val
//}
//
//func (monitor *Monitor) receive(done chan interface{}, receiveChan chan interface{}) {
//	monitor.privateChan = receiveChan
//	go func() {
//		for {
//			select {
//			case val, ok := <-monitor.privateChan:
//				if !ok {
//					return
//				}
//				fmt.Println("Received val from subject", val, ok)
//			case <-done:
//				fmt.Println("Subject closed so will stop receiving")
//				close(monitor.privateChan)
//				return
//			}
//		}
//	}()
//}
//
//func main() {
//	type Stock struct {
//		subject *Subject
//	}
//	type StockerTicker struct {
//		monitor *Monitor
//	}
//
//	appleStock := &Stock{subject: NewSubject()}
//
//	moneyControl := &StockerTicker{monitor: NewMonitor()}
//	appleStock.subject.attach(moneyControl.monitor)
//
//	appleStock.subject.publish("Apple stock price is $230")
//	time.Sleep(1 * time.Millisecond)
//
//	economicTimes := NewMonitor()
//	appleStock.subject.attach(economicTimes)
//
//	appleStock.subject.publish("Apple stock price is $240")
//	time.Sleep(1 * time.Millisecond)
//
//	appleStock.subject.detach(moneyControl.monitor)
//	time.Sleep(1 * time.Millisecond)
//
//	appleStock.subject.publish("Apple stock price is $250")
//	time.Sleep(1 * time.Millisecond)
//	close(appleStock.subject.done)
//	time.Sleep(1 * time.Second)
//
//}
