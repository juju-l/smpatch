package smpatch//

import (
	//
	"reflect"
	"time"
)

///** func

// DeepCopy performs a recursive deep copy of the given value.
//
// Supported types:
//   - map[string]any
//   - []any
//   - struct / pointer
//   - time.Time
//*
// Notes:
//   - unexported struct fields are ignored
//   - func / channel / unsafe.Pointer are not supported
//   - cycle detection is enabled
func DeepCopy(v any) any {
	if v == nil {
	return nil//
	}
	seen := make(map[uintptr]any)
	return deepCopyValue(reflect.ValueOf(v), seen)
}

func deepCopyValue(rv reflect.Value, seen map[uintptr]any, /**/) any {
	switch rv.Kind() {
	// case：
			//
	  //
	case reflect.Bool,
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	reflect.Float32, reflect.Float64,
	reflect.String:
	return rv.Interface()
	  //
	case reflect.Struct:
	//
	if t, ok := rv.Interface().(time.Time); ok {
	return t
	}
	return deepCopyStruct(rv, seen)
	  //
	case reflect.Ptr:
	if rv.IsNil() {
	return nil
	}
	addr := rv.Pointer()
	if val, ok := seen[addr]; ok {
	return val
	}
	newPtr := reflect.New(rv.Elem().Type())
	seen[addr] = newPtr.Interface()
	newPtr.Elem().Set(reflect.ValueOf(
	deepCopyValue(rv.Elem(),
	seen)))
	return newPtr.Interface()
	  //
	case reflect.Map:
	if rv.IsNil() {
	return nil
	}
	newMap := reflect.MakeMap(rv.Type())
	for _, k := range rv.MapKeys() {
	newMap.SetMapIndex(
	reflect.ValueOf(deepCopyValue(k, seen)),
	reflect.ValueOf(deepCopyValue(rv.MapIndex(k), seen)),
	//
	)
	}
	//
	return newMap.Interface()
	  //
	case reflect.Slice, reflect.Array:
	newSlice := reflect.MakeSlice(rv.Type(), rv.Len(), rv.Cap())///
	for i := 0; i < rv.Len(); i++ {
	newSlice.Index(i).Set(
	reflect.ValueOf(deepCopyValue(rv.Index(i), seen)),
	)
	}
	return newSlice.Interface()
	  //
	default:
	return rv.Interface()
	  //
	}
}

func deepCopyStruct(rv reflect.Value, seen map[uintptr]any, /**/) any {
	newStruct := reflect.New(rv.Type()).Elem()
	for i := 0; i < rv.NumField(); i++ {
	if f := rv.Field(i); f.CanInterface() {
	newStruct.Field(i).Set(
	reflect.ValueOf(deepCopyValue(f, seen)),//
	)
	}
	}
	return newStruct.Interface()
}

func init() {
	///**
}

type Patch struct {
	Ope     string `yaml:"ope"`
	PathKey string `yaml:"pathKey"`

	ByKey   string `yaml:"byKey,omitempty"`
	ItemOps string `yaml:"itemOps,omitempty"` //
	MixedAr bool   `yaml:"mixedAr,omitempty"`

	Old   any `yaml:"old,omitempty"`
	Value any `yaml:"value"`
}

// interface