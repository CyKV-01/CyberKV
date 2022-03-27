package db

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBTree(t *testing.T) {
	tree := NewBTree[InternalKey, string](4)

	key1 := NewInternalKey("hello", 1)
	value1 := "world"

	key2 := NewInternalKey("hello", 2)
	value2 := "not overwrite!"
	tree.Insert(key1, value1)
	tree.Insert(key2, value2)

	key, value, ok := tree.Find(key1)
	assert.True(t, ok)
	assert.Equal(t, key, key1)
	assert.Equal(t, value, value1)

	key, value, ok = tree.Find(key2)
	assert.True(t, ok)
	assert.Equal(t, key, key2)
	assert.Equal(t, value, value2)

	key3 := NewInternalKey("hello", 3)
	key, value, ok = tree.Find(key3)
	assert.True(t, ok)
	assert.Equal(t, key, key2)
	assert.Equal(t, value, value2)

	key4 := NewInternalKey("hello", 0)
	key, value, ok = tree.Find(key4)
	assert.False(t, ok)

	// Trigger split
	key3 = NewInternalKey("split0", 1)
	value3 := "ready!"

	key4 = NewInternalKey("split0", 1)
	value4 := "ready!"
	tree.Insert(key3, value3)
	tree.Insert(key4, value4)

	assert.Equal(t, value1, tree.Get(key1))
	assert.Equal(t, value2, tree.Get(key2))
	assert.Equal(t, value3, tree.Get(key3))
	assert.Equal(t, value4, tree.Get(key4))
}

func TestConcurrentBTree(t *testing.T) {
	concurrency := runtime.GOMAXPROCS(0)

	tree := NewBTree[InternalKey, string](8)
	commonKey := NewInternalKey("hello", 5)

	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for c := 0; c < concurrency; c++ {
		go func(idx int) {
			defer wg.Done()

			tree.Insert(commonKey, fmt.Sprint(idx))
			newKey := NewInternalKey(fmt.Sprint("new key from", idx), 10)
			tree.Insert(newKey, fmt.Sprint(idx))
		}(c)
	}

	wg.Wait()
}
