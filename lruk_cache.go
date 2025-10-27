//go:build lruk_cache

package main

type Node struct {
	next              *Node
	prev              *Node
	key               int
	val               string
	access_timestamps []int
}

func newNode(k int) *Node {
	return &Node{
		access_timestamps: make([]int, k),
	}
}

type CachingQueue struct {
	head     *Node
	tail     *Node
	capacity *Node
	index    map[int]*Node
}

type LrukCache struct {
	head         *Node
	tail         *Node
	index        map[int]*Node
	historyQueue *CachingQueue
	currentSize  int
	capacity     int
	k            int
}

func newCache(capacity int) *LrukCache {
	return &LrukCache{
		capacity: capacity,
		index:    make(map[int]*Node, capacity),
	}
}

func (*LrukCache) put(key int, val string) {

}

func main() {

}
