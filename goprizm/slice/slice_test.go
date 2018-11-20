package slice

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		slice  []string
		value  string
		result bool
	}{
		{[]string{}, "a", false},
		{[]string{"a", "b"}, "a", true},
		{[]string{"a", "b"}, "c", false},
	}
	for _, tt := range tests {
		if result := Contains(tt.slice, tt.value); result != tt.result {
			t.Fatalf("Contains tt:%+v result:%v failed", tt, result)
		}
	}
}

func TestAppendUnique(t *testing.T) {
	tests := []struct {
		slice  []string
		values []string
		result []string
		ok     bool
	}{
		{[]string{}, []string{"a"}, []string{"a"}, true},
		{[]string{"a", "b"}, []string{"a"}, []string{"a", "b"}, false},
		{[]string{"a", "b"}, []string{"c"}, []string{"a", "b", "c"}, true},
	}
	for _, tt := range tests {
		if result, ok := AppendUnique(tt.slice, tt.values...); ok != tt.ok || !reflect.DeepEqual(result, tt.result) {
			t.Fatalf("AppendUnique tt:%+v result:%v ok:%v failed", tt, result, ok)
		}
	}
}
