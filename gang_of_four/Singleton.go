package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

//http://marcio.io/2015/07/singleton-pattern-in-go/
type Singleton struct {
	id string
}

var instance *Singleton
var once sync.Once

func getInstance() *Singleton {
	once.Do(func() {
		idStr := strconv.Itoa(rand.Intn(100))
		instance = &Singleton{id: idStr}
		fmt.Println("Created instance with id", instance.id)
	})
	return instance
}
func main() {
	var wg sync.WaitGroup
	fetchInstance := func(wg *sync.WaitGroup) {
		inst := getInstance()
		fmt.Println("Fetched instance with id", inst.id)
		wg.Done()
	}
	wg.Add(3)
	go fetchInstance(&wg)
	go fetchInstance(&wg)
	go fetchInstance(&wg)
	wg.Wait()

}
