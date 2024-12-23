// Version 1.0.0
package randomness

import (
	"crypto/sha512"
	"encoding/binary"
	"math"
	"math/big"
	"strconv"
)

// Randomness provides numbers from a source of randomness.
type Randomness interface {
	// Probability returns a uniformly distributed random number in the range
	// (0.0, 1.0], based on the underlying uint64 value.
	Probability() float64

	// Bits returns a slice of boolean values representing the bits of the
	// underlying byte slice. A whole byte is consumed at a time even if n is
	// less than a multiple of 8.
	Bits(n int) BitArray

	// Bytes returns a slice of bytes from the underlying byte slice.
	Bytes(n int) []byte

	// Uint64 reads a uint64 value from the underlying byte slice.
	Uint64() uint64

	// Uint32 reads a uint32 value from the underlying byte slice.
	Uint32() uint32

	// Uint16 reads a uint16 value from the underlying byte slice.
	Uint16() uint16

	// Uint8 reads a uint8 value from the underlying byte slice.
	Uint8() uint8

	// Int64 reads an int64 value from the underlying byte slice.
	Int64() int64

	// Int32 reads an int32 value from the underlying byte slice.
	Int32() int32

	// Int16 reads an int16 value from the underlying byte slice.
	Int16() int16

	// Int8 reads an int8 value from the underlying byte slice.
	Int8() int8

	// Float64 reads a float64 value from the underlying byte slice.
	Float64() float64

	// Float32 reads a float32 value from the underlying byte slice.
	Float32() float32

	// Select returns a unique (no duplicates) set of ints in the
	// range [0, magnitude).
	Select(n int, magnitude int) []int

	// Numbers returns a readable list of up to count ints in the
	// range (0, magnitude].
	Numbers(count, magnitude int) Numbers
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
func NewRandomness(β string) Randomness {
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
func (b *randomness) get(n int) []byte {
	b.need(n)
	tmp := b.data[b.pos : b.pos+n]
	b.pos += n
	return tmp
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

func (b *randomness) Probability() float64 {
	p, _ := U64ToProbability(b.Uint64()).Float64()
	return p
}

func (b *randomness) Bits(n int) BitArray {
	if n < 0 {
		panic("randomness: bits out of range")
	}

	var bits []bool
	for {
		byte := b.get(1)[0]
		for i := 0; i < 8; i++ {
			bits = append(bits, (byte&(1<<uint(7-i))) != 0)
		}
		if len(bits) >= n {
			break
		}
	}

	return bits[:n]
}

func (b *randomness) Bytes(n int) []byte {
	return b.get(n)
}

func (b *randomness) Uint64() uint64 {
	return binary.BigEndian.Uint64(b.get(8))
}

func (b *randomness) Uint32() uint32 {
	return binary.BigEndian.Uint32(b.get(4))
}

func (b *randomness) Uint16() uint16 {
	return binary.BigEndian.Uint16(b.get(2))
}

func (b *randomness) Uint8() uint8 {
	return b.get(1)[0]
}

func (b *randomness) Int64() int64 {
	v := int64(b.Uint64())
	return v
}

func (b *randomness) Int32() int32 {
	v := int32(b.Uint32())
	return v
}

func (b *randomness) Int16() int16 {
	v := int16(b.Uint16())
	return v
}

func (b *randomness) Int8() int8 {
	v := int8(b.Uint8())
	return v
}

func (b *randomness) Float64() float64 {
	v := math.Float64frombits(b.Uint64())
	return v
}

func (b *randomness) Float32() float32 {
	v := math.Float32frombits(b.Uint32())
	return v
}

func (b *randomness) Select(n int, magnitude int) []int {
	if n < 0 || n > magnitude {
		panic("randomness: select out of range")
	}

	set := make([]int, magnitude)

	for i := range set {
		set[i] = int(i)
	}

	numbers := b.Numbers(n, magnitude)

	selected := make([]int, n)
	for i := 0; i < n; i++ {
		pos := numbers.Read(len(set))
		selected[i] = set[pos]
		set = append(set[:pos], set[pos+1:]...)
	}

	return selected
}

func (b *randomness) Numbers(count, magnitude int) Numbers {
	perNumber, bytesNeeded := numbersNeeds(count, magnitude)
	bits := b.Bits(bytesNeeded * 8)

	return newNumbers(bits, perNumber, count, magnitude)
}

// numbersNeeds returns the number of bits per number and maximum total bytes
// needed to load `count` numbers of `magnitude` size.
func numbersNeeds(count, magnitude int) (bitsPerNumber, bytesNeeded int) {
	bitsPerNumber = int(math.Ceil(math.Log2(float64(magnitude))))
	bitsNeeded := bitsPerNumber * (count + 1)
	bytesNeeded = int(math.Ceil(float64(bitsNeeded) / 8))
	return
}

// TestStringForValue returns a string representing the given values when used
// by NewRandomness and decoded. It's designed to be used for unit testing.
func TestStringForValue(vals ...any) string {
	var buf []byte

	var addVal func(val any)
	addVal = func(val any) {
		switch v := val.(type) {
		case uint64:
			tmp := make([]byte, 8)
			binary.BigEndian.PutUint64(tmp, v)
			buf = append(buf, tmp...)
		case uint32:
			tmp := make([]byte, 4)
			binary.BigEndian.PutUint32(tmp, v)
			buf = append(buf, tmp...)
		case uint16:
			tmp := make([]byte, 2)
			binary.BigEndian.PutUint16(tmp, v)
			buf = append(buf, tmp...)
		case uint8:
			buf = append(buf, v)
		case int64:
			addVal(uint64(v))
		case int32:
			addVal(uint32(v))
		case int:
			if strconv.IntSize == 64 {
				addVal(uint64(v))
			} else {
				addVal(uint32(v))
			}
		default:
			panic("randomness: invalid value type")
		}
	}

	for _, v := range vals {
		addVal(v)
	}

	return string(buf)
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
