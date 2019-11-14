package doer

//Doer doing nothing
type Doer interface {
	DoSomething(int, string) error
	DoThisToo(int, int) int
}
