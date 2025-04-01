// Version 1.0.0
package randomness

type BitArray []bool

func (b BitArray) Get(index int) bool {
	return b[index]
}

func (b BitArray) Length() int {
	return len(b)
}

func (b BitArray) NumberAt(index int, size int) uint64 {
	n := uint64(0)
	for i := range size {
		if b.Get(index + i) {
			n |= 1 << i
		}
	}
	return n
}
