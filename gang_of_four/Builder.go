package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

type Pizza struct {
	crust   string   `json:"crust"`
	base    string   `json:"base"`
	topping []string `json:"topping"`
	size    int      `json:"size"`
	isVeg   bool     `json:"isVeg"`
	isBurst bool     `json:"isBurst"`
}

func (pizza *Pizza) serve() {
	fmt.Printf("%+v\n", pizza)
}

type PizzaBuilder interface {
	chooseCrustAndBase(crust string, base string) PizzaBuilder
	addTopping(string) PizzaBuilder
	setSize(int) PizzaBuilder
	setIsVeg(bool) PizzaBuilder
	setIsBurst(bool) PizzaBuilder
	Build() (*Pizza, error)
}

type VegFeastPizzaBuilder struct {
	pizza *Pizza
}

func PizzaBuilderImpl() PizzaBuilder {
	simplePizza := &Pizza{topping: make([]string, 0)}
	return &VegFeastPizzaBuilder{pizza: simplePizza}
}

func (builder *VegFeastPizzaBuilder) chooseCrustAndBase(crust string, base string) PizzaBuilder {
	builder.pizza.crust = crust
	builder.pizza.base = base
	return builder
}

func (builder *VegFeastPizzaBuilder) addTopping(topping string) PizzaBuilder {
	builder.pizza.topping = append(builder.pizza.topping, topping)
	return builder

}

func (builder *VegFeastPizzaBuilder) setSize(size int) PizzaBuilder {
	builder.pizza.size = size
	return builder

}

func (builder *VegFeastPizzaBuilder) setIsVeg(isVeg bool) PizzaBuilder {
	builder.pizza.isVeg = isVeg
	return builder

}

func (builder *VegFeastPizzaBuilder) setIsBurst(isBurst bool) PizzaBuilder {
	builder.pizza.isBurst = isBurst
	return builder

}

func (builder *VegFeastPizzaBuilder) Build() (*Pizza, error) {
	return builder.pizza, nil
}

func main() {

	vegPizzaBuilder := PizzaBuilderImpl()
	pizza, ok := vegPizzaBuilder.chooseCrustAndBase("thin", "wheat").addTopping("onion").addTopping("tomato").addTopping("capsicum").setSize(9).setIsVeg(true).setIsBurst(false).Build()
	if ok != nil {
		panic("Pizza not baked well")
	}
	//fmt.Println("Veg feast pizza baked ", pizza)
	spew.Dump(pizza)

	mexicanNonVegPizzaBuilder := PizzaBuilderImpl()
	pizza, ok = mexicanNonVegPizzaBuilder.chooseCrustAndBase("thin", "wheat").addTopping("chicken").addTopping("onion").addTopping("jalapenos").addTopping("chillies").setSize(9).setIsVeg(false).setIsBurst(false).Build()
	if ok != nil {
		panic("Pizza not baked well")
	}
	//fmt.Println("Veg feast pizza baked ", pizza)
	spew.Dump(pizza)

}
