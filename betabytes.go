package randomness

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
)

type BetaBytes []byte

func (b BetaBytes) String() string {
	// Hex encode the bytes
	return hex.EncodeToString(b)
}

func (b BetaBytes) Bytes() []byte {
	return b
}

func BetaBytesFromHex(β string) (BetaBytes, error) {
	bytes, err := hex.DecodeString(β)
	if err != nil {
		return nil, err
	}
	return BetaBytes(bytes), nil
}

func MustBetaBytesFromHex(β string) BetaBytes {
	bytes, err := BetaBytesFromHex(β)
	if err != nil {
		panic(err)
	}
	return bytes
}

// BetaValues returns a BetaBytes representing the given values when
// used by NewRandomness and decoded.
// Designed for unit testing.
func BetaValues(vals ...any) BetaBytes {
	if len(vals) == 0 {
		panic("no values provided to BetaValues")
	}

	var buf []byte

	var addVal func(val any)
	addVal = func(val any) {
		if val == nil {
			buf = append(buf, 0)
			return
		}
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
		case int16:
			addVal(uint16(v))
		case int8:
			addVal(uint8(v))
		case int:
			if strconv.IntSize == 64 {
				addVal(uint64(v))
			} else {
				addVal(uint32(v))
			}
		case float64:
			tmp := make([]byte, 8)
			binary.BigEndian.PutUint64(tmp, math.Float64bits(v))
			buf = append(buf, tmp...)
		case float32:
			tmp := make([]byte, 4)
			binary.BigEndian.PutUint32(tmp, math.Float32bits(v))
			buf = append(buf, tmp...)
		case string:
			buf = append(buf, []byte(v)...)
		case []byte:
			buf = append(buf, v...)
		case []uint16, []uint32, []uint64, []int16, []int32, []int64, []float32, []float64, []string:
			// Handle all array types by recursively calling addVal for each element
			rv := reflect.ValueOf(v)
			for i := 0; i < rv.Len(); i++ {
				addVal(rv.Index(i).Interface())
			}
		default:
			panic(fmt.Sprintf("unsupported value type: %T", val))
		}
	}

	for _, v := range vals {
		addVal(v)
	}

	return BetaBytes(buf)
}

// HashValues generates a 32-byte string by hashing the input values using
// SHA-256. The input values are converted to bytes using BetaValues before hashing.
// Designed for unit testing.
func HashValues(values ...any) BetaBytes {
	// Convert values to bytes using BetaValues
	bytes := BetaValues(values...).Bytes()

	// Hash the bytes using SHA-256
	hash := sha256.Sum256(bytes)

	// Convert the hash to a string
	return BetaBytes(hash[:])
}

// BetaString generates a 32-byte string that represents the precise fraction i/n
// between 0x0000000000000000000000000000000000000000000000000000000000000000
// and 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF.
// The string is formatted as the raw bytes.
// Designed for unit testing.
func BetaString(i, n uint64) string {
	if n == 0 {
		panic("denominator cannot be zero")
	}
	if i > n {
		panic("numerator cannot be greater than denominator")
	}

	// Convert to big.Float for precise arithmetic
	bigI := new(big.Float).SetUint64(i)
	bigN := new(big.Float).SetUint64(n)

	// Calculate the fraction
	fraction := new(big.Float).Quo(bigI, bigN)

	// Create a 256-bit number (32 bytes)
	maxValue := new(big.Float).SetUint64(1)
	for i := 0; i < 32; i++ {
		maxValue = new(big.Float).Mul(maxValue, new(big.Float).SetUint64(256))
	}

	// Scale the fraction to the full range
	scaled := new(big.Float).Mul(fraction, maxValue)

	// Convert to big.Int
	result := new(big.Int)
	scaled.Int(result)

	// Convert to bytes
	bytes := result.Bytes()
	output := make([]byte, 32)

	// Pad with leading zeros if necessary
	if len(bytes) < 32 {
		copy(output[32-len(bytes):], bytes)
	} else {
		copy(output, bytes[len(bytes)-32:])
	}

	return string(output)
}
