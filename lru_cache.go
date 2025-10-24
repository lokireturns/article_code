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
	var penultimateNode *Node = lc.tail.prev
	lc.tail = penultimateNode
	lc.tail.next = nil
}

func (lc *LruCache) put(key int, val string) {
	// If key exists the update it and move to MSU
	targetNode, exists := lc.index[key]

	if exists {
		// Assuming its not the head or tail
		fmt.Print("key found: %d", key)
	} else {
		cachedNode := *lc.head
		newNode := Node{next: &cachedNode, prev: lc.tail, contents: val}
		cachedNode.prev = &newNode
	}

}

func main() {

}
