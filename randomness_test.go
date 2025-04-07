package randomness

import (
	"bytes"
	"crypto/sha256"
	"math"
	"testing"
)

func TestBetaString(t *testing.T) {
	tests := []struct {
		name     string
		i        uint64
		n        uint64
		expected string
	}{
		{
			name:     "zero",
			i:        0,
			n:        1,
			expected: string(make([]byte, 32)), // 32 zero bytes
		},
		{
			name:     "one half",
			i:        1,
			n:        2,
			expected: string(append([]byte{0x80}, make([]byte, 31)...)), // 0x80 followed by 31 zeros
		},
		{
			name:     "one",
			i:        1,
			n:        1,
			expected: string(make([]byte, 32)), // 32 0xFF bytes
		},
		{
			name:     "one quarter",
			i:        1,
			n:        4,
			expected: string(append([]byte{0x40}, make([]byte, 31)...)), // 0x40 followed by 31 zeros
		},
		{
			name:     "three quarters",
			i:        3,
			n:        4,
			expected: string(append([]byte{0xC0}, make([]byte, 31)...)), // 0xC0 followed by 31 zeros
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BetaString(tt.i, tt.n)
			if got != tt.expected {
				t.Errorf("BetaString(%d, %d) = %x, want %x", tt.i, tt.n, []byte(got), []byte(tt.expected))
			}
		})
	}
}

func TestBetaStringPanics(t *testing.T) {
	tests := []struct {
		name string
		i    uint64
		n    uint64
	}{
		{
			name: "zero denominator",
			i:    1,
			n:    0,
		},
		{
			name: "numerator greater than denominator",
			i:    2,
			n:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("BetaString(%d, %d) did not panic", tt.i, tt.n)
				}
			}()
			BetaString(tt.i, tt.n)
		})
	}
}

func TestHashValues(t *testing.T) {
	tests := []struct {
		name     string
		i        uint64
		expected BetaBytes
	}{
		{
			name:     "zero",
			i:        0,
			expected: func() BetaBytes { h := sha256.Sum256([]byte{0, 0, 0, 0, 0, 0, 0, 0}); return BetaBytes(h[:]) }(),
		},
		{
			name:     "one",
			i:        1,
			expected: func() BetaBytes { h := sha256.Sum256([]byte{0, 0, 0, 0, 0, 0, 0, 1}); return BetaBytes(h[:]) }(),
		},
		{
			name: "max uint64",
			i:    math.MaxUint64,
			expected: func() BetaBytes {
				h := sha256.Sum256([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
				return BetaBytes(h[:])
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashValues(tt.i)
			// Compare the bytes
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("HashValues(%d) = %x, want %x", tt.i, got, tt.expected)
			}
		})
	}
}
