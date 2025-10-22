//go:build lru_cache

package main

import "fmt"

// Some page in our buffer pool
type Node struct {
	next     *Node
	prev     *Node
	contents string
}

type LruCache struct {
	head     *Node
	tail     *Node
	capacity int
	index    map[int]*Node
}

func (lc *LruCache) put(key int, val string) {
	// If key exists the update it and move to MSU
	_, exists := lc.index[key]

	if exists {
		fmt.Print("key found: %d", key)
	}
	// If not insert as new head

}

func main() {

}
