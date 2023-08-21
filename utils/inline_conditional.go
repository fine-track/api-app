package utils

func InlineConditional[T comparable](c bool, whenTrue T, whenFalse T) T {
	if c {
		return whenTrue
	} else {
		return whenFalse
	}
}
