package gods

var (
	// zScores - percentils to abs(zscores)
	zScores = map[int]float64{
		5:  1.65,
		10: 1.28,
		20: 0.84,
		25: 0.67,
		30: 0.52,
		40: 0.25,
		50: 0.01,
		60: 0.25,
		70: 0.52,
		75: 0.67,
		80: 0.84,
		90: 1.28,
		95: 1.65,
	}

	// zScoreMax - zscores of percentils 5 and 95.
	zScoreMax = 1.65
)

type Percentile struct {
	I     int     `json:"i"`
	Value float64 `json:"value"`
}

// PercentileProbablity - given percentiles for a metric and value, return a score between 0 and 1
// which indicate probablity of value being part of given percentile distribution.
func PercentileProbablity(pcts []Percentile, value float64) (score float64) {
	var (
		pct Percentile
		i   int
	)
	for i, pct = range pcts {
		if pct.Value >= value {
			break
		}
	}

	if i == len(pcts) {
		return 0.0
	}

	if i != 0 {
		if (value - pcts[i-1].Value) < (pcts[i].Value - value) {
			i = i - 1
		}
	}

	return (1.0 - zScores[pcts[i].I]/zScoreMax)
}

// PercentileSimilarity - compare percentiles using
// https://datascience.stackexchange.com/questions/6898/users-percentile-similarity-measure
// - Sorensen Dice coefficient = 2 |A . B| / |A||A| + |B||B|
func PercentileSimilarity(pct1, pct2 []Percentile) float64 {
	var dotProd float64
	for i := range pct1 {
		dotProd += pct1[i].Value * pct2[i].Value
	}

	magnitude := func(vector []Percentile) float64 {
		var m float64
		for _, pct := range vector {
			m += (pct.Value * pct.Value)
		}
		return m
	}

	mag := magnitude(pct1) + magnitude(pct2)
	if mag == 0.0 {
		return 0.0
	}

	return TruncateFloat((2 * dotProd) / mag)
}

func JaccardSimilarity(list1, list2 []string) float64 {
	common := float64(sliceCommon(list1, list2))

	total := float64(len(list1) + len(list2))
	if int(total-common) == 0 {
		return 0.0
	}
	return TruncateFloat(common / (total - common))
}

// SorensenDiceSimilarity - 2 * (A intersect B) / (|A| + |B|)
func SorensenDiceSimilarity(list1, list2 []string) float64 {
	common := float64(sliceCommon(list1, list2))
	total := float64(len(list1) + len(list2))
	if int(total-common) == 0 {
		return 0.0
	}
	return TruncateFloat(2 * common / (total))
}

func sliceCommon(list1, list2 []string) (common int) {
	set1 := NewStringSet(list1...)
	for _, v2 := range list2 {
		if _, ok := set1[v2]; ok {
			common += 1
		}
	}
	return common
}

func TruncateFloat(f float64) float64 {
	return float64(int(f*100.0)) / 100.0
}
