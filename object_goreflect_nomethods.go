//go:build !goja_reflect_methods

package goja

import "reflect"

// This is the DEFAULT variant; it is compiled unless the goja_reflect_methods
// build tag is set. It contains no use of reflect.Type.Method /
// reflect.Value.Method, so the Go linker is free to perform method-level
// dead-code elimination. (Reachable index-based reflect method calls otherwise
// force every exported method of every reachable type to be retained, which can
// easily double a binary's size.)
//
// The trade-off: Go methods are NOT automatically exposed to JavaScript on
// values handed to a Runtime. Expose Go behaviour to JS via native functions
// instead, e.g. obj.Set("name", func(call goja.FunctionCall) goja.Value { ... }).
// Build with -tags goja_reflect_methods to restore reflective method exposure
// (object_goreflect_methods.go).

// buildMethodInfo records no methods: with this tag, Go methods are not exposed
// to JavaScript. info.Methods / info.MethodNames are left nil, which is a safe
// empty method set for the readers (map lookup miss in _getMethod, zero-length
// range in nextMethod).
func (r *Runtime) buildMethodInfo(t reflect.Type, info *reflectTypeInfo) {
}

// reflectGetMethod is never reached with this build tag (buildMethodInfo records
// no methods, so objectGoReflect._getMethod never finds a match). It exists only
// so the package compiles, and deliberately avoids reflect.Value.Method.
func reflectGetMethod(v reflect.Value, idx int) reflect.Value {
	return reflect.Value{}
}
