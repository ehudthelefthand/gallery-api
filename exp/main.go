package main

import "fmt"

type Duck struct {
	Name string
}

func (d *Duck) Fly() {
	fmt.Println(d.Name + " fly")
}

type Bird struct {
	Name string
}

func (b *Bird) Fly() {
	fmt.Println(b.Name + " fly like bird")
}

type Flyable interface {
	Fly()
}

func MakeItFly(d Flyable) {
	d.Fly()
}

func main() {
	duck := Duck{Name: "Pop"}
	// duck.Fly()
	MakeItFly(&duck)

	bird := Bird{Name: "Game"}
	// bird.Fly()
	MakeItFly(&bird)
}
