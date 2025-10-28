// Package lru implements a simple LRU cache
// This implementation is non-thread safe
// thread safety is expected at application or buffer pool level
package lru

// Some page in our buffer pool
type Node struct {
	next *Node
	prev *Node
	val  string
	key  int
}

// Implements LRU cache algorithm
//
// - Hash map + Doubly Linked List
// - O(1) for all operations
type LruCache struct {
	head        *Node
	tail        *Node
	capacity    int
	index       map[int]*Node
	currentSize int
}

func (lc *LruCache) NewHead(node *Node, key int) {
	if node.next != nil {
		node.next.prev = node.prev // What if theres no previous node?
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node == lc.tail {
		lc.tail = node.prev
	}

	cachedHead := lc.head
	lc.head = node

	if cachedHead != nil {
		cachedHead.prev = lc.head
	}

	if lc.tail == nil {
		lc.tail = node
	}

	lc.head.next = cachedHead
	lc.head.prev = nil
	lc.index[key] = lc.head
}

func (lc *LruCache) Evict() {
	if lc.tail == nil {
		return // Nothing to evict
	}

	nodeToRemove := lc.tail
	delete(lc.index, nodeToRemove.key)

	if lc.tail.prev != nil {
		// There's more than one node...
		lc.tail = lc.tail.prev
		lc.tail.next = nil
	} else {
		// Only one node - list becomes empty
		lc.tail = nil
		lc.head = nil
	}

	lc.currentSize--
}

func (lc *LruCache) Put(key int, val string) {
	targetNode, exists := lc.index[key]

	if exists {
		targetNode.val = val
		// Only move to head if not already there
		if targetNode != lc.head {
			lc.NewHead(targetNode, key)
		}
	} else {

		if lc.currentSize == lc.capacity {
			lc.Evict()
		}

		newNode := &Node{val: val, key: key}
		lc.NewHead(newNode, newNode.key)
		lc.currentSize++
	}

}

func (lc *LruCache) Get(key int) (string, bool) {
	targetNode, exists := lc.index[key]
	if !exists {
		return "", false
	}
	lc.NewHead(targetNode, targetNode.key)

	return targetNode.val, true
}

func NewLruCache(capacity int) *LruCache {
	return &LruCache{
		capacity: capacity,
		index:    make(map[int]*Node, capacity),
	}

}
