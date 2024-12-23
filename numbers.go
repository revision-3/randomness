// Version 1.0.0
package randomness

// Numbers allows reading a fixed number of unique random numbers in a
// range [0, magnitude) from a source of random data.
// It does so efficiently by reading a fixed number of bits at a time
// from a BitArray and adding them to an accumulator which is used to
// derive the next number.
type Numbers interface {
	Read(size int) int
	Readn(count, size int) []int
}

type numbers struct {
	bits        BitArray
	pos         int
	accumulator uint64
	perNumber   int
	count       int
	magnitude   int
}

func newNumbers(bits BitArray, perNumber, count, magnitude int) Numbers {
	nums := &numbers{
		bits:      bits,
		perNumber: perNumber,
		count:     count,
		magnitude: magnitude,
	}
	// Siphon off the initial number to prime the overflow and make
	// statistical prediction of the first number difficult.
	nums.Read(magnitude)
	return nums
}

// Read returns a random number in the range [0, magnitude).
func (n *numbers) Read(magnitude int) int {
	if n.pos > n.count+1 {
		panic("out of numbers")
	}
	if magnitude > n.magnitude {
		panic("magnitude too large")
	}
	if magnitude <= 0 {
		panic("magnitude must be positive")
	}
	input := n.bits.NumberAt((n.pos)*n.perNumber, n.perNumber)
	n.pos += 1
	n.accumulator += input
	return int(n.accumulator % uint64(magnitude))
}

// Readn returns a slice of count random numbers in the range [0, magnitude).
func (n *numbers) Readn(count, magnitude int) (nums []int) {
	for i := 0; i < count; i++ {
		nums = append(nums, n.Read(magnitude))
	}
	return
}
