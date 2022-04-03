package db

import (
	"sync"

	. "github.com/yah01/CyberKV/common"
)

const (
	LeafNodeType = iota
	InternalNodeType
)

type NodeType int

type BTreeNode[K Comparable[K], V any] interface {
	Parent() *InternalNode[K, V]
	SetParent(*InternalNode[K, V])
	Split() (BTreeNode[K, V], K)
	Len() int
	Cap() int
	NodeType() NodeType
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

type baseBTreeNode[K Comparable[K], V any] struct {
	rwmutex sync.RWMutex
	parent  *InternalNode[K, V]
	keys    []K
}

func newBaseBTreeNode[K Comparable[K], V any](fanout int) *baseBTreeNode[K, V] {
	return &baseBTreeNode[K, V]{
		rwmutex: sync.RWMutex{},
		parent:  nil,
		keys:    make([]K, 0, fanout-1),
	}
}

func (node *baseBTreeNode[K, V]) Parent() *InternalNode[K, V] {
	return node.parent
}

func (node *baseBTreeNode[K, V]) SetParent(parent *InternalNode[K, V]) {
	node.parent = parent
}

func (node *baseBTreeNode[K, V]) Len() int {
	return len(node.keys)
}

func (node *baseBTreeNode[K, V]) Cap() int {
	return cap(node.keys)
}

func (node *baseBTreeNode[K, V]) Lock() {
	node.rwmutex.Lock()
}

func (node *baseBTreeNode[K, V]) Unlock() {
	node.rwmutex.Unlock()
}

func (node *baseBTreeNode[K, V]) RLock() {
	node.rwmutex.RLock()
}

func (node *baseBTreeNode[K, V]) RUnlock() {
	node.rwmutex.RUnlock()
}

type LeafNode[K Comparable[K], V any] struct {
	*baseBTreeNode[K, V]
	values []V
	right  *LeafNode[K, V]
}

func NewLeafNode[K Comparable[K], V any](fanout int) *LeafNode[K, V] {
	return &LeafNode[K, V]{
		baseBTreeNode: newBaseBTreeNode[K, V](fanout),
		values:        make([]V, 0, fanout-1),
		right:         nil,
	}
}

func (node *LeafNode[K, V]) Insert(key K, value V) {
	if len(node.keys) == 0 {
		node.keys = append(node.keys, key)
		node.values = append(node.values, value)
		return
	}

	if node.Len() == node.Cap() {
		panic("insert into a full leaf node")
	}

	i := Search(node.keys, key)
	if i < node.Len() && node.keys[i].Compare(key) == 0 {
		node.values[i] = value
		return
	}
	node.keys = append(node.keys[:i+1], node.keys[i:]...)
	node.keys[i] = key

	node.values = append(node.values[:i+1], node.values[i:]...)
	node.values[i] = value
}

func (node *LeafNode[K, V]) Get(key K) (value V) {
	if node.Len() == 0 {
		return
	}

	i := Search(node.keys, key)
	if node.keys[i].Compare(key) != 0 {
		return
	}

	return node.values[i]
}

func (node *LeafNode[K, V]) NodeType() NodeType {
	return LeafNodeType
}

func (node *LeafNode[K, V]) Split() (BTreeNode[K, V], K) {
	mid := node.Len()/2 + 1
	parent := node.parent

	newLeaf := NewLeafNode[K, V](node.Cap() + 1)
	newLeaf.keys = append(newLeaf.keys, node.keys[mid:]...)
	newLeaf.values = append(newLeaf.values, node.values[mid:]...)
	newLeaf.SetParent(parent)
	newLeaf.right = node.right

	node.keys = node.keys[:mid]
	node.values = node.values[:mid]
	node.right = newLeaf

	parent.Insert(newLeaf.keys[0], newLeaf)

	return newLeaf, newLeaf.keys[0]
}

type InternalNode[K Comparable[K], V any] struct {
	*baseBTreeNode[K, V]
	children []BTreeNode[K, V]
}

func NewInternalNode[K Comparable[K], V any](fanout int, initChlid ...BTreeNode[K, V]) *InternalNode[K, V] {
	children := make([]BTreeNode[K, V], 0, fanout)
	children = append(children, initChlid...)

	node := &InternalNode[K, V]{
		baseBTreeNode: newBaseBTreeNode[K, V](fanout),
		children:      children,
	}

	for i := range children {
		children[i].SetParent(node)
	}

	return node
}

func (node *InternalNode[K, V]) Insert(key K, child BTreeNode[K, V]) {
	if len(node.children) == 0 {
		panic("InternalNode's children can't be empty")
	}

	if node.Len() == node.Cap() {
		panic("insert into a full leaf node")
	}

	i := Search(node.keys, key)
	node.keys = append(node.keys[:i+1], node.keys[i:]...)
	node.keys[i] = key

	i++
	node.children = append(node.children[:i+1], node.children[i:]...)
	node.children[i] = child
}

func (node *InternalNode[K, V]) Split() (BTreeNode[K, V], K) {
	mid := node.Len()/2 + 1
	parent := node.parent

	newNode := NewInternalNode[K, V](node.Cap() + 1)
	newNode.keys = append(newNode.keys, node.keys[mid+1:]...)
	newNode.children = append(newNode.children, node.children[mid+1:]...)
	newNode.SetParent(parent)
	for _, child := range node.children[mid+1:] {
		child.SetParent(newNode)
	}

	midKey := node.keys[mid]
	node.keys = node.keys[:mid]
	node.children = node.children[:mid+1]

	parent.Insert(midKey, newNode)

	return newNode, midKey
}

func (node *InternalNode[K, V]) NodeType() NodeType {
	return InternalNodeType
}

func (node *InternalNode[K, V]) findChild(key K) BTreeNode[K, V] {
	i := Search(node.keys, key)
	return node.children[i]
}

type BTree[K Comparable[K], V any] struct {
	rwmutex  sync.RWMutex // guard root
	root     BTreeNode[K, V]
	leftmost *LeafNode[K, V]
	height   int
}

func NewBTree[K Comparable[K], V any](fanout int) *BTree[K, V] {
	root := NewLeafNode[K, V](fanout)
	return &BTree[K, V]{
		rwmutex:  sync.RWMutex{},
		root:     root,
		leftmost: root,
		height:   1,
	}
}

func (tree *BTree[K, V]) Get(key K) V {
	leaf := tree.goToLeaf(tree.getRoot(), key)
	leaf.RLock()
	defer leaf.RUnlock()

	return leaf.Get(key)
}

// Find the first element with key not less than the given key
func (tree *BTree[K, V]) Find(key K) (K, V, bool) {
	var (
		defaultKey   K
		defaultValue V
	)
	leaf := tree.goToLeaf(tree.getRoot(), key)
	leaf.RLock()
	defer leaf.RUnlock()
	if leaf.Len() == 0 {
		return defaultKey, defaultValue, false
	}

	i := Search(leaf.keys, key)
	if i >= leaf.Len() {
		return defaultKey, defaultValue, false
	}

	return leaf.keys[i], leaf.values[i], true
}

func (tree *BTree[K, V]) Insert(key K, value V) {
	var root = tree.getRoot()
	path := tree.leafPath(root, key)
	leaf := path[len(path)-1].(*LeafNode[K, V])

	if leaf.Len() < leaf.Cap() {
		leaf.Insert(key, value)
		tree.unlockPath(path)
		return
	}

	ancientIndex := len(path) - 1
	for ancientIndex >= 0 && path[ancientIndex].Len() == path[ancientIndex].Cap() {
		ancientIndex--
	}

	var parent *InternalNode[K, V]
	// The tree is full
	if ancientIndex == -1 {
		parent = NewInternalNode(tree.Fanout(), root)
		parent.Lock()
		root.SetParent(parent)
		tree.setRoot(parent)
	} else {
		parent = path[ancientIndex].(*InternalNode[K, V])
	}

	if ancientIndex > 0 {
		tree.unlockPath(path[:ancientIndex])
	}

	// Split chian
	tree.split(leaf)

	tree.unlockPath(path[ancientIndex+1:])
	parent.Unlock()

	tree.Insert(key, value)
}

func (tree *BTree[K, V]) split(node BTreeNode[K, V]) {
	if node.Parent().Len() == node.Parent().Cap() {
		tree.split(node.Parent())
	}

	node.Split()
}

func (tree *BTree[K, V]) Fanout() int {
	return tree.root.Cap() + 1
}

func (tree *BTree[K, V]) Range(fn func(key K, value V) bool) {
	// leafNode := tree.root
	// for leafNode.NodeType() != LeafNodeType {
	// 	leafNode = leafNode.(*InternalNode[K, V]).children[0]
	// }
	// leaf := leafNode.(*LeafNode[K, V])
	leaf := tree.leftmost
	for leaf != nil {
		leaf.RLock()

		for i := range leaf.keys {
			if !fn(leaf.keys[i], leaf.values[i]) {
				return
			}
		}
		right := leaf.right
		leaf.RUnlock()

		leaf = right
	}
}

func (tree *BTree[K, V]) leafPath(node BTreeNode[K, V], key K) []BTreeNode[K, V] {
	node.Lock()
	nodes := make([]BTreeNode[K, V], 0, tree.height)

	for node.NodeType() != LeafNodeType {
		nodes = append(nodes, node)
		node = node.(*InternalNode[K, V]).findChild(key)
		node.Lock()
	}

	nodes = append(nodes, node)

	return nodes
}

func (tree *BTree[K, V]) goToLeaf(node BTreeNode[K, V], key K) *LeafNode[K, V] {
	node.RLock()
	for node.NodeType() != LeafNodeType {
		child := node.(*InternalNode[K, V]).findChild(key)
		child.RLock()
		node.RUnlock()
		node = child
	}

	node.RUnlock()

	return node.(*LeafNode[K, V])
}

func (tree *BTree[K, V]) unlockPath(path []BTreeNode[K, V]) {
	for i := range path {
		path[len(path)-i-1].Unlock()
	}
}

func (tree *BTree[K, V]) getRoot() BTreeNode[K, V] {
	tree.rwmutex.RLock()
	defer tree.rwmutex.RUnlock()

	return tree.root
}

func (tree *BTree[K, V]) setRoot(node BTreeNode[K, V]) {
	tree.rwmutex.Lock()
	defer tree.rwmutex.Unlock()

	tree.root = node
}
