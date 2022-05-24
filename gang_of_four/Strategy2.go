package main

import "fmt"

//https://faun.pub/design-pattern-with-go-ft-strategy-pattern-a1efc58972e6

//core entity vigilante which is not allowed to be modified
type Vigilante struct {
	name   string
	weapon Weapon
}

func (v Vigilante) intro() {
	fmt.Printf("Hi my name is %v and I am here to kill zombies", v.name)
}

//contract for behaviour that is allowed to vary and which has been extracted out as a strategy
type Weapon interface {
	UseWeapon()
}

//how a particular type of wepaon behaves
type SwordWeapon struct {
	method string
}

func (w *SwordWeapon) UseWeapon() {
	fmt.Printf("A sword will be used to %v the zombies", w.method)
}

type GunWeapon struct {
	barell string
	method string
}

func (w *GunWeapon) UseWeapon() {
	fmt.Printf("A %v gun will be used to %v the zombies", w.barell, w.method)
}

type BombWeapon struct {
	method string
}

func (w *BombWeapon) UseWeapon() {
	fmt.Printf("A bomb will be used to %v the zombies", w.method)
}

func main() {
	v := Vigilante{
		name:   "Shashank",
		weapon: &GunWeapon{method: "shoot", barell: "double"},
	}
	v.intro()
	fmt.Println("")
	v.weapon.UseWeapon()
	v.weapon = &BombWeapon{method: "blow"}
	fmt.Println("")
	v.weapon.UseWeapon()
}
