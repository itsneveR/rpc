package server

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type services struct {
	mu  sync.Mutex
	svc map[string]*object
}

type object struct {
	name    string // dcw the servce name
	methods map[string]*method
}

type method struct {
	//arguments
	fn    reflect.Value
	args  []reflect.Type
	rcver reflect.Value // == object
	//return values
	rn  []reflect.Type
	err string
}

func (s *services) RegisterService(name string, obj any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.svc == nil {
		s.svc = make(map[string]*object)
	}

	_, ok := s.svc[name]

	if ok {
		return errors.New("service name already exists")
	}

	s.svc[name] = getObject(reflect.ValueOf(obj))

	return nil
}

func getObject(obj reflect.Value) *object {

	return &object{
		name:    obj.Type().Name(), //struct/servive name
		methods: getMethods(obj),
	}
}

func getMethods(obj reflect.Value) map[string]*method {

	methods := make(map[string]*method)

	for i := 0; i < obj.NumMethod(); i++ {

		method := obj.Type().Method(i)

		methods[method.Name] = extractMethodInfo(method, obj)

	}
	return methods
}

func extractMethodInfo(fn reflect.Method, obj reflect.Value) *method {

	args := make([]reflect.Type, fn.Type.NumIn())

	for i := 0; i < fn.Type.NumIn(); i++ {
		args = append(args, fn.Type.In(i))
	}

	rn := make([]reflect.Type, fn.Type.NumIn())

	for i := 0; i < fn.Type.NumOut(); i++ {
		rn = append(rn, fn.Type.Out(i))
	}

	return &method{
		fn:    fn.Func,
		args:  args,
		rcver: obj,
		rn:    rn,
		//err:   rn[fn.Type.NumOut()].String(),
	}
}

// ServiceName_method
func (s *services) Call(service_method string, args []reflect.Value) (any, error) {
	srvc, method, _ := strings.Cut(service_method, "_")

	object, ok := s.svc[srvc]

	if !ok {
		return nil, errors.New("srvice not found")
	}

	fn, ok := object.methods[method]

	if !ok {
		return nil, errors.New("method not found in this service")
	}

	// Debug information
	fmt.Printf("Method type: %v\n", fn.fn.Type())
	fmt.Printf("Number of input parameters: %d\n", fn.fn.Type().NumIn())
	fmt.Printf("Number of provided args: %d\n", len(args))
	fmt.Printf("Method type:%d\n", fn.rcver.Type())

	// Check if we have the correct number of arguments
	if fn.fn.Type().NumIn() != len(args)+1 { // +1 for receiver
		return nil, fmt.Errorf("incorrect number of arguments. Expected %d, got %d", fn.fn.Type().NumIn()-1, len(args))
	}

	fullArgs := make([]reflect.Value, 0, len(args)+1)

	fullArgs = append(fullArgs, fn.rcver)

	fullArgs = append(fullArgs, args...)

	result := fn.fn.Call(fullArgs)

	if len(result) == 0 {
		return nil, nil
	}

	return result[0].Interface(), nil
}
