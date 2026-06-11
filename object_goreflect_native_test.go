package goja

import (
	"reflect"
	"testing"
	"time"
)

// These tests exercise RegisterNativeMethods, which must work under the DEFAULT
// build (no goja_reflect_methods tag) -- that is its whole purpose: expose Go
// methods to JS without reflect.Value.Method so the linker can still drop unused
// methods. The file deliberately has no build tag.

func TestRegisterNativeMethods(t *testing.T) {
	const SCRIPT = `
	ts.Unix() === 1257894000 && ts.Year() === 2009 && ("Unix" in ts);
	`
	vm := New()
	vm.RegisterNativeMethods(reflect.TypeOf(time.Time{}), map[string]NativeMethod{
		"Unix": func(this interface{}, call FunctionCall) Value {
			return vm.ToValue(this.(time.Time).Unix())
		},
		"Year": func(this interface{}, call FunctionCall) Value {
			return vm.ToValue(this.(time.Time).Year())
		},
	})
	vm.Set("ts", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))

	v, err := vm.RunString(SCRIPT)
	if err != nil {
		t.Fatal(err)
	}
	if !v.ToBoolean() {
		t.Fatalf("unexpected result: %v", v)
	}
}

// The wrapped value must remain a genuine time.Time, so Export()/round-tripping
// is unaffected by registering native methods.
func TestRegisterNativeMethodsExportRoundTrip(t *testing.T) {
	vm := New()
	vm.RegisterNativeMethods(reflect.TypeOf(time.Time{}), map[string]NativeMethod{
		"Unix": func(this interface{}, call FunctionCall) Value {
			return vm.ToValue(this.(time.Time).Unix())
		},
	})
	want := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	v, err := vm.RunString(`(function(t){ return t; })`)
	if err != nil {
		t.Fatal(err)
	}
	fn, ok := AssertFunction(v)
	if !ok {
		t.Fatal("not a function")
	}
	res, err := fn(Undefined(), vm.ToValue(want))
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := res.Export().(time.Time); !ok || !got.Equal(want) {
		t.Fatalf("round-trip failed: got %#v", res.Export())
	}
}
