# Sentence Cipher

A multi-language library (Go & TypeScript) that encodes arbitrary binary data into natural-looking English sentences. Hide your secret messages in plain sight!

![Sentence Cipher Modes](assets/image-9b9cbf87-6348-49fa-a04e-9c8acf7380b8.png)

## Overview

Sentence Cipher transforms any byte sequence into grammatically correct English sentences using a workplace/office theme vocabulary. Unlike simple substitution ciphers, it constructs meaningful sentence structures that mimic professional communication.

### Example

```
Mode: String
Input:  "Hello"
Output: "ruth trains isabella prints. carl cleans daily."
```
```
Mode: Natural
Input:  "Hello"
Output:
"""
Subject: Code Review

Hi Team,

Ian imports sandra os. Specifically, Leo sorts daily. 

Best regards,
Samantha
"""
```
```
Mode: Binary
Input:  "Hello"
Output: "ruth trains isabella prints. carl cleans daily."
```

## How It Works

The library offers three primary modes of operation to suit different needs:

### 📝 String Mode
Standard encoding that transforms text into office-themed sentences. It handles UTF-8 strings directly, making it perfect for short text messages or chat applications.

### 📧 Natural Mode
Advanced encoding that structures data into a natural-looking email format. It wraps the encoded sentences with realistic:
- **Subjects**: "Code Review", "Meeting Notes", "Project Update"
- **Openers**: "Hi Team,", "Good morning,"
- **Connectors**: "Specifically,", "Additionally,"
- **Closers**: "Best regards,", "Cheers,"

This mode generates output that looks indistinguishable from a real workplace email, providing better cover for your data.

### 💾 Binary Mode
Raw byte encoding for any file type (images, documents, executables). It maintains data integrity by treating the input as a raw byte stream. The output looks the same as String Mode but ensures that binary data is perfectly preserved during the round-trip.

---

### Encoding Scheme

The cipher uses three word lists, each containing exactly 256 words (mapping to byte values 0-255):

| Word List | Purpose | Examples |
|-----------|---------|----------|
| **Names** | Subjects & Indirect Objects | adam, alex, alice... |
| **Verbs** | Actions | helps, assists, supports... |
| **Objects** | Direct Objects | reports, documents, files... |

### Sentence Patterns

Data is encoded using three different sentence patterns depending on the remaining bytes:

1.  **Full Sentence (3 bytes):** `[Subject] [Verb] [IndirectObject] [Object].`
2.  **Short Sentence (2 bytes):** `[Subject] [Verb] daily.`
3.  **Minimal Sentence (1 byte):** `[Subject] works.`

## Installation

### Go

```bash
go get github.com/kittizz/sentence-cipher
```

### Node.js / TypeScript

```bash
npm install sentence-cipher
# or
pnpm add sentence-cipher
```

## Usage

### Go Library

```go
package main

import (
    "fmt"
    sentencecipher "github.com/kittizz/sentence-cipher"
)

func main() {
    // 1. Basic encoding/decoding
    message := "Hello World"
    encoded := sentencecipher.EncodeString(message)
    fmt.Println(encoded) 
    // Output: "ruth trains isabella prints. carl cleans daily..."

    // 2. Natural Mode (Email style)
    naturalEncoded := sentencecipher.EncodeNatural([]byte(message))
    fmt.Println(naturalEncoded)
    // Output: "Subject: Project Update\nHi Team,\nRuth trains..."

    // 3. With Encryption Key
    cipher := sentencecipher.NewCipher("my-secret-key")
    secureEncoded := cipher.EncodeString(message)
    
    // Decode
    decoded, _ := cipher.DecodeString(secureEncoded)
    fmt.Println(decoded)
}
```

### Node.js / TypeScript Library

The JavaScript library provides a Type-Safe API with full support for all modes.

```typescript
import { createCipher, createDefaultCipher } from 'sentence-cipher';

// Initialize
const cipher = createDefaultCipher();

// --- String Mode ---
const encoded = cipher.encodeString("Hello World");
const decoded = cipher.decodeString(encoded);
console.log(encoded); 
// "ruth trains isabella prints. carl cleans daily."


// --- Natural Mode ---
// Natural mode works with Uint8Array to support any data
const input = new TextEncoder().encode("Hello World");
const naturalEncoded = cipher.encodeNatural(input);

console.log(naturalEncoded);
/*
Subject: Team Sync
Hi Everyone,
Ruth trains isabella prints. Carl cleans daily.
Best,
*/

const naturalDecoded = cipher.decodeNatural(naturalEncoded);
console.log(new TextDecoder().decode(naturalDecoded)); // "Hello World"


// --- Binary Mode ---
// Perfect for files or raw data
const data = new Uint8Array([72, 101, 108, 108, 111]); // "Hello" bytes
const binaryEncoded = cipher.encode(data);
const binaryDecoded = cipher.decode(binaryEncoded);


// --- Encryption Key ---
// Provide a key to shuffle the word lists deterministically
const secureCipher = createCipher("my-super-secret-key");
const secureMsg = secureCipher.encodeString("Secret Data");

// Must use the same key to decode
const decryptedMsg = secureCipher.decodeString(secureMsg);
```

### Command Line Tool (CLI)

You can also use the Go CLI to encode/decode files or text from the terminal.

```bash
# Install
go install github.com/kittizz/sentence-cipher/cmd/sentencecipher@latest

# Usage
sentencecipher "Hello World"
sentencecipher -d "ruth trains isabella prints..."
sentencecipher -n "Generate natural email"
sentencecipher -k "my-key" "Encrypted message"
```

## Technical Details

### Key Derivation & Security
When a key is provided, the library uses **SHA-256** to hash the key. The hash is used to seed a **Fisher-Yates shuffle** algorithm, which randomizes the order of the Names, Verbs, and Objects word lists. This ensures that without the correct key, the sentence mapping is completely different, effectively encrypting the message.

### Expansion Ratio
The encoding transforms binary data into English text, which naturally increases the size.
- **Expansion:** Approximately 15-20x original size.
- **Efficiency:** 
    - 3 bytes → ~4 words
    - 2 bytes → ~3 words
    - 1 byte  → ~2 words

## License

MIT License
