package test

//go:generate basil generate testFunc2
func testFunc2(str1 string, str2 string) string {
	return str1 + str2
}
