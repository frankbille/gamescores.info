package utils

import (
	"sort"
)

// Int64Set is a set of int64 types, guaranteeing uniqueness
type Int64Set struct {
	set map[int64]bool
}

// NewInt64Set creates a new instance of the set
func NewInt64Set() *Int64Set {
	return &Int64Set{make(map[int64]bool)}
}

// Add a value to the set
func (set *Int64Set) Add(i int64) bool {
	_, found := set.set[i]
	set.set[i] = true
	return !found
}

// Values as an array
func (set *Int64Set) Values() []int64 {
	v := make([]int64, len(set.set))
	idx := 0
	for i := range set.set {
		v[idx] = i
		idx++
	}
	sort.Sort(int64arr(v))
	return v
}

type int64arr []int64

func (a int64arr) Len() int           { return len(a) }
func (a int64arr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }
