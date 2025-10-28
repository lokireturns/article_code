// Package lruk implements an LRU-K cache replacement policy.
//
// This implementation uses a single data structure approach where each node
// maintains its own access history, trading O(n) eviction scans for
// implementation simplicity.
//
// The replacer manages eviction policy only - capacity enforcement is
// handled at the buffer pool level.
package lruk

import (
	"fmt"
	"time"
)

type Node struct {
	next      *Node
	prev      *Node
	key       int
	access    []int
	evictable bool
}

// TODO Should we refer to these as frames?
// Since the value member is where we store the pointer to the actual page?
func NewNode(key int, k int) *Node {
	return &Node{
		key:    key,
		access: make([]int, k),
	}
}

// LrukCache implements the LRU-K replacement algorithm.
//
// Design decisions:
// - Non-thread safe, safety managed by buffer pool
// - Single history list per node (vs dual-list approach)
// - O(n) eviction scan for simplicity
// - Capacity managed externally by buffer pool
type LrukCache struct {
	head        *Node
	tail        *Node
	index       map[int]*Node
	k           int
	currentSize int
}

func NewLrukCache(capacity int, k int) *LrukCache {
	return &LrukCache{
		k:     k,
		index: make(map[int]*Node),
	}
}

// Called when a page is pinned in the buffer pool
func (lc *LrukCache) RecordAccess(frame_id int, page *Node) {
	targetPage, exists := lc.index[frame_id]
	if exists {
		historySize := len(targetPage.access)
		newHistory := make([]int, historySize)
		for i := 0; i < historySize; i++ {
			newHistory[i] = targetPage.access[i+1]
		}
		targetPage.access[len(targetPage.access)-1] = int(time.Now().UnixMilli())
	} else {
		page.access[0] = int(time.Now().UnixMilli())
		lc.tail.next = page
		lc.tail = page
		lc.index[frame_id] = page
	}
}

// Toggle page as evictable or not
// When page pin count hits 0 page is evictable
// - Controls the cache size value
// - Assumes pin count is known when called from buffer pool
func (lc *LrukCache) MarkEvictable(frame_id int, set_evictable bool) (int, error) {
	targetPage, exists := lc.index[frame_id]
	if !exists {
		return -1, fmt.Errorf("frame ID: %d not found", frame_id)
	}
	if set_evictable {
		targetPage.evictable = true
		lc.currentSize++
		return 0, nil
	}

	targetPage.evictable = false
	lc.currentSize--
	return 0, nil
}

// Evict frame with largest backward k-distance
// among evictable frames.
func (lc *LrukCache) Evict() int {
	if lc.currentSize < 1 {
		return -1
	}
	now := time.Now().UnixMilli()
	var maxkBackwardsDistance int64
	var evictionFrameId int

	for i, p := range lc.index {
		if p.evictable {
			kBackwardsDistance := now - int64(p.access[(lc.k)-1])
			if kBackwardsDistance > maxkBackwardsDistance {
				maxkBackwardsDistance = kBackwardsDistance
				evictionFrameId = i
			}
		}
	}
	targetNode := lc.index[evictionFrameId]

	// TODO do we need to nil the next and prev pointers of our
	// target node to ensure its de allocated by garbage collector?
	if lc.head == targetNode {
		lc.head = targetNode.next
	}

	if lc.tail == targetNode {
		lc.tail = targetNode.prev
	}

	targetNode.prev.next = targetNode.next
	targetNode.next.prev = targetNode.prev
	lc.currentSize--
	return evictionFrameId
}

// Return number of evictable frames currently in the cache
func (lc *LrukCache) Size() int {
	return lc.currentSize
}
