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
	key       int
	val       any
	access    []int
	evictable bool
}

// TODO Should we refer to these as frames?
// Since the value member is where we store the pointer to the actual page?
func NewNode(key int, k int, val any) *Node {
	return &Node{
		key:    key,
		access: make([]int, k),
		val:    val,
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
	index    map[int]*Node
	history  []*Node
	capacity int
	k        int
}

func NewLrukCache(capacity int, k int) *LrukCache {
	return &LrukCache{
		k:        k,
		capacity: capacity,
		index:    make(map[int]*Node, capacity),
		history:  make([]*Node, capacity),
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
		lc.history = append(lc.history, page)
		lc.index[frame_id] = lc.history[0]
	}
}

// Toggle page as evictable or not
// When page pin count hits 0 page is evictable
// - Controls the cache capacity value
// - Assumes pin count is known when called from buffer pool
func (lc *LrukCache) MarkEvictable(frame_id int, set_evictable bool) (int, error) {
	targetPage, exists := lc.index[frame_id]
	if !exists {
		return -1, fmt.Errorf("frame ID: %d not found", frame_id)
	}
	if set_evictable {
		targetPage.evictable = true
		lc.capacity++
		return 0, nil
	}

	targetPage.evictable = false
	lc.capacity--
	return 0, nil
}

// Return number of evictable frames currently in the cache
// func (lc *LrukCache) Size() int {

// }
