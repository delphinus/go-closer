package closer

import "io"

// Close is for using on defer's.  It closes io.Closer and set errors on a
// supplied variable.
func Close(r io.Closer, err *error) {
	if cerr := r.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
