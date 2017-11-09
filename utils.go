package carrot

func InSlice(str string, items []string) bool {
	for _, item := range items {
		if item == str {
			return true
		}
	}

	return false
}
