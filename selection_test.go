package randomness

import (
	"math"
	"testing"
)

type testItem struct {
	value  int
	weight float64
	supply int
}

func (i *testItem) Weight() float64 {
	return i.weight
}

func (i *testItem) Supply() int {
	return i.supply
}

func (i *testItem) Use() {
	if i.supply == -1 {
		return
	}
	i.supply--
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestSelectWeightedWithSupply(t *testing.T) {
	iterations := 100000
	counts := make(map[int]int)

	for i := 0; i < iterations; i++ {
		// Create items with different supplies and equal weights
		items := []Item{
			&testItem{value: 1, weight: 1.0, supply: 3}, // Should have 3/4 chance
			&testItem{value: 2, weight: 1.0, supply: 1}, // Should have 1/4 chance
		}

		r := NewRandomness(BetaValues(uint64(math.MaxUint64 / uint64(iterations) * uint64(i))))
		results, err := r.Selection(SelectionConfig{Items: items, Count: 1})
		if err != nil {
			t.Fatalf("Selection() error = %v", err)
		}
		counts[results[0].Get().(*testItem).value]++
	}

	// Check distribution matches expected probabilities
	expected := map[int]float64{
		1: float64(iterations) * 3.0 / 4.0, // 75% chance
		2: float64(iterations) * 1.0 / 4.0, // 25% chance
	}

	tolerance := 0.05 // 5% tolerance
	for value, count := range counts {
		expectedCount := expected[value]
		deviation := float64(abs(count-int(expectedCount))) / expectedCount
		if deviation > tolerance {
			t.Errorf("Value %d: count = %d, expected ≈ %.0f (deviation: %.2f%%)",
				value, count, expectedCount, deviation*100)
		}
	}
}

func TestSelectInstances(t *testing.T) {
	t.Run("Equal weights and supplies", func(t *testing.T) {
		iterations := 100000
		valueCounts := make(map[int]int)
		instanceCounts := make(map[struct{ value, instance int }]int)

		// Create items once for the entire test case
		items := []Item{
			&testItem{value: 1, weight: 1.0, supply: 2},
			&testItem{value: 2, weight: 1.0, supply: 2},
		}

		// Create a single randomness instance for the test
		r := NewRandomness(BetaValues(GenerateTestRandomValue()))

		// Create a single SelectionConfig for the test
		config := SelectionConfig{Items: items, Count: 1}

		for i := 0; i < iterations; i++ {
			results, err := r.Selection(config)
			if err != nil {
				t.Fatalf("Selection() error = %v", err)
			}

			result := results[0]
			item := result.Get().(*testItem)
			valueCounts[item.value]++
			instanceCounts[struct{ value, instance int }{item.value, result.Instance()}]++
		}

		// Verify distribution of values
		// Each value has 2 instances with weight 1.0, so should get half the iterations
		expectedValueCount := iterations / 2
		tolerance := float64(expectedValueCount) * 0.05 // 5% tolerance
		for value, count := range valueCounts {
			if math.Abs(float64(count)-float64(expectedValueCount)) > tolerance {
				t.Errorf("Value %d: count = %d, expected ≈ %d (deviation: %.2f%%)",
					value, count, expectedValueCount, math.Abs(float64(count)-float64(expectedValueCount))/float64(expectedValueCount)*100)
			}
		}

		// Verify equal distribution of instances within each value
		// Each instance should have 1/4 chance of being selected
		expectedInstanceCount := iterations / 4
		for key, count := range instanceCounts {
			if math.Abs(float64(count)-float64(expectedInstanceCount)) > tolerance {
				t.Errorf("Value %d Instance %d: count = %d, expected ≈ %d (deviation: %.2f%%)",
					key.value, key.instance, count, expectedInstanceCount, math.Abs(float64(count)-float64(expectedInstanceCount))/float64(expectedInstanceCount)*100)
			}
		}
	})

	t.Run("Different weights", func(t *testing.T) {
		iterations := 100000
		valueCounts := make(map[int]int)
		instanceCounts := make(map[struct{ value, instance int }]int)

		// Create items once for the entire test case
		items := []Item{
			&testItem{value: 1, weight: 1.0, supply: 2}, // 2 instances, weight 1.0 each
			&testItem{value: 2, weight: 1.5, supply: 2}, // 2 instances, weight 1.5 each
		}

		// Create a single randomness instance for the test
		r := NewRandomness(BetaValues(GenerateTestRandomValue()))

		// Create a single SelectionConfig for the test
		config := SelectionConfig{Items: items, Count: 1}

		for i := 0; i < iterations; i++ {
			results, err := r.Selection(config)
			if err != nil {
				t.Fatalf("Selection() error = %v", err)
			}

			result := results[0]
			item := result.Get().(*testItem)
			valueCounts[item.value]++
			instanceCounts[struct{ value, instance int }{item.value, result.Instance()}]++
		}

		// Verify distribution matches weights
		// Total weight = 2 * 1.0 + 2 * 1.5 = 5.0
		// Expected counts for value 1: iterations * (2 * 1.0 / 5.0) = iterations * 0.4
		// Expected counts for value 2: iterations * (2 * 1.5 / 5.0) = iterations * 0.6
		expectedValue1Count := float64(iterations) * 0.4
		expectedValue2Count := float64(iterations) * 0.6
		tolerance := float64(iterations) * 0.05 // 5% tolerance

		if math.Abs(float64(valueCounts[1])-expectedValue1Count) > tolerance {
			t.Errorf("Value 1: count = %d, expected ≈ %.0f (deviation: %.2f%%)",
				valueCounts[1], expectedValue1Count, math.Abs(float64(valueCounts[1])-expectedValue1Count)/expectedValue1Count*100)
		}
		if math.Abs(float64(valueCounts[2])-expectedValue2Count) > tolerance {
			t.Errorf("Value 2: count = %d, expected ≈ %.0f (deviation: %.2f%%)",
				valueCounts[2], expectedValue2Count, math.Abs(float64(valueCounts[2])-expectedValue2Count)/expectedValue2Count*100)
		}

		// Verify distribution of instances matches weights
		// Each instance of value 1 should have 1.0/5.0 = 0.2 chance
		// Each instance of value 2 should have 1.5/5.0 = 0.3 chance
		for key, count := range instanceCounts {
			var expectedCount float64
			if key.value == 1 {
				expectedCount = float64(iterations) * 0.2 // 1.0/5.0
			} else {
				expectedCount = float64(iterations) * 0.3 // 1.5/5.0
			}
			if math.Abs(float64(count)-expectedCount) > tolerance {
				t.Errorf("Value %d Instance %d: count = %d, expected ≈ %.0f (deviation: %.2f%%)",
					key.value, key.instance, count, expectedCount, math.Abs(float64(count)-expectedCount)/expectedCount*100)
			}
		}
	})

	t.Run("No duplicate instances", func(t *testing.T) {
		// Create items with multiple instances
		items := []Item{
			&testItem{value: 1, weight: 1.0, supply: 2},
			&testItem{value: 2, weight: 1.0, supply: 2},
		}

		// Select all instances
		r := NewRandomness(BetaBytes("test"))
		results, err := r.Selection(SelectionConfig{Items: items, Count: 4})
		if err != nil {
			t.Fatalf("Selection() error = %v", err)
		}

		// Track which instances were selected
		selectedInstances := make(map[int]map[int]bool)
		for _, result := range results {
			item := result.Get().(*testItem)
			if selectedInstances[item.value] == nil {
				selectedInstances[item.value] = make(map[int]bool)
			}
			instance := result.Instance()
			if selectedInstances[item.value][instance] {
				t.Errorf("Instance %d of item %d was selected multiple times", instance, item.value)
			}
			selectedInstances[item.value][instance] = true
		}

		// Verify all instances were selected exactly once
		for value, instances := range selectedInstances {
			item := items[value-1].(*testItem)
			for i := 1; i <= item.supply; i++ {
				if !instances[i] {
					t.Errorf("Instance %d of item %d was not selected", i, value)
				}
			}
		}
	})
}

func TestResetSelection(t *testing.T) {
	// Create items with multiple instances
	items := []Item{
		&testItem{value: 1, weight: 1.0, supply: 3},
		&testItem{value: 2, weight: 1.0, supply: 2},
	}

	// Select some instances
	r := NewRandomness(BetaBytes("test"))
	c := SelectionConfig{Items: items, Count: 5}
	results, err := r.Selection(c)
	if err != nil {
		t.Fatalf("Selection() error = %v", err)
	}

	// Reset the selection
	c.Reset()

	// Verify we can select all instances again
	results, err = r.Selection(c)
	if err != nil {
		t.Fatalf("Selection() error after reset = %v", err)
	}

	// Track which instances were selected
	selectedInstances := make(map[int]map[int]bool)
	for _, result := range results {
		item := result.Get().(*testItem)
		if selectedInstances[item.value] == nil {
			selectedInstances[item.value] = make(map[int]bool)
		}
		instance := result.Instance()
		if selectedInstances[item.value][instance] {
			t.Errorf("Instance %d of item %d was selected multiple times after reset", instance, item.value)
		}
		selectedInstances[item.value][instance] = true
	}

	// Verify all instances were selected exactly once
	for value, instances := range selectedInstances {
		item := items[value-1].(*testItem)
		for i := 1; i <= item.supply; i++ {
			if !instances[i] {
				t.Errorf("Instance %d of item %d was not selected after reset", i, value)
			}
		}
	}
}
