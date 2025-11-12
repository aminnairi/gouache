// Package number defines a bunch of useful methods for working with numbers
package number

func IntBetween(lowerbound, upperbound, value int) bool {
	if value < lowerbound {
		return false
	}

	if value > upperbound {
		return false
	}

	return true
}
