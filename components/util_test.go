package components

import (
	"regexp"
	"strings"

	"github.com/google/go-cmp/cmp"
)

var (
	// regexpComparer compares to strings, but if one of them starts with a "^", it will be used as a regexp
	regexpComparer = cmp.Comparer(func(x, y string) bool {
		switch {
		case strings.HasPrefix(x, "^"):
			return regexp.MustCompile(x).MatchString(y)
		case strings.HasPrefix(y, "^"):
			return regexp.MustCompile(y).MatchString(x)
		default:
			return x == y
		}
	})

	// positiveInteger returns true if the given integer is positive
	positiveInteger = func(i int) bool { return i > 0 }
)
