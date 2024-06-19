package server

import "testing"

type testObj struct {
	name string
	num  int
}

func (t *testObj) foo() string {
	return t.name
}

func (t *testObj) bar() int {
	return t.num
}

func newService() *services {
	return &services{}
}
func TestRegisterService(t *testing.T) {

	test := &testObj{
		name: "neve",
		num:  1,
	}

	s := newService()
	s.registerService("testService", test)
}
