package user

import "github.com/piapip/Learning-Go/Gomock/doer"

//User duper useless
type User struct {
	Doer doer.Doer
}

//Use same goes with this
func (u *User) Use() error {
	return u.Doer.DoSomething(123, "Hello GoMock")
}

func (u *User) take(x, y int) int {
	return u.Doer.DoThisToo(x, y)
}
