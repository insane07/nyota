package gods

import (
	"reflect"
	"strings"
	"testing"
)

func counterAdd(ctr Counter, valueSets [][]string) {
	for _, values := range valueSets {
		for _, v := range values {
			ctr.Add(v, 1)
		}
	}
}

func TestCounter(t *testing.T) {
	valueSets := [][]string{
		strings.Fields(strings.Repeat("m ", 5)),
		strings.Fields(strings.Repeat("q ", 45)),
		strings.Fields(strings.Repeat("d ", 95)),
		strings.Fields(strings.Repeat("m ", 35)),
		strings.Fields(strings.Repeat("m ", 25)),
	}

	ctr := NewCounter()
	counterAdd(ctr, valueSets)

	exp := []Count{
		{"d", 95},
		{"m", 65},
		{"q", 45},
	}
	if top := ctr.TopN(3); !reflect.DeepEqual(top, exp) {
		t.Fatalf("TopN exp:%v got:%v", exp, top)
	}

	valueSets = [][]string{
		strings.Fields(strings.Repeat("d ", 100)),
		strings.Fields(strings.Repeat("s ", 85)),
		strings.Fields(strings.Repeat("o ", 15)),
	}
	other := NewCounter()
	counterAdd(other, valueSets)

	ctr.Update(other)
	exp = []Count{
		{"d", 195},
		{"s", 85},
		{"m", 65},
	}
	if top := ctr.TopN(3); !reflect.DeepEqual(top, exp) {
		t.Fatalf("TopN exp:%v got:%v", exp, top)
	}

	expCtr := Counter(map[string]uint64{
		"d": 195,
		"s": 85,
	})
	ctr = ctr.TrimN(2)
	if !reflect.DeepEqual(expCtr, ctr) {
		t.Fatalf("Trim exp:%v got:%v", expCtr, ctr)
	}
}
