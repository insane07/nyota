package gods

import "sort"

type Count struct {
	Value string
	Count uint64
}

type Counter map[string]uint64

func NewCounter(values ...string) Counter {
	ctr := make(Counter)
	for _, value := range values {
		ctr.Add(value, 1)
	}
	return ctr
}

func (ctr Counter) Add(value string, n uint64) {
	ctr[value] += uint64(n)
}

func (ctr Counter) Update(other Counter) {
	for value, count := range other {
		ctr[value] += count
	}
}

func (ctr Counter) TrimN(n int) Counter {
	top := ctr.TopN(n)
	ctrNew := NewCounter()
	for _, c := range top {
		ctrNew.Add(c.Value, uint64(c.Count))
	}
	return ctrNew
}

func (ctr Counter) TopN(n int) []Count {
	if n > len(ctr) {
		n = len(ctr)
	}

	var counts []Count
	for value, count := range ctr {
		counts = append(counts, Count{value, count})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	return append([]Count{}, counts[:n]...)
}
