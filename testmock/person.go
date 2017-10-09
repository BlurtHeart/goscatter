package test

import "fmt"

type Person struct {
	name string
}

func NewPerson(name string) *Person {
	return &Person{
		name: name,
	}
}

func (p *Person) Talk(name string) string {
	return fmt.Sprintf("Hello %s, My name is %s", name, p.name)
}
