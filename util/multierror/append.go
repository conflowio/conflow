package multierror

func Append(e1 error, e2 error) error {
	switch {
	case e1 == nil:
		return e2
	case e2 == nil:
		return e1
	}

	if e, ok := e1.(*Error); ok {
		return e.Append(e2)
	}

	e := &Error{errors: []error{e1}}
	return e.Append(e2)
}
