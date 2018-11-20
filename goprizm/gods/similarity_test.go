package gods

import "testing"

func TestPercentileProbablity(t *testing.T) {
	tests := map[float64]float64{
		1.0:  0.0,
		5.0:  0.0,
		9.0:  0.22,
		10.0: 0.22,
		20.0: 0.59,
		25.0: 0.59,
		35.0: 0.59,
		45.0: 0.99,
		50.0: 0.99,
		55.0: 0.99,
		65.0: 0.59,
		75.0: 0.59,
		80.0: 0.59,
		90.0: 0.22,
		95.0: 0.0,
		99.0: 0.0,
	}

	percentiles := []Percentile{
		{5, 5},
		{10, 10},
		{25, 25},
		{50, 50},
		{75, 75},
		{90, 90},
		{95, 95},
	}

	for value, expScore := range tests {
		score := PercentileProbablity(percentiles, value)
		if TruncateFloat(score) != expScore {
			t.Logf("PercentileProbablity value:%v expScore:%v != score:%v", value, expScore, score)
		}
	}
}

func TestPercentileSimilarity(t *testing.T) {
	tests := []struct {
		pct1, pct2 []Percentile
		score      float64
	}{
		{
			[]Percentile{{1, 2}, {25, 4}, {50, 6}, {75, 8}, {9, 10}},
			[]Percentile{{1, 2}, {25, 4}, {50, 6}, {75, 8}, {9, 10}},
			1.0,
		},
		{
			[]Percentile{{1, 2}, {25, 4}, {50, 6}, {75, 8}, {9, 10}},
			[]Percentile{{1, 1}, {25, 3}, {50, 5}, {75, 7}, {9, 9}},
			0.98,
		},
		{
			[]Percentile{{1, 2}, {25, 4}, {50, 6}, {75, 8}, {9, 10}},
			[]Percentile{{1, 12}, {25, 14}, {50, 16}, {75, 18}, {9, 20}},
			0.67,
		},
	}

	for _, tt := range tests {
		if score := PercentileSimilarity(tt.pct1, tt.pct2); score != tt.score {
			t.Fatalf("PercentileSimilarity(Sorensen) tt:%+v score:%v failed", tt, score)
		}
	}
}

func TestListSimilarity(t *testing.T) {
	tests := []struct {
		l1, l2  []string
		jScore  float64
		sdScore float64
	}{
		{
			[]string{},
			[]string{},
			0.0,
			0.0,
		},
		{
			[]string{"b", "a", "c"},
			[]string{"b", "a", "c"},
			1.0,
			1.0,
		},
		{
			[]string{"a", "b", "c", "d", "e"},
			[]string{"a", "b", "c", "d", "x"},
			0.66,
			0.8,
		},
		{
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c", "d", "x"},
			0.6,
			0.75,
		},
	}

	for _, tt := range tests {
		if score := JaccardSimilarity(tt.l1, tt.l2); score != tt.jScore {
			t.Fatalf("JaccardSimilarity tt:%+v score:%v failed", tt, score)
		}
		if score := SorensenDiceSimilarity(tt.l1, tt.l2); score != tt.sdScore {
			t.Fatalf("SorensenDiceSimilarity tt:%+v score:%v failed", tt, score)
		}
	}
}
