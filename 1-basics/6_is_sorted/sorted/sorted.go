package sorted

import (
	"errors"
)

var ErrNotSorted = errors.New("strings are not sorted")

func Check(strings []string) error {
	prevStr := "" // the smallest string

	for _, str := range strings {
		if str < prevStr {
			return ErrNotSorted
		}

		prevStr = str
	}

	return nil
}

func tmp() {
	println(123)

}
