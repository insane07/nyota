package goprof

import (
	"reflect"
	"unsafe"
)

type visit struct {
	a1  unsafe.Pointer
	typ reflect.Type
}

func isComplexKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Array, reflect.Slice, reflect.Interface, reflect.Ptr, reflect.Struct, reflect.Map, reflect.String:
		return true
	}

	return false
}

func deepSizeof(v1 reflect.Value, visited map[visit]bool, skip bool) uint64 {
	var size uint64

	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}

	if v1.CanAddr() && hard(v1.Kind()) {
		if v1.IsNil() {
			return 0
		}

		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, typ}
		if visited[v] {
			return 0
		}

		// Remember for later.
		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v1.Len(); i++ {
			size = size + deepSizeof(v1.Index(i), visited, false)
		}
	case reflect.Interface:
		size = deepSizeof(v1.Elem(), visited, false)
	case reflect.Ptr:
		size = deepSizeof(v1.Elem(), visited, false)
	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if !isComplexKind(v1.Field(i).Kind()) {
				continue
			}

			//size = size - int64(v1.Field(i).Type().Size())
			size = size + deepSizeof(v1.Field(i), visited, true)
		}
	case reflect.Map:
		for _, k := range v1.MapKeys() {
			size = size + deepSizeof(k, visited, false)
			size = size + deepSizeof(v1.MapIndex(k), visited, false)
		}
	case reflect.String:
		size = uint64(len(v1.String()))
	}

	if !skip {
		size = size + uint64(v1.Type().Size())
	}

	return size
}

/*
Compute number of bytes occupied by variable of any arbitrary type
Compute size by visiting all the nested attributes
*/
func Sizeof(x interface{}) uint64 {
	if x == nil {
		return 0
	}

	return deepSizeof(reflect.ValueOf(x), make(map[visit]bool), false)
}
