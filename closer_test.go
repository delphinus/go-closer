// +build go1.7

package closer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type someIO struct {
	err error
}

func (r someIO) Close() error {
	return r.err
}

var (
	errA = errors.New("A! error")
	errB = errors.New("B! error")
)

func TestClose(t *testing.T) {
	for _, c := range []struct {
		name       string
		errOnFunc  error
		errOnClose error
		want       error
	}{
		{
			name:       "if no error, return nil",
			errOnFunc:  nil,
			errOnClose: nil,
			want:       nil,
		},
		{
			name:       "if func has no error & close has error, want close error",
			errOnFunc:  nil,
			errOnClose: errA,
			want:       errA,
		},
		{
			name:       "if func has error & close has no error, want func error",
			errOnFunc:  errA,
			errOnClose: nil,
			want:       errA,
		},
		{
			name:       "if func has error & close has error, want func error",
			errOnFunc:  errA,
			errOnClose: errB,
			want:       errA,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			a := assert.New(t)
			r := someIO{c.errOnClose}
			someFunc := func() (err error) {
				defer Close(r, &err)
				return c.errOnFunc
			}
			a.Equal(c.want, someFunc())
		})
	}
}
