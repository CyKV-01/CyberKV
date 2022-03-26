package db

import (
	. "github.com/yah01/CyberKV/common"
)

const (
	LeafNodeType = iota
	InternalNodeType
)

type NodeType int

type BTreeNode[K Comparable[K], V any] interface {
	Split(parent *InternalNode[K, V]) (BTreeNode[K, V], K)
	InnerKeys() []K
	Len() int
	Cap() int
	NodeType() NodeType
}

type baseBTreeNode[K Comparable[K]] struct {
	keys []K
}

func (node *baseBTreeNode[K]) InnerKeys() []K {
	return node.keys
}

func (node *baseBTreeNode[K]) Len() int {
	return len(node.keys)
}

func (node *baseBTreeNode[K]) Cap() int {
	return cap(node.keys)
}

type LeafNode[K Comparable[K], V any] struct {
	*baseBTreeNode[K]
	values []V
	right  *LeafNode[K, V]
}

func NewLeafNode[K Comparable[K], V any](fanout int) *LeafNode[K, V] {
	return &LeafNode[K, V]{
		baseBTreeNode: &baseBTreeNode[K]{make([]K, 0, fanout-1)},
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

func (node *LeafNode[K, V]) Get(key K) *V {
	if node.Len() == 0 {
		return nil
	}

	i := Search(node.keys, key)
	if node.keys[i].Compare(key) != 0 {
		return nil
	}

	return &node.values[i]
}

func (node *LeafNode[K, V]) NodeType() NodeType {
	return LeafNodeType
}

func (node *LeafNode[K, V]) Split(parent *InternalNode[K, V]) (BTreeNode[K, V], K) {
	mid := node.Len()/2 + 1

	newLeaf := NewLeafNode[K, V](node.Cap() + 1)
	newLeaf.keys = append(newLeaf.keys, node.keys[mid:]...)
	newLeaf.values = append(newLeaf.values, node.values[mid:]...)

	node.keys = node.keys[:mid]
	node.values = node.values[:mid]
	node.right = newLeaf

	parent.Insert(newLeaf.keys[0], newLeaf)

	return newLeaf, newLeaf.keys[0]
}

type InternalNode[K Comparable[K], V any] struct {
	*baseBTreeNode[K]
	children []BTreeNode[K, V]
}

func NewInternalNode[K Comparable[K], V any](fanout int, initChlid ...BTreeNode[K, V]) *InternalNode[K, V] {
	children := make([]BTreeNode[K, V], 1, fanout)
	children = append(children, initChlid...)

	return &InternalNode[K, V]{
		baseBTreeNode: &baseBTreeNode[K]{make([]K, 0, fanout-1)},
		children:      children,
	}
}

func (node *InternalNode[K, V]) Insert(key K, child BTreeNode[K, V]) {

}

func (node *InternalNode[K, V]) Split(parent *InternalNode[K, V]) (BTreeNode[K, V], K) {
	mid := node.Len() / 2

	newNode := NewInternalNode[K, V](node.Cap() + 1)
	newNode.keys = append(newNode.keys, node.keys[mid+1:]...)
	newNode.children = append(newNode.children, node.children[mid+1:]...)

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
	root   BTreeNode[K, V]
	height int
}

func NewBTree[K Comparable[K], V any](fanout int) *BTree[K, V] {
	return &BTree[K, V]{
		root:   NewLeafNode[K, V](fanout),
		height: 1,
	}
}

func (tree *BTree[K, V]) Get(key K) *V {
	leaf := tree.GoToLeaf(tree.root, key)
	leaf.Get(key)

	return leaf.Get(key)
}

// Find the first element with key not less than the given key
func (tree *BTree[K, V]) Find(key K) (K, *V) {
	var resultKey K
	leaf := tree.GoToLeaf(tree.root, key)
	if leaf.Len() == 0 {
		return resultKey, nil
	}

	i := Search(leaf.keys, key)
	if i >= leaf.Len() {
		return resultKey, nil
	}

	return leaf.keys[i], &leaf.values[i]
}

func (tree *BTree[K, V]) Insert(key K, value V) {
	path := tree.leafPath(tree.root, key)
	leaf := path[len(path)-1].(*LeafNode[K, V])

	if leaf.Len() < leaf.Cap() {
		leaf.Insert(key, value)
		return
	}

	i := len(path) - 1
	for i >= 0 && path[i].Len() == path[i].Cap() {
		i--
	}

	var parent *InternalNode[K, V]
	// The tree is full
	if i == -1 {
		parent = NewInternalNode(tree.Fanout(), tree.root)
	} else {
		parent = path[i].(*InternalNode[K, V])
	}

	// Split chian
	for i++; i < len(path); i++ {
		node := path[i]

		newNode, midKey := node.Split(parent)

		if node.NodeType() == InternalNodeType {
			if key.Compare(midKey) < 0 {
				parent = node.(*InternalNode[K, V])
			} else {
				parent = newNode.(*InternalNode[K, V])
			}
		} else {
			if key.Compare(midKey) < 0 {
				node.(*LeafNode[K, V]).Insert(key, value)
			} else {
				newNode.(*LeafNode[K, V]).Insert(key, value)
			}
		}
	}
}

func (tree *BTree[K, V]) Fanout() int {
	return tree.root.Cap() + 1
}

func (tree *BTree[K, V]) leafPath(node BTreeNode[K, V], key K) []BTreeNode[K, V] {
	nodes := make([]BTreeNode[K, V], 0, tree.height)

	for ; node.NodeType() != LeafNodeType; node = node.(*InternalNode[K, V]).findChild(key) {
		nodes = append(nodes, node)
	}

	nodes = append(nodes, node)

	return nodes
}

func (tree *BTree[K, V]) GoToLeaf(node BTreeNode[K, V], key K) *LeafNode[K, V] {
	for node.NodeType() != LeafNodeType {
		node = node.(*InternalNode[K, V]).findChild(key)
	}

	return node.(*LeafNode[K, V])
}
