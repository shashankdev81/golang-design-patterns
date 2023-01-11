package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

//https://brandur.org/go-worker-pool
type Task struct {
	dummy   chan string
	result  interface{}
	err     error
	request interface{}
	f       func(req interface{}) (interface{}, error)
}

func NewTask(f func(r interface{}) (interface{}, error), req interface{}) *Task {
	testChan := make(chan string, 1)
	return &Task{dummy: testChan, request: req, f: f}
}

func (task *Task) run(wg *sync.WaitGroup) {
	res, err := task.f(task.request)
	task.result = res
	task.err = err
	wg.Done()
}

type Pool struct {
	tasks       []*Task
	concurrency int
	tasksChan   chan *Task
	wg          sync.WaitGroup
}

func NewPool(tasks []*Task, con int) *Pool {
	return &Pool{tasks: tasks, concurrency: con, tasksChan: make(chan *Task)}
}

func (pool *Pool) Run() {
	//many consumers - start go routine workers that can pull tasks and start working
	pool.wg.Add(len(pool.tasks))

	for i := 0; i < pool.concurrency; i++ {
		go pool.Work()
	}

	//single producer - add tasks to tasks chan
	for _, task := range pool.tasks {
		pool.tasksChan <- task
	}
	close(pool.tasksChan)
	pool.wg.Wait()
}

func (pool *Pool) Work() {
	for task := range pool.tasksChan {
		task.run(&pool.wg)
	}
}

func (pool *Pool) PrintResults() {
	for _, task := range pool.tasks {
		if task.err != nil {
			fmt.Println("Error occurred for ", task.request, " ", task.err)
		} else {
			fmt.Println("Result for ", task.request, " ", task.result)
		}
	}
}

func main() {
	done := make(chan interface{})
	http.HandleFunc("/start", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		go LotOfConcurrentWork(done)
		//fmt.Printf(writer, "hello")
	})

	http.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		close(done)
	})
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatalln(err)
	}
}

func LotOfConcurrentWork(done chan interface{}) {
	for i := 0; i < 100; i++ {
		select {
		case <-done:
			break
		default:
			go work(i)
		}

	}
}

func work(i int) {
	println("Count=", i)
	divide := func(object interface{}) (interface{}, error) {
		a, _ := object.(int)
		if a != 0 {
			return a / 2, nil
		} else {
			return -1, errors.New("Divide by zero")
		}
	}
	tasks := []*Task{
		NewTask(divide, 4),
		NewTask(divide, 0),
		NewTask(divide, 6),
		NewTask(divide, 8),
		NewTask(divide, 10),
	}

	pool := NewPool(tasks, 2)
	pool.Run()
	pool.PrintResults()
}
