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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run multiple iterations to check distribution
			iterations := 10000
			counts := make(map[int]int)

			for i := 0; i < iterations; i++ {
				r := NewRandomness(TestStringForValue(GenerateTestRandomValue()))
				numbers, err := r.Numbers(tt.count, tt.magnitude)
				if err != nil {
					t.Errorf("Numbers(%d, %d) returned error: %v", tt.count, tt.magnitude, err)
				}

				// Generate numbers and count their occurrences
				for j := 0; j < tt.count; j++ {
					num, err := numbers.Read(tt.magnitude)
					if err != nil {
						t.Errorf("Numbers(%d, %d) returned error: %v", tt.count, tt.magnitude, err)
					}
					counts[num]++
				}
			}

			// Calculate expected count for each number
			expectedCount := float64(iterations*tt.count) / float64(tt.magnitude)
			tolerance := expectedCount * 0.2 // 20% tolerance for randomness

			// Check distribution
			for num := 0; num < tt.magnitude; num++ {
				count := counts[num]
				deviation := float64(count) - expectedCount
				if deviation < 0 {
					deviation = -deviation
				}
				if deviation > tolerance {
					t.Errorf("number %d appears %d times, expected approximately %.0f (deviation: %.2f%%)",
						num, count, expectedCount, (deviation/expectedCount)*100)
				}
			}
		})
	}
}
