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

func TestLists(t *testing.T) {
	us := []user{
		{Id: 1, Name: "a"},
		{Id: 2, Name: "b"},
		{Id: 3, Name: "c"},
	}
	col := lists.Column(us, func(v user) int { return v.Id })
	assert.Equal(t, col, []int{1, 2, 3})
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
