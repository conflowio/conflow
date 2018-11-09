package function

import (
	"encoding/base64"

	"github.com/opsidian/basil/function"
)

// Base64Decode returns the string represented by the base64 string s.
//go:generate basil generate Base64Decode
func Base64Decode(src string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", function.NewError(0, err)
	}

	return string(res), nil
}
