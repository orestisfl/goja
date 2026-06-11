//go:build !goja_reflect_methods

package goja

import (
	"reflect"
	"testing"
	"time"
)

// Under the default build, methods that were not registered via
// RegisterNativeMethods stay absent (the exposure is an explicit allowlist).
// This is build-tag specific: with goja_reflect_methods the full method set is
// exposed reflectively, so an unregistered method is reachable.
func TestRegisterNativeMethodsUnregistered(t *testing.T) {
	vm := New()
	vm.RegisterNativeMethods(reflect.TypeOf(time.Time{}), map[string]NativeMethod{
		"Unix": func(this interface{}, call FunctionCall) Value {
			return vm.ToValue(this.(time.Time).Unix())
		},
	})
	vm.Set("ts", time.Now())
	if _, err := vm.RunString(`ts.Month()`); err == nil {
		t.Fatal("expected a TypeError for an unregistered method")
	}
}
