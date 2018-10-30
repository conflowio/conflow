package util

// StringSliceContains tells whether a contains x.
func StringSliceContains(slice []string, search string) bool {
	for _, str := range slice {
		if str == search {
			return true
		}
	}
	return false
}
