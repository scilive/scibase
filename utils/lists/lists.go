package lists

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Unique returns a new slice containing the unique values in list.
func Unique[T comparable](list []T) []T {
	m := make(map[T]bool, len(list))
	var r []T
	for _, v := range list {
		if _, ok := m[v]; !ok {
			m[v] = true
			r = append(r, v)
		}
	}
	return r
}

// Concat arrays with copy
func Concat[T any](arrs ...[]T) []T {
	n := 0
	for _, arr := range arrs {
		n += len(arr)
	}
	res := make([]T, n)
	i := 0
	for _, arr := range arrs {
		copy(res[i:], arr)
		i += len(arr)
	}
	return res
}

// Map returns a new slice containing the results of applying fn to each
func Map[T any, V any](arr []T, fn func(v T, i int) V) []V {
	res := make([]V, len(arr))
	for i, v := range arr {
		res[i] = fn(v, i)
	}
	return res
}

// Reduce returns a new slice containing the results of applying fn to each
func Reduce[T any, V any](arr []T, fn func(v T, i int, acc V) V, acc V) V {
	for i, v := range arr {
		acc = fn(v, i, acc)
	}
	return acc
}

// Filter returns a new slice containing the results of applying fn to each, keep filter(v, i) == true
func Filter[T any](arr []T, filter func(v T, i int) bool) []T {
	var res []T
	for i, v := range arr {
		if filter(v, i) {
			res = append(res, v)
		}
	}
	return res
}

// FilterInt returns a new slice containing the results of applying fn to each
func FilterInt[T constraints.Integer](arr []T) []T {
	return Filter(arr, func(v T, i int) bool { return v != 0 })
}

// FilterStr returns a new slice containing the results of applying fn to each
func FilterStr(arr []string) []string {
	return Filter(arr, func(v string, i int) bool { return v != "" })
}

// Sort sorts the slice according to the fileds function.
// less: func(i, j int) , v[i]-v[j] < 0 meand {v[i], v[j]}
func Sort(arr any, less func(i, j int) []int) {
	sort.Slice(arr, func(i, j int) bool {
		result := less(i, j)
		for _, v := range result {
			if v == 0 {
				continue
			}
			return v < 0
		}
		return false
	})
}

func Index[T comparable](arr []T, value T) int {
	for i, v := range arr {
		if v == value {
			return i
		}
	}
	return -1
}

func IndexBy[T any](arr []T, fn func(v T, i int) bool) int {
	for i, v := range arr {
		if fn(v, i) {
			return i
		}
	}
	return -1
}
