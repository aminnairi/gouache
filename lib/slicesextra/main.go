// Package slicesextra helps enhancing the already existing slices package by providing additionnal and useful functions.
package slicesextra

func At[Type comparable](index int, fallback Type, items []Type) Type {
	itemsLength := len(items)

	if itemsLength == 0 {
		return fallback
	}

	if itemsLength < index {
		return fallback
	}

	return items[index]
}
