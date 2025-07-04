package sqlengine

import (
	"bytes"
	"encoding/binary"
)

// BNode /*
/**
A node consists of:
1. A fixed-sized header containing the type of the node (leaf node or internal node) and the number of keys.
2. A list of pointers to the child nodes. (Used by internal nodes).
3. A list of offsets pointing to each key-value pair.
4. Packed KV pairs.
| type | nkeys | pointers | offsets | key-values
| 2B   | 2B    |nkeys*8B   |nkeys*2B|...
This is the format of the KV pair. Lengths followed by data.
| klen | vlen | key | val |
| 2B | 2B |...|...|
*/
type BNode struct {
	data []byte //can be dumped to the disk
}

const (
	BNODE_NODE         = 1 //internal nodes without value
	BNODE_LEAF         = 2 // leaf nodes with values
	HEADER             = 4
	BTREE_PAGE_SIZE    = 4096 // 4K bytes one page
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
)

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	assert(node1max <= BTREE_PAGE_SIZE)
}

// header
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

// 有多少个key
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

// pointers
func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[pos:])
}
func (node BNode) setPtr(idx uint16, val uint64) {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

// offset list
func offsetPos(node BNode, idx uint16) uint16 {
	assert(1 <= idx && idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx):])
}
func (node BNode) setOffset(idx uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx):], offset)
}

// key-values
func (node BNode) kvPos(idx uint16) uint16 {
	assert(idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos:])
	return node.data[pos+4:][:klen]
}
func (node BNode) getVal(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos+0:])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])
	// 4代表4个存储key长度和value长度的字节数
	return node.data[pos+4+klen:][:vlen]
}

// node size in bytes: node的总大小
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

// look up the key

func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)
	// the first key is a copy from the parent node
	// thus it's always less than or equal to the key
	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}
	return found
}

// add leaf node
func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	// update header
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

// The nodeAppendKV function copies a KV pair to the new node.
func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	// ptrs
	new.setPtr(idx, ptr)
	// KVs
	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new.data[pos+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new.data[pos+2:], uint16(len(val)))
	copy(new.data[pos+4:], key)
	copy(new.data[pos+4+uint16(len(key)):], val)
	// the offset of the next key
	new.setOffset(idx+1, new.getOffset(idx)+4+uint16(len(key)+len(val)))
}

// The nodeAppendRange function copies keys from an old node to a new node.
func nodeAppendRange(new BNode, old BNode, dstNew uint16, srcOld uint16, n uint16) {
	assert(srcOld+n <= old.nkeys())
	assert(dstNew+n <= new.nkeys())
	if n == 0 {
		return
	}
	// pointers
	for i := uint16(0); i < n; i++ {
		new.setPtr(dstNew+i, old.getPtr(srcOld+i))
	}
	// offsets
	dstBegin := new.getOffset(dstNew)
	srcBegin := old.getOffset(srcOld)
	for i := uint16(1); i < n; i++ { // NOTE: the range is [1,n]
		offset := dstBegin + old.getOffset(srcOld+i) - srcBegin
		new.setOffset(dstNew+i, offset)
	}
	// KVs
	begin := old.kvPos(srcOld)
	end := old.kvPos(srcOld + n)
	copy(new.data[new.kvPos(dstNew):], old.data[begin:end])
}

func treeInsert(tree *BTree, node BNode, key []byte, val []byte) BNode {
	// the result node it's allowed to be bigger than 1 page and will be split if so
	newNode := BNode{data: make([]byte, 2*BTREE_PAGE_SIZE)}
	//where to insert the key
	idx := nodeLookupLE(node, key)
	// act depending on the node type
	switch node.btype() {
	case BNODE_LEAF:
		// leaf, node.getKey(idx)<=key
		if bytes.Equal(key, node.getKey(idx)) {
			// found the key, update it
			leafUpdate(newNode, node, idx, key, val)
		} else {
			// insert it after the position
			leafInsert(newNode, node, idx, key, val)
		}

	case BNODE_NODE:
		// internal node, insert it to kid node
		nodeInsert(tree, newNode, node, idx, key, val)
	default:
		panic("bad node!")
	}
	return newNode
}

func leafUpdate(node BNode, node2 BNode, idx uint16, key []byte, val []byte) {

}

// KV insertion to a internal node
func nodeInsert(tree *BTree, new BNode, node BNode, idx uint16, key []byte, val []byte) {
	// get and deallocate the kid node
	kptr := node.getPtr(idx)
	knode := tree.get(kptr)
	tree.del(kptr)
	// recursive insertion to the kid node
	knode = treeInsert(tree, knode, key, val)
	// split the result
	nsplit, splited := nodeSplit3(knode)
	// update the kid links
	nodeReplaceKidN(tree, new, node, idx, splited[:nsplit]...)

}

// replace a link with multiple links
func nodeReplaceKidN(tree *BTree, new BNode, old BNode, idx uint16, kids ...BNode) {
	inc := uint16(len(kids))
	new.setHeader(BNODE_NODE, old.nkeys()+inc-1)
	nodeAppendRange(new, old, 0, 0, idx)
	for i, node := range kids {
		nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
	}
	nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys()-(idx+1))

}

// split a node if it's too big, the results are 1-3 nodes.
func nodeSplit3(old BNode) (nodeCount uint16, nodes [3]BNode) {
	if old.nbytes() <= BTREE_PAGE_SIZE {
		old.data = old.data[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old}
	}
	left := BNode{make([]byte, 2*BTREE_PAGE_SIZE)} // might be split later
	right := BNode{make([]byte, BTREE_PAGE_SIZE)}
	nodeSplit2(left, right, old)
	if left.nbytes() <= BTREE_PAGE_SIZE {
		left.data = left.data[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}
	}
	// the left node is still too large
	leftleft := BNode{make([]byte, BTREE_PAGE_SIZE)}
	middle := BNode{make([]byte, BTREE_PAGE_SIZE)}
	nodeSplit2(leftleft, middle, left)
	assert(leftleft.nbytes() <= BTREE_PAGE_SIZE)
	return 3, [3]BNode{leftleft, middle, right}
}

func nodeSplit2(left BNode, right BNode, old BNode) {

}

// remove a key from a leaf node
func leafDelete(new BNode, old BNode, idx uint16) {
	new.setHeader(BNODE_LEAF, old.nkeys()-1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendRange(new, old, idx, idx+1, old.nkeys()-(idx+1))
}

// delete a key from tree, code is similar to the insertion
func treeDelete(tree *BTree, node BNode, key []byte) BNode {
	// where to find the key?
	idx := nodeLookupLE(node, key)
	// act depending on the node type
	switch node.btype() {
	case BNODE_LEAF:
		if !bytes.Equal(key, node.getKey(idx)) {
			return BNode{} // not found
		}
		// delete the key in the leaf
		nNode := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		leafDelete(nNode, node, idx)
		return nNode

	case BNODE_NODE:
		return nodeDelete(tree, node, idx, key)
	default:
		panic("bad node")
	}
}

// handle internal nodes
// The difference is that we need to merge nodes instead of splitting nodes.
// Part of treeDelete()
func nodeDelete(tree *BTree, node BNode, idx uint16, key []byte) BNode {
	// recurse into the kid
	kPtr := node.getPtr(idx)
	updated := treeDelete(tree, tree.get(kPtr), key)
	if len(updated.data) == 0 {
		return BNode{}
	}
	tree.del(kPtr)
	newNode := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
	// check for merging
	mergeDir, sibling := shouldMerge(tree, node, idx, updated)
	switch {
	case mergeDir < 0: //left
		merged := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		nodeMerge(merged, sibling, updated)
		tree.del(node.getPtr(idx - 1))
		nodeReplace2Kid(newNode, node, idx, tree.new(merged), merged.getKey(0))
	case mergeDir > 0: // right
		merged := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		nodeMerge(merged, updated, sibling)
		tree.del(node.getPtr(idx + 1))
		nodeReplace2Kid(newNode, node, idx, tree.new(merged), merged.getKey(0))
	case mergeDir == 0:
		assert(updated.nkeys() > 0)
		nodeReplaceKidN(tree, newNode, node, idx, updated)
	}
	return newNode

}

func nodeReplace2Kid(new BNode, node BNode, idx uint16, u uint64, key []byte) {

}

// merge 2 nodes into 1
func nodeMerge(newNode BNode, leftNode BNode, rightNode BNode) {
	newNode.setHeader(leftNode.btype(), leftNode.nkeys()+rightNode.nkeys())
	nodeAppendRange(newNode, leftNode, 0, 0, leftNode.nkeys())
	nodeAppendRange(newNode, rightNode, leftNode.nkeys(), 0, rightNode.nkeys())

}

// should the updated kid be merged with a sibling
func shouldMerge(tree *BTree, node BNode, idx uint16, updated BNode) (int, BNode) {
	if updated.nkeys() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}
	}
	if idx > 0 {
		sibling := tree.get(node.getPtr(idx - 1))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return -1, sibling
		}
	}
	if idx+1 < node.nkeys() {
		sibling := tree.get(node.getPtr(idx + 1))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return +1, sibling
		}
	}

	return 0, BNode{}

}
