//go:build lruk_cache

package articlecontent

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

type LrukCache struct {
	head        *Node
	tail        *Node
	index       map[int]*Node
	currentSize int
	capacity    int
	k           int
}
