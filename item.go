package randomness

// Item represents a selectable item with a weight and supply
type Item interface {
	Weight() float64
	Supply() int
}

// BaseItem provides a base implementation of the Item interface
type BaseItem struct {
	weight float64
	supply int
}

// GenericItem is a generic implementation of the Item interface as a generic type
type GenericItem[T any] struct {
	BaseItem
	Value T `json:"value"`
}

// Weight returns the weight of the item
func (i BaseItem) Weight() float64 {
	return i.weight
}

// Supply returns the supply of the item
func (i BaseItem) Supply() int {
	return i.supply
}

func SingleItem() BaseItem {
	return BaseItem{
		weight: 1,
		supply: 1,
	}
}

func InfiniteItem() BaseItem {
	return BaseItem{
		weight: 1,
		supply: -1,
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
