package maths

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](values ...T) T {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func Max[T constraints.Ordered](values ...T) T {
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
