package gods

import (
	"reflect"
	"sort"
	"testing"
)

func TestStringSet(t *testing.T) {
	ss := NewStringSet("a", "b")
	ss.Add("c", "d")
	ss.Add("c", "d")

	slice := ss.ToSlice()
	sort.Strings(slice)

	if len(ss) != 4 || !reflect.DeepEqual(slice, []string{"a", "b", "c", "d"}) {
		t.Fatalf("StringSet len:%v != 4 slice:%v", len(ss), slice)
	}

	ss.Remove("c")
	slice = ss.ToSlice()
	sort.Strings(slice)
	if len(ss) != 3 || !reflect.DeepEqual(slice, []string{"a", "b", "d"}) {
		t.Fatalf("StringSet len:%v != 3 slice:%v", len(ss), slice)
	}
}
