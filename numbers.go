// Version 1.0.0
package randomness

import "fmt"

// Numbers allows reading a fixed number of unique random numbers in a
// range [0, magnitude) from a source of random data.
// It does so efficiently by reading a fixed number of bits at a time
// from a BitArray and adding them to an accumulator which is used to
// derive the next number.
type Numbers interface {
	Read(size int) (int, error)
	Readn(count, size int) ([]int, error)
}

type numbers struct {
	bits        BitArray
	pos         int
	accumulator uint64
	perNumber   int
	count       int
	magnitude   int
}

func NewNumbers(bits BitArray, perNumber, count, magnitude int) Numbers {
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
func (n *numbers) Read(magnitude int) (int, error) {
	if n.pos > n.count+1 {
		return 0, fmt.Errorf("out of numbers")
	}
	if magnitude > n.magnitude {
		return 0, fmt.Errorf("magnitude too large")
	}
	if magnitude <= 0 {
		return 0, fmt.Errorf("magnitude must be positive")
	}
	input := n.bits.NumberAt((n.pos)*n.perNumber, n.perNumber)
	n.pos += 1
	n.accumulator += input
	return int(n.accumulator % uint64(magnitude)), nil
}

// Readn returns a slice of count random numbers in the range [0, magnitude).
func (n *numbers) Readn(count, magnitude int) (nums []int, err error) {
	for range count {
		num, err := n.Read(magnitude)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}
	return nums, nil
}
