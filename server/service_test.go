package server

import (
	"fmt"
	"reflect"
	"testing"
)

type testObj struct {
	name string
	num  int
}

func (t *testObj) Foo() string {
	return t.name
}

func (t *testObj) Bar(x int) int {
	return t.num + x
}

func TestRegisterService(t *testing.T) {

	test := &testObj{
		name: "neve",
		num:  1,
	}

	s := newService()

	err := s.RegisterService("testService", test)

	if err != nil {
		t.Errorf("something's wrong %v ", err)
	}

}

// Test struct
type TestStruct struct {
	Value int
}

func (t *TestStruct) Add(a, b int) int {
	return a + b
}

func TestCall(t *testing.T) {

	testStruct := &TestStruct{Value: 10}
	/*add := func(t *TestStruct, a, b int) int {
		return a + b
	}*/
	/*
		s := &services{
			svc: map[string]*object{
				"TestService": {
					methods: map[string]*method{
						"Add": {
							rcver: reflect.ValueOf(testStruct),
							fn:    reflect.ValueOf(add),
						},
					},
				},
			},
		}
	*/
	s := newService()

	s.RegisterService("TestService", testStruct)
	// Debug information
	fmt.Printf("Add function type: %v\n", reflect.TypeOf(testStruct.Add))

	tests := []struct {
		name          string
		serviceMethod string
		args          []reflect.Value
		want          any
		wantErr       bool
	}{
		{
			name:          "Valid call",
			serviceMethod: "TestService_Add",
			args:          []reflect.Value{reflect.ValueOf(5), reflect.ValueOf(7)},
			want:          12,
			wantErr:       false,
		},
		{
			name:          "Service not found",
			serviceMethod: "NonexistentService_Add",
			args:          []reflect.Value{},
			want:          nil,
			wantErr:       true,
		},
		{
			name:          "Method not found",
			serviceMethod: "TestService_NonexistentMethod",
			args:          []reflect.Value{},
			want:          nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Call(tt.serviceMethod, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Call() got = %v, want %v", got, tt.want)
			}
		})
	}
}
