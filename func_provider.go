package dig

import (
	"fmt"
	"reflect"
)

type funcProviderHandler struct {
	target reflect.Type
}

func (handler funcProviderHandler) handle(args []reflect.Value) []reflect.Value {
	t := reflect.New(handler.target)

	for i := 0; i < t.Elem().Type().NumField(); i++ {
		if t.Elem().Type().Field(i).Anonymous {
			continue
		}
		field := t.Elem().Field(i)
		found := false
		for _, arg := range args {
			if arg.Type() == field.Type() {
				field.Set(arg)
				found = true
			}
		}
		if !found {
			return []reflect.Value{
				reflect.Zero(handler.target),
				reflect.ValueOf(fmt.Errorf("unable to find arg for field=%s", t.Elem().Type().Field(i).Name)),
			}
		}
	}

	return []reflect.Value{
		t,
		reflect.Zero(errorType()),
	}
}

func errorType() reflect.Type {
	return reflect.TypeOf(make([]error, 1)).Elem()
}

func funcProvider(t reflect.Type) interface{} {

	args := make([]reflect.Type, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.Anonymous {
			args = append(args, field.Type)
		}
	}

	fn := reflect.FuncOf(args, []reflect.Type{reflect.PtrTo(t), errorType()}, false)

	handler := funcProviderHandler{
		target: t,
	}
	newFunc := reflect.MakeFunc(fn, handler.handle)
	return newFunc.Interface()
}
