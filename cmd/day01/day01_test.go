package day01

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseElfCalorieList_EmptyString(t *testing.T) {
	list, err := ParseElfCalorieList("")
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}

func TestParseElfCalorieList_Newlines(t *testing.T) {
	list, err := ParseElfCalorieList("\n\n\n")
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}

func TestParseElfCalorieList_OneElf_OneItem(t *testing.T) {
	list, err := ParseElfCalorieList("1000")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, list[0], 1000)
}

func TestParseElfCalorieList_OneElf_UnexpectedItem(t *testing.T) {
	list, err := ParseElfCalorieList("asdf")
	assert.Error(t, err)
	assert.Len(t, list, 0)
}

func TestParseElfCalorieList_OneElf_MultipleItemsOnLine(t *testing.T) {
	list, err := ParseElfCalorieList("1000 2000")
	assert.Error(t, err)
	assert.Len(t, list, 0)
}

func TestParseElfCalorieList_OneElf_TwoItems(t *testing.T) {
	list, err := ParseElfCalorieList("1000\n3000")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, list[0], 4000)
}

func TestParseElfCalorieList_TwoElves_TwoItems(t *testing.T) {
	list, err := ParseElfCalorieList("1000\n3000\n\n4000\n5000")
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, list[0], 4000)
	assert.Equal(t, list[1], 9000)
}
