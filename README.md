# Randomness Package

A robust implementation for generating random numbers and performing weighted random selections, designed for verifying the fairness of games.

## Overview
The `randomness` package provides a comprehensive set of tools for generating random values and performing weighted random selections. It's designed to be used for verifying the fairness of games, with a focus on cryptographic security and precise probability control.

## Core Components

### 1. Randomness Interface
The main interface that provides various methods for generating random values:

```go
type Randomness interface {
    Probability() float64        // Returns a random number in (0.0, 1.0]
    Bits(n int) BitArray        // Returns n random bits
    Bytes(n int) []byte         // Returns n random bytes
    Uint64() uint64             // Returns a random uint64
    Uint32() uint32             // Returns a random uint32
    Uint16() uint16             // Returns a random uint16
    Uint8() uint8               // Returns a random uint8
    Int64() int64               // Returns a random int64
    Int32() int32               // Returns a random int32
    Int16() int16               // Returns a random int16
    Int8() int8                 // Returns a random int8
    Float64() float64           // Returns a random float64
    Float32() float32           // Returns a random float32
    Select(n int, magnitude int) []int  // Returns n unique random integers in [0, magnitude)
    Numbers(count, magnitude int) Numbers  // Returns a Numbers interface for reading random numbers
    Selection(cfg SelectionConfig) ([]SelectionResult, error)  // Performs weighted random selection
}
```

### 2. Initialization
To create a new Randomness instance:

```go
// Create from a hex string
entropy, err := BetaBytesFromHex("deadbeef")
if err != nil {
    // Handle error
}
r := NewRandomness(entropy)

// Or use the convenience function for hex strings
r := NewRandomness(MustBetaBytesFromHex("deadbeef"))

// Or create from raw bytes
r := NewRandomness(BetaBytes([]byte("raw bytes")))
```

Where the entropy source is provided as `BetaBytes`. The implementation uses SHA-512 to amplify this entropy as needed.

## Usage Examples

### 1. Basic Random Number Generation

```go
r := NewRandomness("initial entropy")

// Generate random numbers
randomUint64 := r.Uint64()
randomFloat := r.Probability()  // Returns a number in (0.0, 1.0]
randomBytes := r.Bytes(16)      // Get 16 random bytes
```

### 2. Weighted Random Selection

```go
// Create items with different weights and supplies
items := []Item{
    WeightedFiniteItem(1.0, 3),  // 3 items with weight 1.0
    WeightedFiniteItem(1.5, 2),  // 2 items with weight 1.5
    WeightedFiniteItem(2.0, 1),  // 1 item with weight 2.0
}

// Create selection configuration
cfg := SelectionConfig{
    Items: items,
    Count: 3,  // Select 3 items
}

// Perform selection
results, err := r.Selection(cfg)
if err != nil {
    // Handle error
}

// Process results
for _, result := range results {
    item := result.Get()
    instance := result.Instance()
    // Use the selected item and its instance number
}
```

### 3. Generating Unique Random Numbers

```go
// Generate 5 unique random numbers in the range [0, 10)
uniqueNumbers := r.Select(5, 10)

// Or use the Numbers interface for more control
numbers := r.Numbers(5, 10)
for i := 0; i < 5; i++ {
    randomNum := numbers.Read(10)  // Get a number in [0, 10)
}
```

### 4. Working with Bits

```go
// Get 8 random bits
bits := r.Bits(8)
for i := 0; i < bits.Length(); i++ {
    if bits.Get(i) {
        // Bit is 1
    } else {
        // Bit is 0
    }
}
```

## Important Notes

1. **Entropy Amplification**: The implementation automatically amplifies entropy using SHA-512 when needed, ensuring a continuous supply of random values.

2. **Probability Range**: The `Probability()` method returns values in the range (0.0, 1.0], never returning 0.0 to avoid potential issues in probability calculations.

3. **Selection Features**:
   - Supports both finite and infinite supply items
   - Handles weighted selection with precise probability control
   - Tracks instance numbers for finite supply items
   - Prevents duplicate selection of instances
   - Can be reset to allow re-selection of items

4. **Thread Safety**: The implementation is not thread-safe by default. If you need thread safety, you should implement appropriate synchronization mechanisms.

5. **License Restrictions**: This package is licensed specifically for verifying the fairness of games and may not be used for commercial purposes.

## Best Practices

1. Always use a cryptographically secure source for the initial entropy string when creating a new Randomness instance.

2. When using weighted selection:
   - Keep weights non-negative
   - Don't mix finite and infinite supply items in the same selection
   - Ensure total supply is sufficient for finite items

3. For probability-based decisions, use the `Probability()` method rather than converting from other random number types.

4. When selecting multiple items, use the `Selection()` method rather than multiple independent selections to ensure proper distribution.

## License
This package is licensed solely to be used for the purpose of verifying the fairness of games. It may not be used in any other context or for any commercial purposes. See the LICENSE file for full details.

## Verifying Releases

All releases are signed with GPG. To verify a release:

1. Import the public key:
   ```bash
   gpg --import pubkey.asc
   ```

2. Verify a file's signature:
   ```bash
   gpg --verify file.asc file
   ```

The public key used for signing releases is available in `pubkey.asc` in this repository. 