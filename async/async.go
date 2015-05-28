package async

import (
	"errors"
	"github.com/wayt/happyngine"
	"reflect"
)

var (
	// precomputed types
	contextType = reflect.TypeOf((*happyngine.Context)(nil))
)

type Function struct {
	fv reflect.Value // Kind() == reflect.Func
}

type FunctionResult struct {
	errs   chan error
	result []reflect.Value
}

func New(i interface{}) *Function {

	f := &Function{fv: reflect.ValueOf(i)}

	t := f.fv.Type()
	if t.Kind() != reflect.Func {
		panic(errors.New("not a function"))
	}
	if t.NumIn() == 0 || t.In(0) != contextType {
		panic(errors.New("first argument must be *happyngine.Context"))
	}

	return f
}

func (f *Function) Call(c *happyngine.Context, args ...interface{}) *FunctionResult {

	ft := f.fv.Type()
	in := []reflect.Value{reflect.ValueOf(c)}
	for _, arg := range args {
		var v reflect.Value
		if arg != nil {
			v = reflect.ValueOf(arg)
		} else {
			// Task was passed a nil argument, so we must construct
			// the zero value for the argument here.
			n := len(in) // we're constructing the nth argument
			var at reflect.Type
			if !ft.IsVariadic() || n < ft.NumIn()-1 {
				at = ft.In(n)
			} else {
				at = ft.In(ft.NumIn() - 1).Elem()
			}
			v = reflect.Zero(at)
		}
		in = append(in, v)
	}

	result := new(FunctionResult)
	result.errs = make(chan error)

	go func() {

		result.result = f.fv.Call(in)

		close(result.errs)
	}()

	return result
}

func (r *FunctionResult) Wait() []reflect.Value {

	<-r.errs

	return r.result
}
