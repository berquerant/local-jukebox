package slicex

import "slices"

// Split splits a given slice into two parts based on the first occurrence of a specified value,
// returning the sub-slices that precede and follow the value.
func Split[S ~[]E, E comparable](s S, v E) (S, S) {
	i := slices.Index(s, v)
	if i < 0 {
		return s, nil
	}
	return s[:i], s[i+1:]
}
