package main

import "fmt"

type pizza interface {
	getPrice() int
}

type SimpleCheeseTomatoPizza struct {
}

func (p *SimpleCheeseTomatoPizza) getPrice() int {
	return 15
}

type OnionToppingPizza struct {
	pizza *SimpleCheeseTomatoPizza
}

func (p *OnionToppingPizza) getPrice() int {
	return p.pizza.getPrice() + 100
}

type VegFeastPizza struct {
	pizza *OnionToppingPizza
}

func (p *VegFeastPizza) getPrice() int {
	return p.pizza.getPrice() + 100
}

func main() {
	pizza := &SimpleCheeseTomatoPizza{}
	pizzaWithOnionTopping := &OnionToppingPizza{pizza: pizza}
	vegFeastPizza := &VegFeastPizza{pizza: pizzaWithOnionTopping}
	fmt.Println("Price of veg feast pizza=", vegFeastPizza.getPrice())

}
