package randomness

import (
	"fmt"
)

type SelectionResult interface {
	Item
	Instance() int
	Get() Item
	Fraction() float64 // Returns the fractional position within the instance (0-1)
}

// SelectionResult represents a single selected item with its instance number
type selectionResult struct {
	Item             // The selected item
	instance int     // The instance number of the selected item (for finite supplies)
	fraction float64 // The fractional position within the instance (0-1)
}

func (s *selectionResult) Instance() int {
	return s.instance
}

func (s *selectionResult) Get() Item {
	return s.Item
}

func (s *selectionResult) Fraction() float64 {
	return s.fraction
}

// SelectionConfig represents the configuration for a selection operation.
// It contains the items to select from, the number of items to select,
// and tracks the state of each item's instances.
type SelectionConfig struct {
	Items      []Item
	Count      int
	ItemStates []itemState
}

// itemState tracks the state of an item's instances during selection.
// It maintains the original supply, remaining supply, and which instances have been used.
type itemState struct {
	Item            Item
	OriginalSupply  int
	RemainingSupply int
	UsedInstances   map[int]bool
	IsConsumed      bool // Whether this item's supply is consumed (finite) or not (infinite)
}

// Reset resets the selection state, allowing all instances to be selected again.
// For infinite supply items, this has no effect since they can always be selected.
func (c *SelectionConfig) Reset() {
	for i := range c.ItemStates {
		c.ItemStates[i].RemainingSupply = c.ItemStates[i].OriginalSupply
		c.ItemStates[i].UsedInstances = make(map[int]bool)
	}
}

// Selection performs a weighted random selection of items based on their weights and supplies.
// Each instance of an item is treated as a separate selectable item with the same weight as its parent.
// For example, if there are 2 apples (weight 1.0) and 2 oranges (weight 1.0), each instance has a 1/4 chance
// of being selected.
// If there's then a single grapefruit added (weight 1.5), its instances each have a 1.5/5.5 chance and the apples and
// oranges each have a 1/5.5 chance.
// Negative supplies represent infinite items where the magnitude affects their relative weight.
// For example, if there are 3 finite apples and 3 infinite oranges (both weight 1.0), the first selection
// has a 3/6 chance of being an orange or apple. If an apple is selected, the next selection has a 3/5
// chance of being an orange since there are 2 apples and 3 infinite oranges remaining.
func (r *randomness) Selection(cfg SelectionConfig) ([]SelectionResult, error) {
	// Validate the configuration
	if err := ValidateSelectionConfig(cfg); err != nil {
		return nil, err
	}

	// Initialize results slice
	results := make([]SelectionResult, 0, cfg.Count)

	// Initialize item states if not already done
	if cfg.ItemStates == nil {
		cfg.ItemStates = make([]itemState, len(cfg.Items))
		for i, item := range cfg.Items {
			supply := item.Supply()
			originalSupply := supply
			if supply < 0 {
				originalSupply = -supply // Convert negative to positive for tracking
			}
			cfg.ItemStates[i] = itemState{
				Item:            item,
				OriginalSupply:  originalSupply,
				RemainingSupply: originalSupply,
				UsedInstances:   make(map[int]bool),
				IsConsumed:      supply >= 0,
			}
		}
	}

	// Perform the requested number of selections
	for range cfg.Count {
		// Calculate total weight of all available instances
		totalWeight := 0.0
		for i := range cfg.ItemStates {
			state := &cfg.ItemStates[i]
			if !state.IsConsumed {
				// For infinite supply items, their weight is multiplied by their supply magnitude
				totalWeight += state.Item.Weight() * float64(state.OriginalSupply)
			} else if state.RemainingSupply > 0 {
				// For finite supply items, add weight for each unused instance
				for instance := 1; instance <= state.OriginalSupply; instance++ {
					if !state.UsedInstances[instance] {
						totalWeight += state.Item.Weight()
					}
				}
			}
		}

		// Check if there are any instances available to select
		if totalWeight == 0 {
			return nil, fmt.Errorf("no items remaining with non-zero supply")
		}

		// Generate random value between 0 and 1
		prob, err := r.Probability()
		if err != nil {
			return nil, err
		}

		// Select an instance using normalized weights
		accumulatedProb := 0.0
		var selectedState *itemState
		var selectedInstance int
		var fractionalPos float64

		for i := range cfg.ItemStates {
			state := &cfg.ItemStates[i]
			if !state.IsConsumed {
				// For infinite supply items, their weight is multiplied by their supply magnitude
				weight := state.Item.Weight() * float64(state.OriginalSupply)
				accumulatedProb += weight / totalWeight
				if prob <= accumulatedProb {
					selectedState = state
					// For infinite supply items, we can reuse any instance number
					selectedInstance = 1
					// Calculate fractional position
					if i > 0 {
						prevAccumulated := accumulatedProb - weight/totalWeight
						fractionalPos = (prob - prevAccumulated) / (weight / totalWeight)
					} else {
						fractionalPos = prob / (weight / totalWeight)
					}
					break
				}
			} else if state.RemainingSupply > 0 {
				// For finite supply items
				for instance := 1; instance <= state.OriginalSupply; instance++ {
					if !state.UsedInstances[instance] {
						instanceWeight := state.Item.Weight() / totalWeight
						accumulatedProb += instanceWeight
						if prob <= accumulatedProb {
							selectedState = state
							selectedInstance = instance
							// Calculate fractional position
							fractionalPos = (prob - (accumulatedProb - instanceWeight)) / instanceWeight
							break
						}
					}
				}
				if selectedState != nil {
					break
				}
			}
		}

		// Find the next available instance number for the selected item
		if selectedState.IsConsumed {
			selectedInstance = 0
			for instance := 1; instance <= selectedState.OriginalSupply; instance++ {
				if !selectedState.UsedInstances[instance] {
					selectedInstance = instance
					break
				}
			}
			if selectedInstance == 0 {
				return nil, fmt.Errorf("no available instances for selected item")
			}
		}

		// Mark the selected instance as used (only for finite supply items)
		if selectedState.IsConsumed {
			selectedState.UsedInstances[selectedInstance] = true
			selectedState.RemainingSupply--
		}

		// Add the selected instance to the results
		results = append(results, &selectionResult{
			Item:     selectedState.Item,
			instance: selectedInstance,
			fraction: fractionalPos,
		})
	}

	return results, nil
}

// ValidateSelectionConfig validates a selection configuration
func ValidateSelectionConfig(cfg SelectionConfig) error {
	if len(cfg.Items) == 0 {
		return fmt.Errorf("no items to select from")
	}

	if cfg.Count <= 0 {
		return fmt.Errorf("count must be positive")
	}

	// Check for negative weights and supplies
	hasFinite := false
	hasInfinite := false
	for _, item := range cfg.Items {
		if item.Weight() < 0 {
			return fmt.Errorf("weights must be non-negative")
		}
		if item.Supply() < 0 {
			hasInfinite = true
		} else {
			hasFinite = true
		}
	}

	// Only check total supply if there are no infinite items
	if !hasInfinite && hasFinite && cfg.Count > 0 {
		totalSupply := 0
		for _, item := range cfg.Items {
			if item.Supply() >= 0 {
				totalSupply += item.Supply()
			}
		}
		if totalSupply < cfg.Count {
			return fmt.Errorf("total supply (%d) is less than requested count (%d)", totalSupply, cfg.Count)
		}
	}

	return nil
}
