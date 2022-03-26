package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBTree(t *testing.T) {
	tree := NewBTree[InternalKey, string](32)

	key1 := NewInternalKey("hello", 1)
	value1 := "world"

	key2 := NewInternalKey("hello", 2)
	value2 := "not overwrite!"
	tree.Insert(key1, value1)
	tree.Insert(key2, value2)

	key, value := tree.Find(key1)
	assert.Equal(t, key, key1)
	assert.Equal(t, *value, value1)

	key, value = tree.Find(key2)
	assert.Equal(t, key, key2)
	assert.Equal(t, *value, value2)

	key3 := NewInternalKey("hello", 3)
	key, value = tree.Find(key3)
	assert.Equal(t, key, key2)
	assert.Equal(t, *value, value2)
}
