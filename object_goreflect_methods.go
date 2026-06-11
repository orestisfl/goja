//go:build goja_reflect_methods

package goja

import (
	"go/ast"
	"reflect"
)

// This variant is selected by the OPT-IN goja_reflect_methods build tag. It
// restores goja's historical behaviour: Go methods are exposed to JavaScript on
// values handed to a Runtime (e.g. via ToValue).
//
// It relies on reflect.Type.Method / reflect.Value.Method to look methods up by
// a name chosen at runtime. The Go linker treats any reachable call to these as
// "the program may call an arbitrary method by reflection" and therefore
// retains every exported method of every reachable type, disabling method-level
// dead-code elimination for the whole binary. That is why method reflection is
// NOT the default: build without this tag to compile the no-op variant in
// object_goreflect_nomethods.go and keep dead-code elimination enabled (at the
// cost of not exposing Go methods to JavaScript).

// buildMethodInfo records the exported methods of t in info so they can be
// resolved by name from JavaScript.
func (r *Runtime) buildMethodInfo(t reflect.Type, info *reflectTypeInfo) {
	n := t.NumMethod()
	info.Methods = make(map[string]int, n)
	info.MethodNames = make([]string, 0, n)
	for i := 0; i < n; i++ {
		method := t.Method(i)
		name := method.Name
		if !ast.IsExported(name) {
			continue
		}
		if r.fieldNameMapper != nil {
			name = r.fieldNameMapper.MethodName(t, method)
			if name == "" {
				continue
			}
		}

		if _, exists := info.Methods[name]; !exists {
			info.MethodNames = append(info.MethodNames, name)
		}

		info.Methods[name] = i
	}
}

// reflectGetMethod returns the bound method value at idx of v.
func reflectGetMethod(v reflect.Value, idx int) reflect.Value {
	return v.Method(idx)
}
