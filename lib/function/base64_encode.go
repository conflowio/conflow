package function

import (
	"encoding/base64"
)

// Base64Encode returns the base64 encoding of src.
//go:generate basil generate Base64Encode
func Base64Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}
