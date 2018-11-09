package variable

// ID is the identifier type
type ID string

// String returns with the ID string
func (i ID) String() string {
	return string(i)
}
