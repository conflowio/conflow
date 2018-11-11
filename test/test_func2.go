package test

//go:generate basil generate
func testFunc2(str1 string, str2 string) string {
	return str1 + str2
}
