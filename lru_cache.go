//go:build lru_cache

package main

import "fmt"

// Some page in our buffer pool
type Node struct {
	next *Node
	prev *Node
	val  string
}

type LruCache struct {
	head        *Node
	tail        *Node
	capacity    int
	index       map[int]*Node
	currentSize int
}

func (lc *LruCache) newHead(node *Node, key int) {
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
	if 
}

func (lc *LruCache) put(key int, val string) {
	// If key exists the update it and move to MSU
	targetNode, exists := lc.index[key]

	if exists {
		// Assuming its not the head or tail
		fmt.Print("key found: %d", key)

		// Deal with targets next node
		var targetNodeNext *Node = targetNode.next
		targetNodeNext.prev = targetNode.prev

		var targetNodePrev *Node = targetNode.prev
		targetNodePrev.next = targetNode.next

		lc.newHead(targetNode, key)
		lc.head.val = val

	} else {
		// If not check if full & evict if needed
		if lc.currentSize == lc.capacity {
			// Evict
		} else {
			// insert as new head
			newNode := Node{}
			lc.newHead(&newNode, key)
		}
	}
}

func main() {

}
