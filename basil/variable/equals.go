package variable

// Equals returns true if the two given values are deeply equal
func Equals(v1 interface{}, v2 interface{}) bool {
	type1 := GetType(v1)
	type2 := GetType(v2)
	if type1 != type2 {
		return false
	}

	switch v1t := v1.(type) {
	case []interface{}:
		v2t := v2.([]interface{})
		if len(v1t) != len(v2t) {
			return false
		}
		for i := range v1t {
			if !Equals(v1t[i], v2t[i]) {
				return false
			}
		}
		return true
	case map[string]interface{}:
		v2t := v2.(map[string]interface{})
		if len(v1t) != len(v2t) {
			return false
		}
		for key, v1e := range v1t {
			v2e, ok := v2t[key]
			if !ok {
				return false
			}
			if !Equals(v1e, v2e) {
				return false
			}
		}
		return true
	case Union:
		v2t := v2.(Union)
		return Equals(v1t.Value(), v2t.Value())
	default:
		return v1 == v2
	}
}
