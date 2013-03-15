package util

// Splits a string into a slice of space (U+0020) separated rune slices.
func SplitForUI(s string) (split [][]rune) {
	runes := []rune(s)

	start := -1
	for i, r := range runes {
		if r == ' ' {
			if start != -1 {
				split = append(split, runes[start:i])
				start = -1
			}
		} else {
			if start == -1 {
				start = i
			}
		}
	}

	if start != -1 {
		split = append(split, runes[start:])
	}

	return
}
