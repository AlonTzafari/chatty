package main

func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findIndex[T comparable](s []T, item T) int {
	for i, v := range s {
		if v == item {
			return i
		}
	}
	return -1
}
