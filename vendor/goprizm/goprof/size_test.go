package goprof

import (
	"goprizm/log"
	"testing"
)

type aStruct struct {
	name string
	b    *bStruct
}

type bStruct struct {
	name string
	a    *aStruct
}

func TestSizeof(t *testing.T) {
	b := &bStruct{
		name: "B",
	}
	cyclicStruct := &aStruct{
		name: "A",
		b:    b,
	}
	b.a = cyclicStruct

	tests := []struct {
		name string
		typ  interface{}
		size uint64
	}{
		{
			name: "int32",
			typ:  int32(123),
			size: 4,
		},
		{
			name: "int64",
			typ:  int64(123),
			size: 8,
		},
		{
			name: "bool",
			typ:  true,
			size: 1,
		},
		{
			name: "string",
			typ:  "Madhu",
			size: 21,
		},
		{
			name: "map",
			typ:  map[int]int{1: 1, 2: 2},
			size: 40,
		},
		{
			name: "map[int]map[int]int",
			typ:  map[int]map[int]int{1: {1: 1}, 2: {2: 2}},
			size: 72,
		},
		{
			name: "struct",
			typ: struct {
				a string
				b int
				c bool
				d map[int]int
			}{},
			size: 40,
		},
		{
			name: "cyclic struct",
			typ:  cyclicStruct,
			size: 83,
		},
		{
			name: "slice",
			typ:  []int{1, 2, 3},
			size: 48,
		},
		{
			name: "array",
			typ:  make([]int, 3),
			size: 48,
		},
		{
			name: "ptr",
			typ: &struct {
				a string
				b int
				c bool
				d map[int]int
			}{},
			size: 48,
		},
	}
	for _, test := range tests {
		log.Printf("compute size of %s", test.name)
		size := Sizeof(test.typ)
		if size != test.size {
			t.Fatalf("(%s) actual size: %d expected size: %d", test.name, size, test.size)
		}
	}
}
