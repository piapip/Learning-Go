package main

//Employee testing
type Employee struct {
	name string
	age  int
}

//NewEmployee initiation
func NewEmployee() *Employee {
	p := &Employee{}
	return p
}

//PrintEmployee something random
func PrintEmployee(p *Employee) string {
	return "Hello world!"
}
