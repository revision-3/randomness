# Randomness WASM Package

This package provides a WebAssembly implementation of the Randomness library for use in web browsers.

## Installation

Download the latest release from the GitHub releases page. The release package includes:
- `randomness.wasm` - The WebAssembly binary
- `wasm_exec.js` - The Go WebAssembly support file
- `pubkey.asc` - The public key for verifying releases
- `index.html` - A demo implementation
- `README.md` - This documentation

## Basic Usage

1. First, import the required files in your HTML:

```html
<script src="wasm_exec.js"></script>
<script>
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("randomness.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
        console.log("WASM loaded successfully");
    }).catch((err) => {
        console.error("Error loading WASM:", err);
    });
</script>
```

2. Create a Randomness instance:

```javascript
// Generate a random seed (32 bytes = 64 hex characters)
const array = new Uint8Array(32);
crypto.getRandomValues(array);
const seed = Array.from(array)
    .map(b => b.toString(16).padStart(2, '0'))
    .join('');

// Create the Randomness instance with the hex string
// The hex string will be automatically converted to BetaBytes
const result = Randomness(seed);
if (result.err) {
    console.error("Error creating Randomness instance:", result.err);
    return;
}
const randomness = result;
```

Note: The entropy string must be a valid hex string. The Randomness constructor will automatically convert it to BetaBytes.

## Available Methods

All methods return an object with the following structure:
```javascript
{
    value: any,  // The result value
    err: string  // Error message if any, null if successful
}
```

### Basic Random Number Generation

```javascript
// Get a random number in (0.0, 1.0]
const result = randomness.probability();

// Get n random bits
const result = randomness.bits(n);

// Get n random bytes
const result = randomness.bytes(n);
```

### Integer Types

```javascript
// Get random integers of various sizes
const result = randomness.uint64();  // Random uint64
const result = randomness.uint32();  // Random uint32
const result = randomness.uint16();  // Random uint16
const result = randomness.uint8();   // Random uint8
const result = randomness.int64();   // Random int64
const result = randomness.int32();   // Random int32
const result = randomness.int16();   // Random int16
const result = randomness.int8();    // Random int8
```

### Floating Point Types

```javascript
// Get random floating point numbers
const result = randomness.float64();  // Random float64
const result = randomness.float32();  // Random float32
```

### Selection Methods

```javascript
// Pick n unique random integers in [0, magnitude)
const result = randomness.pickDistinct(n, magnitude);

// Pick n random integers in [0, magnitude) (may include duplicates)
const result = randomness.pick(n, magnitude);

// Perform weighted random selection
const items = [
    { value: "apple", weight: 1.0, supply: 2 },
    { value: "banana", weight: 1.0, supply: 3 },
    { value: "orange", weight: 1.5, supply: -1 },  // -1 means infinite supply
    { value: "grape", weight: 0.5, supply: 4 },
    { value: "kiwi", weight: 1.0, supply: 2 }
];
const result = randomness.selection(items, count);
```

### Utility Methods

```javascript
// Get the package name/version
const name = randomness.name;
const version = randomness.version;

// Get memory statistics
const stats = randomness.memStats();
```

## Verifying Releases

To verify the integrity of the release files:

1. Import the public key from a keyserver:
   ```bash
   gpg --recv-keys 76B5D604BF678DB4
   ```

2. Verify the files:
   ```bash
   gpg --verify randomness.wasm.asc randomness.wasm
   gpg --verify wasm_exec.js.asc wasm_exec.js
   ```

## Example Implementation

A complete example implementation is included in `index.html`. This demonstrates all available methods and includes a user interface for testing the functionality.

## Error Handling

All methods return an object with `value` and `err` properties. Always check the `err` property before using the `value`:

```javascript
const result = randomness.probability();
if (result.err) {
    console.error("Error:", result.err);
    return;
}
const value = result.value;
```

## Memory Management

The WASM implementation includes memory statistics that can be monitored:

```javascript
const stats = randomness.memStats();
if (stats.err) {
    console.error("Error getting memory stats:", stats.err);
    return;
}
console.log("Memory Statistics:", stats.value);
```

## License

This package is licensed solely to be used for the purpose of verifying the fairness of games. It may not be used in any other context or for any commercial purposes without express permission.
See the LICENSE file for full details. 