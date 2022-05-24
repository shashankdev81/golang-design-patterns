package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("Testing lfu")
	cache := initCache(&lfu{}, 3)
	test(cache)
	fmt.Println("")
	fmt.Println("Testing fifo")
	cache = initCache(&fifo{}, 3)
	test(cache)

}

func test(cache *cache) {
	cache.add("a", "1")
	cache.add("b", "2")
	cache.add("c", "3")
	cache.print()
	cache.get("a")
	cache.get("a")
	cache.get("c")
	cache.add("d", "4")
	cache.print()
}

type cache struct {
	storage      map[string]string
	evictionAlgo evictionAlgo
	currCapacity int
	maxCapacity  int
}

type evictionAlgo interface {
	evict(*cache)
	reorder(k string)
}

type fifo struct {
	//metadata for book-keeping
	order []string
}

//lifo implementation  of the eviction algorithm
func (l *fifo) evict(c *cache) {
	//remove first element from cache
	first := l.order[0]
	l.order = l.order[1:]
	delete(c.storage, first)
}

func (l *fifo) reorder(k string) {
	l.order = append(l.order, k)
}

//lru implementation  of the algorithm
type lfu struct {
	accessedTimes map[string]int
}
type AccessedTimesPair struct {
	key   string
	value int
}
type AccessedTimesPairList []AccessedTimesPair

func (p AccessedTimesPairList) Len() int {
	return len(p)
}

func (p AccessedTimesPairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p AccessedTimesPairList) Less(i, j int) bool {
	return p[i].value < p[j].value
}
func (l *lfu) evict(c *cache) {
	p := make(AccessedTimesPairList, len(l.accessedTimes))
	i := 0
	for k, v := range l.accessedTimes {
		p[i] = AccessedTimesPair{key: k, value: v}
		i++
	}
	//https://medium.com/@kdnotes/how-to-sort-golang-maps-by-value-and-key-eedc1199d944
	sort.Sort(p)
	deleteKey := p[0].key
	delete(l.accessedTimes, deleteKey)
	delete(c.storage, deleteKey)
}

func (l *lfu) reorder(k string) {
	if l.accessedTimes == nil {
		l.accessedTimes = make(map[string]int)
	}
	l.accessedTimes[k] = l.accessedTimes[k] + 1
}

type lru struct {
}

func (l *lru) evict(c *cache) {

}

func (l *lru) reorder(k string) {
}

func initCache(e evictionAlgo, capcity int) *cache {
	storage := make(map[string]string)
	return &cache{
		storage:      storage,
		evictionAlgo: e,
		currCapacity: 0,
		maxCapacity:  capcity,
	}
}

func (c *cache) add(k string, v string) {
	if c.currCapacity == c.maxCapacity {
		c.evictionAlgo.evict(c)
	} else {
		c.evictionAlgo.reorder(k)
		c.currCapacity++
	}
	c.storage[k] = v
}

func (c *cache) get(k string) string {
	val := c.storage[k]
	c.evictionAlgo.reorder(k)
	return val
}

func (c *cache) print() {
	fmt.Println("")
	for k, v := range c.storage {
		fmt.Printf("key=%v value=%v, ", k, v)
	}
}
