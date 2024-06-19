package server

import (
	"errors"
	"reflect"
	"sync"
)

type services struct {
	mu  sync.Mutex
	srv map[string]*object
}

type object struct {
	name    string
	methods map[string]*method
}

type method struct {
	//method's name will be stored as a string in -> methods map[string]Method
	//arguments
	args  []reflect.Type
	fn    reflect.Value
	rcver reflect.Value
	//return values
	rn  reflect.Value
	err string
}

func (s *services) RegisterService(name string, object any) error {

	if _, ok := s.srv[name]; ok {
		return errors.New("service name already exists")
	}

	obj := reflect.ValueOf(object)
	s.srv[name] = getObj(obj)

}

func getObj(obj reflect.Value) *object {
	objType := obj.Type()

	return &object{
		name:    objType.Name(),
		methods: getMethods(obj),
	}
}

func getMethods(obj reflect.Value) map[string]*method {
	methods := make(map[string]*method)

	mtd := &method{
		rcver: obj,
	}

	for i := 0; i < obj.NumMethod(); i++ {
		fn := obj.Method(i)
		methods[fn.Type().Name()] = extractMethod(fn, mtd)
	}
	return methods
}

func extractMethod(fn reflect.Value, mtd *method) *method {
	fnType := fn.Type()

	args := make([]reflect.Type, fnType.NumIn()+1)

	for i := 0; i < fnType.NumIn(); i++ {
		args = append(args, fnType.In(i))
	}

	return &method{
		args: args,
		fn:   fn,
	}
}
