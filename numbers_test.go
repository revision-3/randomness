package randomness

import (
	"testing"
)

func TestNumbersUniqueness(t *testing.T) {
	tests := []struct {
		name      string
		count     int
		magnitude int
	}{
		{"small range", 5, 10},
		{"medium range", 10, 20},
		{"large range", 100, 200},
		{"full range", 10, 10}, // Testing when count equals magnitude
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test randomness instance with known values
			r := NewRandomness(TestStringForValue(GenerateTestRandomValue()))
			generated, err := r.PickDistinct(tt.count, tt.magnitude)
			if err != nil {
				t.Errorf("PickDistinct(%d, %d) returned error: %v", tt.count, tt.magnitude, err)
			}

			// Check that all numbers are within range
			for i, num := range generated {
				if num < 0 || num >= tt.magnitude {
					t.Errorf("number %d at index %d is out of range [0, %d)", num, i, tt.magnitude)
				}
			}

			// Check for uniqueness using a map
			seen := make(map[int]bool)
			for _, num := range generated {
				if seen[num] {
					t.Errorf("duplicate number %d found", num)
				}
				seen[num] = true
			}

			// Verify we got exactly the expected count of unique numbers
			if len(seen) != tt.count {
				t.Errorf("got %d unique numbers, want %d", len(seen), tt.count)
			}
		})
	}
}

func TestNumbersDistribution(t *testing.T) {
	tests := []struct {
		name      string
		count     int
		magnitude int
	}{
		{"small range", 5, 10},
		{"medium range", 10, 20},
		{"large range", 100, 200},
		{"full range", 10, 10}, // Testing when count equals magnitude
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test randomness instance with known values
			r := NewRandomness(TestStringForValue(GenerateTestRandomValue()))
			numbers, err := r.Numbers(tt.count, tt.magnitude)
			if err != nil {
				t.Errorf("Numbers(%d, %d) returned error: %v", tt.count, tt.magnitude, err)
			}

			// Generate all numbers
			generated := make([]int, tt.count)
			for i := 0; i < tt.count; i++ {
				generated[i], err = numbers.Read(tt.magnitude)
				if err != nil {
					t.Errorf("Numbers(%d, %d) returned error: %v", tt.count, tt.magnitude, err)
				}
			}

			// Check that all numbers are within range
			for i, num := range generated {
				if num < 0 || num >= tt.magnitude {
					t.Errorf("number %d at index %d is out of range [0, %d)", num, i, tt.magnitude)
				}
			}

			// Check distribution
			distribution := make(map[int]int)
			for _, num := range generated {
				distribution[num]++
			}

			// Verify we got the expected count of numbers
			if len(generated) != tt.count {
				t.Errorf("got %d numbers, want %d", len(generated), tt.count)
			}
		})
	}
}
