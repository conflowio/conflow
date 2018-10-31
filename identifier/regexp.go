package identifier

import "regexp"

// RegExpPattern is the regular expression for a valid identifier
const RegExpPattern = "[a-z][a-z0-9]*(?:_[a-z0-9]+)*"

// RegExp is a compiled regular expression for a valid identifier
var RegExp = regexp.MustCompile("^" + RegExpPattern + "$")
