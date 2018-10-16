package slice

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func AppendUnique(slice []string, values ...string) ([]string, bool) {
	var add []string
	for _, v := range values {
		if Contains(slice, v) {
			continue
		}
		add = append(add, v)
	}
	if len(add) == 0 {
		return slice, false
	}
	return append(slice, add...), true
}
