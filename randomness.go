// Version 1.0.0
package randomness

import (
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"slices"
)

// Randomness provides numbers from a source of randomness.
type Randomness interface {
	// Probability returns a uniformly distributed random number in the range
	// (0.0, 1.0], based on the underlying uint64 value.
	Probability() (float64, error)

	// Bits returns a slice of boolean values representing the bits of the
	// underlying byte slice. A whole byte is consumed at a time even if n is
	// less than a multiple of 8.
	Bits(n int) (BitArray, error)

	// Bytes returns a slice of bytes from the underlying byte slice.
	Bytes(n int) ([]byte, error)

	// Uint64 reads a uint64 value from the underlying byte slice.
	Uint64() (uint64, error)

	// Uint32 reads a uint32 value from the underlying byte slice.
	Uint32() (uint32, error)

	// Uint16 reads a uint16 value from the underlying byte slice.
	Uint16() (uint16, error)

	// Uint8 reads a uint8 value from the underlying byte slice.
	Uint8() (uint8, error)

	// Int64 reads an int64 value from the underlying byte slice.
	Int64() (int64, error)

	// Int32 reads an int32 value from the underlying byte slice.
	Int32() (int32, error)

	// Int16 reads an int16 value from the underlying byte slice.
	Int16() (int16, error)

	// Int8 reads an int8 value from the underlying byte slice.
	Int8() (int8, error)

	// Float64 reads a float64 value from the underlying byte slice.
	Float64() (float64, error)

	// Float32 reads a float32 value from the underlying byte slice.
	Float32() (float32, error)

	// Numbers returns a readable list of random numbers in the
	// range [0, magnitude). Note that this method does not guarantee
	// uniqueness of the returned values. For unique numbers, use
	// the PickDistinct method instead.
	Numbers(count, magnitude int) (Numbers, error)

	// Selection returns a selection of items from the list of items, and can handle complex item configurations, such as weights and limited supplies.
	Selection(cfg SelectionConfig) ([]SelectionResult, error)

	// PickDistinct returns n unique random integers in [0, magnitude)
	PickDistinct(n int, magnitude int) ([]int, error)

	// Pick returns n random integers in [0, magnitude) (may include duplicates)
	Pick(n int, magnitude int) ([]int, error)
}

// randomness implements the Randomness interface.
type randomness struct {
	data          []byte
	pos           int
	originalLen   int
	amplification int
}

// Randomness is implemented by *randomness.
var _ Randomness = (*randomness)(nil)

// NewRandomness creates a new Randomness instance from a random string.
func NewRandomness(β BetaBytes) Randomness {
	return &randomness{
		data:        []byte(β),
		pos:         0,
		originalLen: len(β),
	}
}

// have tells you how much entropy you have remaining.
func (b *randomness) have() int {
	return len(b.data) - b.pos
}

// need ensures you have enough entropy to perform your current operation.
func (b *randomness) need(n int) {
	for b.have() < n {
		b.amplify()
	}
}

// get returns a slice of bytes from the underlying byte slice.
func (b *randomness) get(n int) ([]byte, error) {
	b.need(n)
	tmp := b.data[b.pos : b.pos+n]
	b.pos += n
	return tmp, nil
}

// amplify extends the underlying byte slice using SHA-512.
func (b *randomness) amplify() {
	b.amplification++
	sha := sha512.New()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(b.amplification))
	sha.Write(b.data[:b.originalLen])
	sha.Write(buf)
	tmp := sha.Sum(nil)
	b.data = append(b.data, tmp...)
}

func (b *randomness) Probability() (float64, error) {
	u64, err := b.Uint64()
	if err != nil {
		return 0, err
	}
	p, _ := U64ToProbability(u64).Float64()
	return p, nil
}

func (b *randomness) Bits(n int) (BitArray, error) {
	if n < 0 {
		return nil, fmt.Errorf("cannot generate %d bits: count must be non-negative", n)
	}

	var bits []bool
	for {
		bytes, err := b.get(1)
		if err != nil {
			return nil, err
		}
		byte := bytes[0]
		for i := range 8 {
			bits = append(bits, (byte&(1<<uint(7-i))) != 0)
		}
		if len(bits) >= n {
			break
		}
	}

	return bits[:n], nil
}

func (b *randomness) Bytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, fmt.Errorf("cannot generate %d bytes: count must be non-negative", n)
	}
	return b.get(n)
}

func (b *randomness) Uint64() (uint64, error) {
	bytes, err := b.get(8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(bytes), nil
}

func (b *randomness) Uint32() (uint32, error) {
	bytes, err := b.get(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(bytes), nil
}

func (b *randomness) Uint16() (uint16, error) {
	bytes, err := b.get(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(bytes), nil
}

func (b *randomness) Uint8() (uint8, error) {
	bytes, err := b.get(1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}

func (b *randomness) Int64() (int64, error) {
	u64, err := b.Uint64()
	if err != nil {
		return 0, err
	}
	return int64(u64), nil
}

func (b *randomness) Int32() (int32, error) {
	u32, err := b.Uint32()
	if err != nil {
		return 0, err
	}
	return int32(u32), nil
}

func (b *randomness) Int16() (int16, error) {
	u16, err := b.Uint16()
	if err != nil {
		return 0, err
	}
	return int16(u16), nil
}

func (b *randomness) Int8() (int8, error) {
	u8, err := b.Uint8()
	if err != nil {
		return 0, err
	}
	return int8(u8), nil
}

func (b *randomness) Float64() (float64, error) {
	u64, err := b.Uint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(u64), nil
}

func (b *randomness) Float32() (float32, error) {
	u32, err := b.Uint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(u32), nil
}

func (b *randomness) Numbers(count, magnitude int) (Numbers, error) {
	if count < 0 {
		return nil, fmt.Errorf("cannot generate %d numbers: count must be non-negative", count)
	}
	if magnitude <= 0 {
		return nil, fmt.Errorf("cannot generate numbers in range [0, %d): magnitude must be positive", magnitude)
	}

	perNumber, bytesNeeded, err := numbersNeeds(count, magnitude)
	if err != nil {
		return nil, err
	}
	bits, err := b.Bits(bytesNeeded * 8)
	if err != nil {
		return nil, err
	}

	return NewNumbers(bits, perNumber, count, magnitude), nil
}

// numbersNeeds returns the number of bits per number and maximum total bytes
// needed to load `count` numbers of `magnitude` size.
func numbersNeeds(count, magnitude int) (bitsPerNumber, bytesNeeded int, err error) {
	if count <= 0 {
		err = fmt.Errorf("invalid count %d: must be positive", count)
		return
	}
	if magnitude <= 0 {
		err = fmt.Errorf("invalid magnitude %d: must be positive", magnitude)
		return
	}

	bitsPerNumber = int(math.Ceil(math.Log2(float64(magnitude))))
	bitsNeeded := bitsPerNumber * (count + 1)
	bytesNeeded = int(math.Ceil(float64(bitsNeeded) / 8))
	return bitsPerNumber, bytesNeeded, nil
}

var bigMaxUint64 = big.NewFloat(float64(math.MaxUint64))
var bigSmallestNonzeroFloat64 = big.NewFloat(math.SmallestNonzeroFloat64)

func U64ToProbability(u uint64) *big.Float {
	// Convert the uint64 to a float in the range [0.0, 1.0).
	p := new(big.Float)
	p.SetUint64(u)
	p.Quo(p, bigMaxUint64)
	p.Add(p, bigSmallestNonzeroFloat64)
	// Add the smallest nonzero float64 to the result to avoid returning 0.0:
	// (0.0, 1.0) -> [0.0, 1.0).
	return p
}

// PickDistinct returns n unique random integers in [0, magnitude)
func (b *randomness) PickDistinct(n int, magnitude int) ([]int, error) {
	if n < 0 {
		return nil, fmt.Errorf("cannot generate %d numbers: count must be non-negative", n)
	}
	if magnitude <= 0 {
		return nil, fmt.Errorf("cannot generate numbers in range [0, %d): magnitude must be positive", magnitude)
	}
	if n > magnitude {
		return nil, fmt.Errorf("cannot pick %d distinct numbers from a range of only %d numbers", n, magnitude)
	}

	set := make([]int, magnitude)
	for i := range set {
		set[i] = int(i)
	}

	numbers, err := b.Numbers(n, magnitude)
	if err != nil {
		return nil, err
	}

	selected := make([]int, n)
	for i := range n {
		pos, err := numbers.Read(len(set))
		if err != nil {
			return nil, err
		}
		selected[i] = set[pos]
		set = slices.Delete(set, pos, pos+1)
	}

	return selected, nil
}

// Pick returns n random integers in [0, magnitude) (may include duplicates)
func (b *randomness) Pick(n int, magnitude int) ([]int, error) {
	if n < 0 {
		return nil, fmt.Errorf("cannot generate %d numbers: count must be non-negative", n)
	}
	if magnitude <= 0 {
		return nil, fmt.Errorf("cannot generate numbers in range [0, %d): magnitude must be positive", magnitude)
	}

	numbers, err := b.Numbers(n, magnitude)
	if err != nil {
		return nil, err
	}

	result := make([]int, n)
	for i := range n {
		result[i], err = numbers.Read(magnitude)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
