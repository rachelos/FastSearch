package btree

import (
	"github.com/google/btree"
)

const deg = 8 // B-Tree的度

type DataItem string

type DataIndex struct {
	tree *btree.BTree
}

func NewDataIndex() *DataIndex {
	return &DataIndex{
		tree: btree.New(deg),
	}
}

func (di *DataIndex) Insert(key string, value DataItem) {
	di.tree.ReplaceOrInsert(makeItem(key, value))
}

func (di *DataIndex) Get(key string) (DataItem, bool) {
	item := di.tree.Get(makeItem(key, ""))
	if item == nil {
		return "", false
	}
	return item.(Item).value, true
}

func makeItem(key string, value DataItem) Item {
	return Item{
		key:   key,
		value: value,
	}
}

type Item struct {
	key   string
	value DataItem
}

func (i Item) Less(than btree.Item) bool {
	return i.key < than.(Item).key
}
