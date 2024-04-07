package sqlengine

type BTree struct {
	// pointer (a nonzero page number)
	root uint64
	// callbacks for managing on-disk pages
	get func(uint64) BNode //dereference a pointer
	new func(BNode) uint64 // allocate a new page
	del func(uint64)       // deallocate a page
}

func (tree *BTree) Delete(key []byte) bool {
	assert(len(key) != 0)
	assert(len(key) <= BTREE_MAX_KEY_SIZE)
	if tree.root == 0 {
		return false
	}
	updated := treeDelete(tree, tree.get(tree.root), key)
	if len(updated.data) == 0 {
		return false //not found
	}
	tree.del(tree.root)
	if updated.btype() == BNODE_NODE && updated.nkeys() == 1 {
		// remove a level
		tree.root = updated.getPtr(0)
	} else {
		tree.root = tree.new(updated)

	}
	return true

}

// Insert  the interface for insertion
// It does two things:
// 1. A new root node is created when the old root is split into multiple nodes.
// 2. When inserting the first key, create the first leaf node as the root.
func (tree *BTree) Insert(key []byte, val []byte) {
	assert(len(key) != 0)
	assert(len(key) <= BTREE_MAX_KEY_SIZE)
	assert(len(val) <= BTREE_MAX_VAL_SIZE)
	if tree.root == 0 {
		// create the first node
		root := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		root.setHeader(BNODE_LEAF, 2)
		// a dummy key, this makes the tree cover the whole key space. // thus a lookup can always find a containing node.
		nodeAppendKV(root, 0, 0, nil, nil)
		nodeAppendKV(root, 1, 0, key, val)
		tree.root = tree.new(root)
		return
	}
	node := tree.get(tree.root)
	tree.del(tree.root)
	node = treeInsert(tree, node, key, val)
	nsplit, splitted := nodeSplit3(node)
	if nsplit > 1 {
		// the root was split, add a new level.
		root := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		root.setHeader(BNODE_NODE, nsplit)
		for i, knode := range splitted[:nsplit] {
			ptr, key := tree.new(knode), knode.getKey(0)
			nodeAppendKV(root, uint16(i), ptr, key, nil)
		}
		tree.root = tree.new(root)
	} else {
		// the root was not split, just replace it.
		tree.root = tree.new(splitted[0])
	}
}

func (tree *BTree) Get(key []byte) ([]byte, bool) {
	// todo:
	return nil, false
}
