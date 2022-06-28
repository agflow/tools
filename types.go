package main

import (
	"fmt"
	"reflect"
)

func BaseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = DeRef(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}

func DeRef(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
