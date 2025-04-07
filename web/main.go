//go:build js && wasm

package main

import (
	"fmt"
	"runtime"
	"syscall/js"

	"github.com/revision-3/randomness"
)

// RandomnessWrapper wraps the Randomness interface for JavaScript
type RandomnessWrapper struct {
	r randomness.Randomness
}

// NewRandomnessWrapper creates a new RandomnessWrapper
func NewRandomnessWrapper(entropy randomness.BetaBytes) *RandomnessWrapper {
	return &RandomnessWrapper{
		r: randomness.NewRandomness(entropy),
	}
}

type Result = map[string]any

func ErrResult(err string) Result {
	return Result{
		"value": nil,
		"err":   err,
	}
}

func ValueResult(value any) Result {
	return Result{
		"value": value,
		"err":   nil,
	}
}

func ValueOf(value Result) js.Value {
	var v map[string]any = value
	return js.ValueOf(v)
}

func panicSafe[T any](fn func() T, defaultVal T) (ret T) {
	defer func() {
		if r := recover(); r != nil {
			ret = defaultVal
		}
	}()
	return fn()
}

// panicHandler wraps a function that might panic and converts the panic to a JavaScript error
func panicHandler(fn func() Result) (ret js.Value) {
	defer func() {
		if r := recover(); r != nil {
			// Convert the panic value to a string
			fmt.Printf("caught panic: %v\n", r)
			var errMsg string
			switch v := r.(type) {
			case string:
				errMsg = v
			case error:
				errMsg = v.Error()
			default:
				errMsg = fmt.Sprintf("%v", v)
			}
			ret = ValueOf(ErrResult(errMsg))
		}
	}()
	return ValueOf(fn())
}

// Probability returns a random number in (0.0, 1.0]
func (w *RandomnessWrapper) Probability() (float64, error) {
	return w.r.Probability()
}

// Bits returns n random bits
func (w *RandomnessWrapper) Bits(n int) ([]bool, error) {
	return w.r.Bits(n)
}

// Bytes returns n random bytes
func (w *RandomnessWrapper) Bytes(n int) ([]byte, error) {
	return w.r.Bytes(n)
}

// Uint64 returns a random uint64
func (w *RandomnessWrapper) Uint64() (uint64, error) {
	return w.r.Uint64()
}

// Uint32 returns a random uint32
func (w *RandomnessWrapper) Uint32() (uint32, error) {
	return w.r.Uint32()
}

// Uint16 returns a random uint16
func (w *RandomnessWrapper) Uint16() (uint16, error) {
	return w.r.Uint16()
}

// Uint8 returns a random uint8
func (w *RandomnessWrapper) Uint8() (uint8, error) {
	return w.r.Uint8()
}

// Int64 returns a random int64
func (w *RandomnessWrapper) Int64() (int64, error) {
	return w.r.Int64()
}

// Int32 returns a random int32
func (w *RandomnessWrapper) Int32() (int32, error) {
	return w.r.Int32()
}

// Int16 returns a random int16
func (w *RandomnessWrapper) Int16() (int16, error) {
	return w.r.Int16()
}

// Int8 returns a random int8
func (w *RandomnessWrapper) Int8() (int8, error) {
	return w.r.Int8()
}

// Float64 returns a random float64
func (w *RandomnessWrapper) Float64() (float64, error) {
	return w.r.Float64()
}

// Float32 returns a random float32
func (w *RandomnessWrapper) Float32() (float32, error) {
	return w.r.Float32()
}

// PickDistinct returns n unique random integers in [0, magnitude)
func (w *RandomnessWrapper) PickDistinct(n int, magnitude int) ([]int, error) {
	return w.r.PickDistinct(n, magnitude)
}

// Pick returns n random integers in [0, magnitude) (may include duplicates)
func (w *RandomnessWrapper) Pick(n int, magnitude int) ([]int, error) {
	return w.r.Pick(n, magnitude)
}

// Selection performs a weighted random selection of items based on their weights and supplies
func (w *RandomnessWrapper) Selection(items []randomness.Item, count int) ([]randomness.SelectionResult, error) {
	cfg := randomness.SelectionConfig{
		Items: items,
		Count: count,
	}
	return w.r.Selection(cfg)
}

// NewRandomness creates a new Randomness instance
func NewRandomness(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(ErrResult("error: entropy string required"))
	}

	hexEntropy := args[0].String()
	entropy, err := randomness.BetaBytesFromHex(hexEntropy)
	if err != nil {
		return js.ValueOf(ErrResult(fmt.Sprintf("error: invalid hex string: %v", err)))
	}
	wrapper := NewRandomnessWrapper(entropy)

	var m runtime.MemStats
	stats := map[string]any{}
	mr := map[string]any{
		"value": stats,
	}

	// Create a JavaScript object with all the methods
	lib := Result{}

	lib["name"] = js.ValueOf(randomness.Name)
	lib["version"] = js.ValueOf(randomness.Version)

	lib["memStats"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		runtime.GC()
		runtime.ReadMemStats(&m)
		stats["alloc"] = m.Alloc
		stats["totalAlloc"] = m.TotalAlloc
		stats["sys"] = m.Sys
		stats["heapAlloc"] = m.HeapAlloc
		stats["heapSys"] = m.HeapSys
		stats["heapIdle"] = m.HeapIdle
		stats["heapInuse"] = m.HeapInuse
		stats["heapReleased"] = m.HeapReleased
		stats["mallocs"] = m.Mallocs
		stats["frees"] = m.Frees
		stats["gcSys"] = m.GCSys
		return js.ValueOf(mr)
	})

	lib["probability"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Probability()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["bits"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return js.ValueOf(ErrResult("error: n parameter required"))
		}
		n := args[0].Int()
		return panicHandler(func() Result {
			bits, err := wrapper.Bits(n)
			if err != nil {
				return ErrResult(err.Error())
			}
			result := make([]any, len(bits))
			for i, bit := range bits {
				result[i] = bit
			}
			return ValueResult(result)
		})
	})
	lib["bytes"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return js.ValueOf(ErrResult("error: n parameter required"))
		}
		n := args[0].Int()
		return panicHandler(func() Result {
			bytes, err := wrapper.Bytes(n)
			if err != nil {
				return ErrResult(err.Error())
			}
			result := make([]any, len(bytes))
			for i, b := range bytes {
				result[i] = b
			}
			return ValueResult(result)
		})
	})
	lib["uint64"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Uint64()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["uint32"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Uint32()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["uint16"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Uint16()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["uint8"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Uint8()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["int64"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Int64()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["int32"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Int32()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["int16"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Int16()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["int8"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Int8()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["float64"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Float64()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["float32"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		return panicHandler(func() Result {
			value, err := wrapper.Float32()
			if err != nil {
				return ErrResult(err.Error())
			}
			return ValueResult(value)
		})
	})
	lib["selection"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 2 {
			return js.ValueOf(ErrResult("error: items and count parameters required"))
		}
		itemsArg := args[0]
		count := args[1].Int()

		// Convert JavaScript array of items to Go slice
		items := make([]randomness.Item, itemsArg.Length())
		for i := 0; i < itemsArg.Length(); i++ {
			item := itemsArg.Index(i)
			if item.Truthy() {
				genericItem := randomness.GenericItem[string]{
					BaseItem: randomness.NewBaseItem(
						panicSafe(func() float64 { return item.Get("weight").Float() }, 1.0),
						panicSafe(func() int { return item.Get("supply").Int() }, -1),
					),
					Value: item.Get("value").String(),
				}
				items[i] = genericItem
			}
		}

		return panicHandler(func() Result {
			results, err := wrapper.Selection(items, count)
			if err != nil {
				return ErrResult(err.Error())
			}
			jsResults := make([]any, len(results))
			for i, result := range results {
				item := result.Get().(randomness.GenericItem[string])
				jsResults[i] = map[string]any{
					"value":    item.Value,
					"weight":   result.Weight(),
					"supply":   result.Supply(),
					"instance": result.Instance(),
				}
			}
			return ValueResult(jsResults)
		})
	})
	lib["pickDistinct"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 2 {
			return js.ValueOf(ErrResult("error: n and magnitude parameters required"))
		}
		n := args[0].Int()
		magnitude := args[1].Int()
		return panicHandler(func() Result {
			selected, err := wrapper.PickDistinct(n, magnitude)
			if err != nil {
				return ErrResult(err.Error())
			}
			result := make([]any, len(selected))
			for i, s := range selected {
				result[i] = s
			}
			return ValueResult(result)
		})
	})
	lib["pick"] = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 2 {
			return js.ValueOf(ErrResult("error: n and magnitude parameters required"))
		}
		n := args[0].Int()
		magnitude := args[1].Int()
		return panicHandler(func() Result {
			selected, err := wrapper.Pick(n, magnitude)
			if err != nil {
				return ErrResult(err.Error())
			}
			result := make([]any, len(selected))
			for i, s := range selected {
				result[i] = s
			}
			return ValueResult(result)
		})
	})

	return ValueOf(lib)
}

func main() {
	wait := make(chan struct{}, 0)
	// Console log that we're in the main function
	fmt.Println("Loading Randomness")
	js.Global().Set("Randomness", js.FuncOf(NewRandomness))
	<-wait
}
