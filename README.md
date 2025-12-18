# Sentence Cipher | sentence-cipher

A Go library that encodes arbitrary binary data into natural-looking English sentences. Hide your secret messages in plain sight!

## Overview

Sentence Cipher transforms any byte sequence into grammatically correct English sentences using a workplace/office theme vocabulary. Each sentence follows predictable patterns that allow perfect round-trip encoding and decoding.

### Example

```
Input:  "Hello"
Output: "ruth trains isabella prints. carl cleans daily."
```

The encoded text looks like normal English sentences about office activities, making it suitable for steganographic purposes.

## How It Works

### Encoding Scheme

The cipher uses three word lists, each containing exactly 256 words (one for each possible byte value):

| Word List | Purpose | Examples |
|-----------|---------|----------|
| **Names** (256) | Subjects & Indirect Objects | adam, alex, alice, amanda... |
| **Verbs** (256) | Actions | helps, assists, supports, guides... |
| **Objects** (256) | Direct Objects | reports, documents, files, records... |

### Sentence Patterns

Data is encoded using three different sentence patterns depending on how many bytes remain:

#### Pattern 1: Full Sentence (3 bytes)
```
[Subject] [Verb] [IndirectObject] [Object].
```
- **Byte 1** → Subject (name index 0-255)
- **Byte 2** → Verb (verb index 0-255)
- **Byte 3** → Object (object index 0-255)
- **IndirectObject** = derived from `(byte1 + byte2) % 256` for natural flow

#### Pattern 2: Short Sentence (2 bytes)
```
[Subject] [Verb] daily.
```
- **Byte 1** → Subject (name index)
- **Byte 2** → Verb (verb index)

#### Pattern 3: Minimal Sentence (1 byte)
```
[Subject] works.
```
- **Byte 1** → Subject (name index)

### Key-Based Encryption

For added security, you can provide an encryption key that shuffles the word lists deterministically using SHA-256:

```go
// Without key - uses default word order
cipher := sentencecipher.NewDefaultCipher()

// With key - word lists are shuffled based on the key
cipher := sentencecipher.NewCipher("my-secret-key")
```

The same key must be used for both encoding and decoding.

## Installation

```bash
go get github.com/kittizz/sentence-cipher
```

## Usage

### As a Library

```go
package main

import (
    "fmt"
    sentencecipher "github.com/kittizz/sentence-cipher"
)

func main() {
    // Basic encoding/decoding
    message := "Hello World"
    encoded := sentencecipher.EncodeString(message)
    fmt.Println("Encoded:", encoded)
    
    decoded, err := sentencecipher.DecodeString(encoded)
    if err != nil {
        panic(err)
    }
    fmt.Println("Decoded:", decoded)
    
    // With encryption key
    cipher := sentencecipher.NewCipher("my-secret-key")
    secureEncoded := cipher.EncodeString(message)
    secureDecoded, _ := cipher.DecodeString(secureEncoded)
    
    // Natural mode (more varied sentences)
    naturalEncoded := sentencecipher.EncodeNatural([]byte(message))
    naturalDecoded, _ := sentencecipher.DecodeNatural(naturalEncoded)
}
```

### Natural Encoding Mode

Natural mode adds variety to sentences using time prefixes:

```
Standard: "ruth trains isabella prints. carl cleans hannah reports."
Natural:  "ruth trains isabella prints. Today carl cleans hannah reports. Later tom helps mary files."
```

Supported prefixes: `Today`, `Later`, `Then`, `Now`

### Command Line Tool

Build and install the CLI:

```bash
go build -o grammarcipher ./cmd/sentencecipher
```

#### CLI Usage

```bash
# Encode text
grammarcipher "Hello World"

# Decode text
grammarcipher -d "ruth trains isabella prints."

# Use encryption key
grammarcipher -k "my-secret-key" "Secret message"
grammarcipher -d -k "my-secret-key" "encoded text here"

# Natural mode
grammarcipher -n "Hello World"
grammarcipher -d -n "Today ruth trains isabella prints."

# File I/O
grammarcipher -i input.txt -o encoded.txt
grammarcipher -d -i encoded.txt -o decoded.txt

# Pipe from stdin
echo "Secret" | grammarcipher
echo "encoded text" | grammarcipher -d
```

#### CLI Options

| Flag | Description |
|------|-------------|
| `-d` | Decode mode (default is encode) |
| `-n` | Use natural encoding (more varied sentences) |
| `-k KEY` | Encryption key (shuffles word lists) |
| `-i FILE` | Read input from file |
| `-o FILE` | Write output to file |
| `-v` | Show version |
| `-h` | Show help |

## API Reference

### Types

```go
type Cipher struct {
    // Contains shuffled word lists
}
```

### Functions

#### Cipher Constructors

| Function | Description |
|----------|-------------|
| `NewCipher(key string) *Cipher` | Create cipher with key-shuffled word lists |
| `NewDefaultCipher() *Cipher` | Create cipher with default word lists |

#### Cipher Methods

| Method | Description |
|--------|-------------|
| `Encode(data []byte) string` | Encode bytes to sentences |
| `Decode(encoded string) ([]byte, error)` | Decode sentences to bytes |
| `EncodeString(s string) string` | Encode string to sentences |
| `DecodeString(encoded string) (string, error)` | Decode sentences to string |
| `EncodeNatural(data []byte) string` | Encode with varied sentence templates |
| `DecodeNatural(encoded string) ([]byte, error)` | Decode natural-encoded text |

#### Package-Level Functions (Backward Compatibility)

| Function | Description |
|----------|-------------|
| `Encode(data []byte) string` | Encode using default cipher |
| `Decode(encoded string) ([]byte, error)` | Decode using default cipher |
| `EncodeString(s string) string` | Encode string using default cipher |
| `DecodeString(encoded string) (string, error)` | Decode string using default cipher |
| `EncodeNatural(data []byte) string` | Natural encode using default cipher |
| `DecodeNatural(encoded string) ([]byte, error)` | Natural decode using default cipher |

## Technical Details

### Encoding Efficiency

| Input Size | Output Pattern | Bytes Encoded |
|------------|----------------|---------------|
| 3+ bytes | Full sentence (4 words) | 3 bytes |
| 2 bytes | Short sentence (3 words) | 2 bytes |
| 1 byte | Minimal sentence (2 words) | 1 byte |

**Expansion ratio:** Approximately 15-20x (varies based on word lengths)

### Key Derivation

When a key is provided:
1. SHA-256 hash is computed from the key
2. First 8 bytes of hash become the seed
3. Each word list is shuffled using Fisher-Yates with a unique seed offset

```go
hash := sha256.Sum256([]byte(key))
seed := int64(binary.BigEndian.Uint64(hash[:8]))

names   = shuffleWithSeed(defaultNames, seed)
verbs   = shuffleWithSeed(defaultVerbs, seed+1)
objects = shuffleWithSeed(defaultObjects, seed+2)
```

### Word List Statistics

- **Total vocabulary:** 768 unique words
- **Names:** 256 common first names (gender-neutral selection)
- **Verbs:** 256 workplace action verbs
- **Objects:** 256 workplace items and concepts

## Testing

```bash
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

## Use Cases

- **Steganography:** Hide messages in seemingly innocent text
- **Data obfuscation:** Make binary data look like natural language
- **CTF challenges:** Create encoding puzzles
- **Educational:** Demonstrate encoding concepts

## Limitations

- Output is ~15-20x larger than input
- Not cryptographically secure without a key
- Encoded text may look repetitive for long messages
- Word lists are fixed (office/workplace theme)

## License

MIT License
