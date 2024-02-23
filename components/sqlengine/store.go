package sqlengine

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

func mmapInit(fp *os.File) (size int, chunks []byte, err error) {
	file, err := fp.Stat()
	if err != nil {
		return 0, nil, fmt.Errorf("stat: %w", err)
	}
	if file.Size()%BTREE_PAGE_SIZE != 0 {
		return 0, nil, errors.New("fileSize size is not a multiple of page size")
	}
	// 64MB
	mmapSize := 64 << 20
	for mmapSize < int(file.Size()) {
		mmapSize = mmapSize * 2
	}

	chunk, err := syscall.Mmap(
		int(fp.Fd()),
		0,
		mmapSize,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("mmap: %w", err)
	}
	return int(file.Size()), chunk, nil
}

// extend the mmap by adding new mappings
// size increase exponentially,2^n
func extendMmap(db *KV, npages int) error {
	if db.mmap.total >= npages*BTREE_PAGE_SIZE {
		return nil // no need to extend
	}
	// double the address space
	chunk, err := syscall.Mmap(
		int(db.fp.Fd()),
		int64(db.mmap.total),
		db.mmap.total,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		return fmt.Errorf("mmap: %w", err)
	}
	db.mmap.total += db.mmap.total
	db.mmap.chunks = append(db.mmap.chunks, chunk)
	return nil
}

// callback for BTree, dereference a pointer.
func (db *KV) pageGet(ptr uint64) BNode {
	start := uint64(0)
	for _, chunk := range db.mmap.chunks {
		end := start + uint64(len(chunk))/BTREE_PAGE_SIZE
		if ptr < end {
			offset := BTREE_PAGE_SIZE * (ptr - start)
			return BNode{chunk[offset : offset+BTREE_PAGE_SIZE]}
		}
		start = end
	}
	panic("bad ptr")
}

type KV struct {
	Path string
	// internals
	fp   *os.File
	tree BTree
	mmap struct {
		fileSize int      // fileSize size, can be lager than the database
		total    int      //mmap size, can be larger than the fileSize size
		chunks   [][]byte // multiple mmaps, can non-continuous
	}
	page struct {
		flushed uint64   // database size in number of pages
		temp    [][]byte // newly allocated pages
	}
}

const DbSig = "jimmy-sql-123456"

// the master page format
// it contains the pointer to root and other important bits.
// |sig| btree_root| page_used|
// |16B| 8B        | 8B       |
func masterPageLoad(db *KV) error {
	if db.mmap.fileSize == 0 {
		// empty fileSize, the master page will be created on the first write.
		db.page.flushed = 1 // reserved for the master page
		return nil
	}
	data := db.mmap.chunks[0]
	root := binary.LittleEndian.Uint64(data[16:])
	used := binary.LittleEndian.Uint64(data[24:])
	// verify the page
	if !bytes.Equal([]byte(DbSig), data[:16]) {
		return errors.New("bad signature")
	}
	bad := !(1 <= used && used <= uint64(db.mmap.fileSize/BTREE_PAGE_SIZE))
	bad = bad || !(0 <= root && root <= used)
	if bad {
		return errors.New("bad master page")
	}
	db.tree.root = root
	db.page.flushed = used
	return nil
}

func masterStore(db *KV) error {
	var data [32]byte
	copy(data[:16], DbSig)
	binary.LittleEndian.PutUint64(data[16:], db.tree.root)
	binary.LittleEndian.PutUint64(data[24:], db.page.flushed)
	// NOTE: updating the page via mmap is not atomic.
	// use the `pwrite()` syscall instead.
	_, err := db.fp.WriteAt(data[:], 0)
	if err != nil {
		return fmt.Errorf("write master page: %w", err)
	}
	return nil
}

// callback for BTree, allocate new page
func (db *KV) pageNew(node BNode) uint64 {
	// TODO: reuse deallocated pages
	assert(len(node.data) <= BTREE_PAGE_SIZE)
	ptr := db.page.flushed + uint64(len(db.page.temp))
	db.page.temp = append(db.page.temp, node.data)
	return ptr
}

// callback for BTree , deallocate a page
func (db *KV) pageDel(uint642 uint64) {
	// todo:
}

// extend the file to at least `npages`
func extendFile(db *KV, npages int) error {
	filePages := db.mmap.fileSize / BTREE_PAGE_SIZE
	if filePages >= npages {
		return nil
	}
	for filePages < npages {
		// the file size is increased exponentially,
		// so that we don't have to extend the file for every update.
		inc := filePages / 8
		if inc < 1 {
			inc = 1
		}
		filePages += inc
	}
	fileSize := filePages * BTREE_PAGE_SIZE
	// In Linux, use err := syscall.Fallocate(int(db.fp.Fd()), 0, 0, int64(fileSize))
	err := unix.Ftruncate(int(db.fp.Fd()), int64(fileSize))
	if err != nil {
		return fmt.Errorf("ftruncate: %w", err)
	}
	db.mmap.fileSize = fileSize
	return nil
}

// Open : initializing the database
func (db *KV) Open() error {
	// open or create the DB file
	fp, err := os.OpenFile(db.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("OpenFile:%w", err)
	}
	db.fp = fp
	// create the initial mmap
	size, chunk, err := mmapInit(db.fp)
	if err != nil {
		goto fail
	}
	db.mmap.fileSize = size
	db.mmap.total = len(chunk)
	db.mmap.chunks = [][]byte{chunk}

	// BTree callbacks
	db.tree.get = db.pageGet
	db.tree.new = db.pageNew
	db.tree.del = db.pageDel

	// read the master page
	err = masterPageLoad(db)
	if err != nil {
		goto fail
	}

	// done
	return nil
fail:
	db.Close()
	return fmt.Errorf("KV.Open:%w", err)
}

func (db *KV) Close() {
	for _, chunk := range db.mmap.chunks {
		err := syscall.Munmap(chunk)
		assert(err == nil)
	}
	_ = db.fp.Close()
}

// Get :read the db
func (db *KV) Get(key []byte) ([]byte, bool) {
	// todo:
	return nil, false
}

// Set :update the db
func (db *KV) Set(key []byte, val []byte) error {
	db.tree.Insert(key, val)
	return flushPages(db)
}
func (db *KV) Del(key []byte) (bool, error) {
	deleted := db.tree.Delete(key)
	return deleted, flushPages(db)
}

// persist the newly allocated pages after updates
func flushPages(db *KV) error {
	if err := writePages(db); err != nil {
		return err
	}
	return syncPages(db)
}

// It is split into two phases as mentioned earlier.
func writePages(db *KV) error {
	// extend the file & mmap if needed
	npages := int(db.page.flushed) + len(db.page.temp)
	if err := extendFile(db, npages); err != nil {
		return err
	}
	if err := extendMmap(db, npages); err != nil {
		return err
	}
	// copy data to the file
	for i, page := range db.page.temp {
		ptr := db.page.flushed + uint64(i)
		copy(db.pageGet(ptr).data, page)
	}
	return nil
}

// Persist to Disk And the fsync is in between and after them.

func syncPages(db *KV) error {
	// flush data to the disk. must be done before updating the master page.
	if err := db.fp.Sync(); err != nil {
		return fmt.Errorf("fsync: %w", err)
	}
	db.page.flushed += uint64(len(db.page.temp))
	db.page.temp = db.page.temp[:0]
	// update & flush the master page
	if err := masterStore(db); err != nil {
		return err
	}
	if err := db.fp.Sync(); err != nil {
		return fmt.Errorf("fsync: %w", err)
	}
	return nil
}
