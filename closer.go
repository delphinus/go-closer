package closer

import "io"

// Close is for using on defer's.  It closes io.Closer and set errors on a
// supplied variable.
func Close(r io.Closer, err *error) {
	Check(func() error { return r.Close() }, err)
}

// Check is for using on defer's.  It checks the func error and set it on a
// supplied variable.
func Check(f func() error, err *error) {
	if ferr := f(); ferr != nil && *err == nil {
		*err = ferr
	}
}
