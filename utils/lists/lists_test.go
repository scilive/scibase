package lists_test

import (
	"strings"
	"testing"

	"github.com/scilive/scibase/utils/lists"
	"github.com/stretchr/testify/assert"
)

type user struct {
	Id   int
	Name string
}

func TestSort(t *testing.T) {
	us := []user{
		{Id: 1, Name: "a"},
		{Id: 2, Name: "b2"},
		{Id: 2, Name: "b1"},
		{Id: 3, Name: "c"},
	}
	lists.Sort(us, func(i, j int) []int {
		idCmp := us[i].Id - us[j].Id
		nameCmp := strings.Compare(us[i].Name, us[j].Name)
		return []int{idCmp, nameCmp}
	})
	assert.Equal(t, us[1].Name, "b1")
	assert.Equal(t, us[2].Name, "b2")
}

func TestConcat(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	c := []int{7, 8, 9}
	r := lists.Concat(a, b, c)
	assert.Equal(t, r, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestUnique(t *testing.T) {
	a := []int{1, 2, 3, 3, 4, 5, 6, 7, 8, 9}
	r := lists.Unique(a)
	assert.Equal(t, r, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestMap(t *testing.T) {
	a := []int{1, 2, 3}
	r := lists.Map(a, func(v int, i int) int {
		return v + i
	})
	assert.Equal(t, r, []int{1, 3, 5})
}

func TestReduce(t *testing.T) {
	a := []int{1, 2, 3}
	r := lists.Reduce(a, func(v int, i int, acc int) int {
		return acc + v + i
	}, 0)
	assert.Equal(t, r, 9)
}
func TestFilter(t *testing.T) {
	a := []int{1, 2, 3}
	r := lists.Filter(a, func(v int, i int) bool {
		return v > 1
	})
	assert.Equal(t, r, []int{2, 3})
}

func TestIndex(t *testing.T) {
	a := []int{1, 2, 3}
	r := lists.Index(a, 2)
	assert.Equal(t, r, 1)
}

func TestIndexBy(t *testing.T) {
	a := []int{1, 2, 3}
	r := lists.IndexBy(a, func(v int, i int) bool {
		return v == 2
	})
	assert.Equal(t, r, 1)
}

func TestContains(t *testing.T) {
	a := []int{1, 2, 3}
	r := lists.Contains(a, 2)
	assert.Equal(t, r, true)
}

func TestStrs2Int64s(t *testing.T) {
	a := []string{"1", "2", "3"}
	r, err := lists.Strs2Int64s(a)
	assert.Nil(t, err)
	assert.Equal(t, r, []int64{1, 2, 3})
}

func TestInt64s2Strs(t *testing.T) {
	a := []int64{1, 2, 3}
	r := lists.Int64s2Strs(a)
	assert.Equal(t, r, []string{"1", "2", "3"})
}
