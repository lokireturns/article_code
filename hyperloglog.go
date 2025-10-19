package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/bits"
)

/*
HyperLogLog Struct:

- Initialised with m buckets (array if ints?) which if of size 2^b
- Init with `b` number of leading bits we want to use (how do we decide that?)

1. Hash the input value
2. Convert it to binary representation
3. Find the `b` bits from MSB to find bucket number
   b - Number of initial bits in a binary representation of a hash value
4. Find `p` from remaining bits
5. Store in bucker number the max(m[b],p)
*/

type HyperLogLog struct {
	constant float64
	b        int
	buckets  []int
}

func makeHyperLogLog(b int) *HyperLogLog {
	return &HyperLogLog{
		constant: 0.79402,
		b:        b,
		buckets:  make([]int, 1<<b),
	}
}

func (hll *HyperLogLog) addItem(value string) {
	hasher := fnv.New64a() // Maybe this can be passed in as a pointer?
	hasher.Write([]byte(value))
	hashedValue := hasher.Sum64()

	// Extract b bits from MSB
	var bucketNumber uint64 = hashedValue >> (64 - hll.b) //(j)

	// Find p
	var remainingBits uint64 = hashedValue << hll.b
	p := bits.LeadingZeros64(remainingBits) + 1

	hll.buckets[bucketNumber] = max(hll.buckets[bucketNumber], p)
}

func (hll *HyperLogLog) estimateCardinality() float64 {
	var harmonicSum float64
	for i := 0; i < len(hll.buckets); i++ {
		if hll.buckets[i] > 0 {
			rj := float64(hll.buckets[i])
			harmonicSum += math.Pow(2, -rj)
		}
	}

	m := float64(len(hll.buckets))
	cardinality := hll.constant * m * (m / harmonicSum)
	return cardinality

}

func main() {

	data := []string{"a", "b", "c", "d"}
	estimator := makeHyperLogLog(14)
	for i := range data {
		estimator.addItem(data[i])
	}
	fmt.Printf("%1.f", estimator.estimateCardinality())

}
