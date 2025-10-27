//go:build lru_cache

package main

import "sync"

// Some page in our buffer pool
type Node struct {
	next *Node
	prev *Node
	val  string
	key  int
}

type LruCache struct {
	mu          sync.Mutex
	head        *Node
	tail        *Node
	capacity    int
	index       map[int]*Node
	currentSize int
}

func (lc *LruCache) newHead(node *Node, key int) {
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

func (lc *LruCache) evict() {
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

func (lc *LruCache) put(key int, val string) {
	targetNode, exists := lc.index[key]

	if exists {
		targetNode.val = val
		// Only move to head if not already there
		if targetNode != lc.head {
			lc.newHead(targetNode, key)
		}
	} else {

		if lc.currentSize == lc.capacity {
			lc.evict()
		}

		newNode := &Node{val: val, key: key}
		lc.newHead(newNode, newNode.key)
		lc.currentSize++
	}

}

func (lc *LruCache) get(key int) (string, bool) {
	targetNode, exists := lc.index[key]
	if !exists {
		return "", false
	}
	lc.newHead(targetNode, targetNode.key)

	return targetNode.val, true
}

func newCache(capacity int) *LruCache {
	return &LruCache{
		capacity: capacity,
		index:    make(map[int]*Node, capacity),
	}

}

func main() {
	cache := newCache(2)
	cache.put(35, "foo")
	cache.put(22, "bar")
	cache.put(22, "Cheese")
}
