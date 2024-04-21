package sqlengine

import (
	"testing"
	"unsafe"
)

type C struct {
	tree  BTree
	ref   map[string]string
	pages map[uint64]BNode
}

func newC() *C {
	pages := map[uint64]BNode{}
	return &C{
		tree: BTree{
			get: func(ptr uint64) BNode {
				node, ok := pages[ptr]
				assert(ok)
				return node
			},
			new: func(node BNode) uint64 {
				assert(node.nbytes() <= BTREE_PAGE_SIZE)
				key := uint64(uintptr(unsafe.Pointer(&node.data[0])))
				assert(pages[key].data == nil)
				pages[key] = node
				return key
			},
			del: func(ptr uint64) {
				_, ok := pages[ptr]
				assert(ok)
				delete(pages, ptr)
			},
		},
		ref:   map[string]string{},
		pages: pages,
	}
}

func (c *C) add(key string, val string) {
	c.tree.Insert([]byte(key), []byte(val))
	c.ref[key] = val
}

func (c *C) del(key string) bool {
	delete(c.ref, key)
	return c.tree.Delete([]byte(key))
}

func TestAdd(t *testing.T) {
	c := newC()
	key := "key"
	val := "val"
	c.add(key, val)
	c.add("key1", "val1")
	if c.ref[key] != val {
		t.Errorf("Expected value: %s, but got: %s", val, c.ref[key])
	}
	if c.del(key) != true {
		t.Errorf("Expected value: %s, but got: %s", val, c.ref[key])
	}
	if c.del(key) != false {
		t.Errorf("Expected value: %s, but got: %s", val, c.ref[key])
	}

}
