package randomness

import "encoding/json"

type Itemer interface {
	Weight() float64
	Supply() int
	Any() any
}

type TypedItemer[T any] interface {
	Itemer
	Value() T
}

// Item represents a selectable item with a weight and supply
type Item interface {
	Weight() float64
	Supply() int
	Any() any
}

// BaseItem provides a base implementation of the Item interface
type BaseItem struct {
	weight float64
	supply int
}

// GenericItem is a generic implementation of the Item interface as a generic type
type GenericItem[T any] struct {
	BaseItem
	value T
}

func NewGenericItem[T any](value T, weight float64, supply int) *GenericItem[T] {
	return &GenericItem[T]{
		BaseItem: BaseItem{weight: weight, supply: supply},
		value:    value,
	}
}

func (i *GenericItem[T]) Any() any {
	return i.value
}

func (i *GenericItem[T]) Value() T {
	return i.value
}

func (i *GenericItem[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"value":  i.value,
		"weight": i.weight,
		"supply": i.supply,
	})
}

func (i *GenericItem[T]) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	i.weight = 1.0
	i.supply = 1

	if value, ok := m["value"].(T); ok {
		i.value = value
	}

	if weight, ok := m["weight"].(float64); ok {
		i.weight = weight
	}

	if supply, ok := m["supply"].(int); ok {
		i.supply = supply
	}

	return nil
}

var _ TypedItemer[any] = &GenericItem[any]{}

// Weight returns the weight of the item
func (i BaseItem) Weight() float64 {
	return i.weight
}

// Supply returns the supply of the item
func (i BaseItem) Supply() int {
	return i.supply
}

func (i BaseItem) Value() any {
	return nil
}

func SingleItem() BaseItem {
	return BaseItem{
		weight: 1,
		supply: 1,
	}
}

func InfiniteSingle() BaseItem {
	return BaseItem{
		weight: 1,
		supply: -1,
	}
}

func InfiniteItem(supply int) BaseItem {
	return BaseItem{
		weight: 1,
		supply: -supply,
	}
}

func FiniteItem(supply int) BaseItem {
	return BaseItem{
		weight: 1,
		supply: supply,
	}
}

func WeightedSingleItem(weight float64) BaseItem {
	return BaseItem{
		weight: weight,
		supply: 1,
	}
}

func WeightedInfiniteItem(weight float64) BaseItem {
	return BaseItem{
		weight: weight,
		supply: -1,
	}
}

func WeightedFiniteItem(weight float64, supply int) BaseItem {
	return BaseItem{
		weight: weight,
		supply: supply,
	}
}

func NewBaseItem(weight float64, supply int) BaseItem {
	return BaseItem{
		weight: weight,
		supply: supply,
	}
}
