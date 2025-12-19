package sentencecipher

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
)

const Version = "2.0.1"

// Cipher holds the shuffled word lists based on a key
type Cipher struct {
	names   []string
	verbs   []string
	objects []string
	key     string // Store key for regenerating themed ciphers
}

// NewCipher creates a new Cipher with word lists shuffled based on the provided key
func NewCipher(key string) (*Cipher, error) {
	if key == "" {
		return nil, errors.New("key is required")
	}
	hash := sha256.Sum256([]byte(key))
	seed := int64(binary.BigEndian.Uint64(hash[:8]))

	return &Cipher{
		names:   shuffleWithSeed(defaultNames, seed),
		verbs:   shuffleWithSeed(defaultVerbs, seed+1),
		objects: shuffleWithSeed(defaultObjects, seed+2),
		key:     key,
	}, nil
}

// NewDefaultCipher creates a Cipher with default (unshuffled) word lists
func NewDefaultCipher() *Cipher {
	return &Cipher{
		names:   copySlice(defaultNames),
		verbs:   copySlice(defaultVerbs),
		objects: copySlice(defaultObjects),
		key:     "",
	}
}

// NewThemedCipher creates a Cipher based on a specific theme
func NewThemedCipher(key string, theme string) *Cipher {
	var verbs, objects []string

	switch theme {
	case "tech":
		verbs = techVerbs
		objects = techObjects
	case "business":
		verbs = defaultVerbs // default is business
		objects = defaultObjects
	default:
		verbs = defaultVerbs
		objects = defaultObjects
	}

	if key != "" {
		// If key is provided, shuffle the themed lists
		hash := sha256.Sum256([]byte(key))
		seed := int64(binary.BigEndian.Uint64(hash[:8]))

		return &Cipher{
			names:   shuffleWithSeed(defaultNames, seed),
			verbs:   shuffleWithSeed(verbs, seed+1),
			objects: shuffleWithSeed(objects, seed+2),
			key:     key,
		}
	}

	// Default (unshuffled) but themed
	return &Cipher{
		names:   copySlice(defaultNames),
		verbs:   copySlice(verbs),
		objects: copySlice(objects),
		key:     "",
	}
}

func copySlice(src []string) []string {
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func shuffleWithSeed(src []string, seed int64) []string {
	dst := copySlice(src)
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(dst), func(i, j int) {
		dst[i], dst[j] = dst[j], dst[i]
	})
	return dst
}

// Sentence pattern: "Subject verb IndirectObject object."
// Each sentence encodes 3 bytes:
// - Byte 1: Subject (name index 0-255)
// - Byte 2: Verb (verb index 0-255)
// - Byte 3: Object (object index 0-255)
// IndirectObject is derived from (byte1 + byte2) % 256 for natural flow

// Encode converts bytes to English sentences (Cipher method)
func (c *Cipher) Encode(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	var sentences []string
	i := 0

	for i < len(data) {
		remaining := len(data) - i

		if remaining >= 3 {
			// Full pattern: S + V + IO + O (encodes 3 bytes)
			b1 := data[i]   // subject
			b2 := data[i+1] // verb
			b3 := data[i+2] // object

			// Rotation logic with position offset
			offset := i
			idx1 := (int(b1) + offset) % 256
			idx2 := (int(b2) + offset + 1) % 256
			idx3 := (int(b3) + offset + 2) % 256

			subject := c.names[idx1]
			verb := c.verbs[idx2]
			// IO derived for natural flow using rotated indices
			ioIdx := (idx1 + idx2) % 256
			indirectObj := c.names[ioIdx]
			obj := c.objects[idx3]

			sentences = append(sentences, fmt.Sprintf("%s %s %s %s.", subject, verb, indirectObj, obj))
			i += 3
		} else if remaining == 2 {
			// Pattern: S + V + O (encodes 2 bytes)
			b1 := data[i]
			b2 := data[i+1]

			offset := i
			idx1 := (int(b1) + offset) % 256
			idx2 := (int(b2) + offset + 1) % 256

			subject := c.names[idx1]
			verb := c.verbs[idx2]

			sentences = append(sentences, fmt.Sprintf("%s %s daily.", subject, verb))
			i += 2
		} else {
			// Pattern: S + works (encodes 1 byte)
			b := data[i]
			offset := i
			idx := (int(b) + offset) % 256

			subject := c.names[idx]

			sentences = append(sentences, fmt.Sprintf("%s works.", subject))
			i++
		}
	}

	return strings.Join(sentences, " ")
}

// Encode converts bytes to English sentences (package-level function for backward compatibility)
func Encode(data []byte) string {
	return NewDefaultCipher().Encode(data)
}

// Decode converts English sentences back to bytes (Cipher method)
func (c *Cipher) Decode(encoded string) ([]byte, error) {
	if encoded == "" {
		return []byte{}, nil
	}

	sentences := splitSentences(encoded)
	var result []byte
	byteCount := 0 // Track byte position for rotation offset

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		words := strings.Fields(sentence)
		if len(words) == 0 {
			continue
		}

		// Remove trailing punctuation from last word
		lastIdx := len(words) - 1
		words[lastIdx] = strings.TrimSuffix(words[lastIdx], ".")

		// Get rotated bytes from decodeSentence (which returns indices basically)
		chunkBytes, err := c.decodeSentence(words)
		if err != nil {
			return nil, err
		}

		// Un-rotate bytes using position offset
		for j, bRotated := range chunkBytes {
			offset := byteCount + j
			// Reverse rotation: val = (rotated - offset) % 256
			val := (int(bRotated) - offset) % 256
			if val < 0 {
				val += 256
			}
			result = append(result, byte(val))
		}
		byteCount += len(chunkBytes)
	}

	return result, nil
}

// Decode converts English sentences back to bytes (package-level function for backward compatibility)
func Decode(encoded string) ([]byte, error) {
	return NewDefaultCipher().Decode(encoded)
}

func splitSentences(text string) []string {
	var sentences []string
	var current strings.Builder

	for _, r := range text {
		current.WriteRune(r)
		if r == '.' {
			sentences = append(sentences, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		sentences = append(sentences, current.String())
	}

	return sentences
}

func (c *Cipher) decodeSentence(words []string) ([]byte, error) {
	// Clean up words to ensure lowercase for matching
	cleanWords := make([]string, len(words))
	for i, w := range words {
		cleanWords[i] = strings.ToLower(w)
	}
	words = cleanWords

	if len(words) < 2 {
		return nil, errors.New("invalid sentence: too few words")
	}

	// Pattern: "subject works" (1 byte)
	if len(words) == 2 && words[1] == "works" {
		sIdx := findIndex(c.names, words[0])
		if sIdx == -1 {
			return nil, errors.New("unknown name: " + words[0])
		}
		return []byte{byte(sIdx)}, nil
	}

	// Pattern: "subject verb daily" (2 bytes)
	if len(words) == 3 && words[2] == "daily" {
		sIdx := findIndex(c.names, words[0])
		if sIdx == -1 {
			return nil, errors.New("unknown name: " + words[0])
		}
		vIdx := findIndex(c.verbs, words[1])
		if vIdx == -1 {
			return nil, errors.New("unknown verb: " + words[1])
		}
		return []byte{byte(sIdx), byte(vIdx)}, nil
	}

	// Pattern: "subject verb indirectObject object" (3 bytes)
	if len(words) == 4 {
		sIdx := findIndex(c.names, words[0])
		if sIdx == -1 {
			return nil, errors.New("unknown name: " + words[0])
		}
		vIdx := findIndex(c.verbs, words[1])
		if vIdx == -1 {
			return nil, errors.New("unknown verb: " + words[1])
		}
		// IO is derived, we don't decode it
		oIdx := findIndex(c.objects, words[3])
		if oIdx == -1 {
			return nil, errors.New("unknown object: " + words[3])
		}
		return []byte{byte(sIdx), byte(vIdx), byte(oIdx)}, nil
	}

	return nil, errors.New("unrecognized sentence pattern: " + strings.Join(words, " "))
}

// decodeSentence for backward compatibility (uses default word lists)
func decodeSentence(words []string) ([]byte, error) {
	return NewDefaultCipher().decodeSentence(words)
}

func findIndex(list []string, word string) int {
	for i, w := range list {
		if strings.EqualFold(w, word) {
			return i
		}
	}
	return -1
}

// EncodeString encodes a string (Cipher method)
func (c *Cipher) EncodeString(s string) string {
	return c.Encode([]byte(s))
}

// DecodeString decodes to a string (Cipher method)
func (c *Cipher) DecodeString(encoded string) (string, error) {
	data, err := c.Decode(encoded)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// EncodeString is a convenience function for encoding strings (backward compatibility)
func EncodeString(s string) string {
	return Encode([]byte(s))
}

// DecodeString is a convenience function for decoding to strings (backward compatibility)
func DecodeString(encoded string) (string, error) {
	data, err := Decode(encoded)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// EncodeNatural creates more varied, natural-looking sentences structured as an email (Cipher method)
func (c *Cipher) EncodeNatural(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// Calculate a simple seed from data to make deterministic choices
	// Mix in key hash so that different keys produce different seeds/themes
	seed := 0
	if c.key != "" {
		hash := sha256.Sum256([]byte(c.key))
		// Use modulo on uint64 before casting to int to avoid negative values from overflow
		seed = int(binary.BigEndian.Uint64(hash[:8]) % 10000)
	}

	for _, b := range data {
		seed = (seed + int(b)) % 10000
	}

	// Determine Theme
	// 50% chance for Tech, 50% for Business
	theme := "business"
	if seed%2 == 0 {
		theme = "tech"
	}

	// Select Subject based on Theme
	var subj string
	if theme == "tech" {
		subj = techSubjects[seed%len(techSubjects)]
	} else {
		subj = businessSubjects[seed%len(businessSubjects)]
	}

	// Create Themed Cipher to encode the body
	// This ensures the words match the theme
	themedCipher := NewThemedCipher(c.key, theme)

	// Generate basic sentences first using the themed cipher
	var sentences []string
	i := 0
	for i < len(data) {
		remaining := len(data) - i
		var s string

		if remaining >= 3 {
			b1 := data[i]
			b2 := data[i+1]
			b3 := data[i+2]

			// Rotation logic to prevent repeating words for same bytes
			// We use position (i) to offset the index
			offset := i
			idx1 := (int(b1) + offset) % 256
			idx2 := (int(b2) + offset + 1) % 256
			idx3 := (int(b3) + offset + 2) % 256

			subject := themedCipher.names[idx1]
			verb := themedCipher.verbs[idx2]
			ioIdx := (idx1 + idx2) % 256
			indirectObj := themedCipher.names[ioIdx]
			obj := themedCipher.objects[idx3]
			s = fmt.Sprintf("%s %s %s %s", subject, verb, indirectObj, obj)
			i += 3
		} else if remaining == 2 {
			b1 := data[i]
			b2 := data[i+1]

			offset := i
			idx1 := (int(b1) + offset) % 256
			idx2 := (int(b2) + offset + 1) % 256

			subject := themedCipher.names[idx1]
			verb := themedCipher.verbs[idx2]
			s = fmt.Sprintf("%s %s daily", subject, verb)
			i += 2
		} else {
			b := data[i]

			offset := i
			idx := (int(b) + offset) % 256

			subject := themedCipher.names[idx]
			s = fmt.Sprintf("%s works", subject)
			i++
		}
		sentences = append(sentences, s)
	}

	// Construct Email
	var sb strings.Builder

	// Subject Line
	sb.WriteString("Subject: " + subj + "\n\n")

	// Opener
	opener := emailOpeners[seed%len(emailOpeners)]
	sb.WriteString(opener + "\n\n")

	// Body
	for idx, s := range sentences {
		// Capitalize first letter of the sentence logic happens after connector logic
		prefix := ""

		// Add connector mostly for subsequent sentences, not the first one
		if idx > 0 {
			// Use deterministic pseudo-randomness based on index + seed
			r := (seed + idx) % 10
			if r < 6 { // 60% chance to add a connector
				conn := sentenceConnectors[(seed+idx)%len(sentenceConnectors)]
				prefix = conn + " "
			}
		}

		// Capitalize sentence
		s = capitalize(s)
		sb.WriteString(prefix + s + ".")

		// Paragraph spacing
		if (idx+1)%3 == 0 && idx < len(sentences)-1 {
			sb.WriteString("\n\n")
		} else {
			sb.WriteString(" ")
		}
	}

	// Closer
	sb.WriteString("\n\n")
	closer := emailClosers[seed%len(emailClosers)]
	sb.WriteString(closer + "\n")

	// Random Sender (use one of the names based on seed)
	senderIdx := (seed * 7) % len(themedCipher.names)
	sender := capitalize(themedCipher.names[senderIdx])
	sb.WriteString(sender)

	return sb.String()
}

// EncodeNatural creates more varied, natural-looking sentences (backward compatibility)
func EncodeNatural(data []byte) string {
	return NewDefaultCipher().EncodeNatural(data)
}

// DecodeNatural decodes text created by EncodeNatural (Cipher method)
func (c *Cipher) DecodeNatural(encoded string) ([]byte, error) {
	if encoded == "" {
		return []byte{}, nil
	}

	lines := strings.Split(encoded, "\n")
	var bodyLines []string

	// Detect Theme from Subject
	theme := "business" // default

	// Scan for Subject first
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "subject:") {
			subj := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(line), "subject:"))

			// Check if subject belongs to Tech
			for _, s := range techSubjects {
				if strings.EqualFold(s, subj) {
					theme = "tech"
					break
				}
			}
			break
		}
	}

	// Use Themed Cipher for decoding
	themedCipher := NewThemedCipher(c.key, theme)

	// Filter structure
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip Subject line
		if strings.HasPrefix(strings.ToLower(line), "subject:") {
			continue
		}

		// Skip Openers
		isOpener := false
		for _, op := range emailOpeners {
			if strings.TrimSpace(op) == line {
				isOpener = true
				break
			}
		}
		if isOpener {
			continue
		}

		// Skip Closers
		isCloser := false
		for _, cl := range emailClosers {
			if strings.TrimSpace(cl) == line {
				isCloser = true
				break
			}
		}
		if isCloser {
			continue
		}

		// Check if it's just a name (Sender signature)
		// This is tricky because a name could be a valid 1-byte sentence "Name works." but here signature is just "Name"
		// Signature usually doesn't end with "."
		// Our sentences always end with "."
		if !strings.HasSuffix(line, ".") {
			// Likely a signature or subject garbage
			continue
		}

		bodyLines = append(bodyLines, line)
	}

	fullBody := strings.Join(bodyLines, " ")

	// Split by sentences
	// Note: Connectors usually end with comma, sentences end with period.
	rawSentences := splitSentences(fullBody)

	var result []byte
	byteCount := 0 // Track byte position for rotation offset

	for _, sentence := range rawSentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		// Remove connectors
		for _, conn := range sentenceConnectors {
			// Connector usually has space after it
			// Check case-insensitive matching for connectors?
			// In Encode, we use them as is (Title case).
			if strings.HasPrefix(sentence, conn) {
				sentence = strings.TrimPrefix(sentence, conn)
				sentence = strings.TrimSpace(sentence)
				break
			}
		}

		// Clean up punctuation just in case
		sentence = strings.TrimSuffix(sentence, ".")

		words := strings.Fields(sentence)
		if len(words) == 0 {
			continue
		}

		// Lowercase everything for decoding
		for i := range words {
			words[i] = strings.ToLower(words[i])
		}

		// Decode using the new logic that respects offset
		// We use themedCipher instead of 'c' to lookup verbs/objects

		// Inline logic for rotation awareness:
		if len(words) < 2 {
			return nil, fmt.Errorf("invalid sentence: too few words: %s", sentence)
		}

		// Pattern: "subject works" (1 byte)
		if len(words) == 2 && words[1] == "works" {
			sIdx := findIndex(themedCipher.names, words[0])
			if sIdx == -1 {
				return nil, fmt.Errorf("unknown name: %s", words[0])
			}

			offset := byteCount
			b := (sIdx - offset + 256) % 256
			result = append(result, byte(b))
			byteCount++
			continue
		}

		// Pattern: "subject verb daily" (2 bytes)
		if len(words) == 3 && words[2] == "daily" {
			sIdx := findIndex(themedCipher.names, words[0])
			if sIdx == -1 {
				return nil, fmt.Errorf("unknown name: %s", words[0])
			}
			vIdx := findIndex(themedCipher.verbs, words[1])
			if vIdx == -1 {
				return nil, fmt.Errorf("unknown verb: %s", words[1])
			}

			offset := byteCount
			b1 := (sIdx - offset + 256) % 256
			b2 := (vIdx - offset - 1 + 256) % 256

			result = append(result, byte(b1), byte(b2))
			byteCount += 2
			continue
		}

		// Pattern: "subject verb indirectObject object" (3 bytes)
		if len(words) == 4 {
			sIdx := findIndex(themedCipher.names, words[0])
			if sIdx == -1 {
				return nil, fmt.Errorf("unknown name: %s", words[0])
			}
			vIdx := findIndex(themedCipher.verbs, words[1])
			if vIdx == -1 {
				return nil, fmt.Errorf("unknown verb: %s", words[1])
			}
			oIdx := findIndex(themedCipher.objects, words[3])
			if oIdx == -1 {
				return nil, fmt.Errorf("unknown object: %s", words[3])
			}

			offset := byteCount
			b1 := (sIdx - offset + 256) % 256
			b2 := (vIdx - offset - 1 + 256) % 256
			b3 := (oIdx - offset - 2 + 256) % 256

			result = append(result, byte(b1), byte(b2), byte(b3))
			byteCount += 3
			continue
		}

		return nil, fmt.Errorf("unrecognized sentence pattern: %s", sentence)
	}

	return result, nil
}

// DecodeNatural decodes text created by EncodeNatural (backward compatibility)
func DecodeNatural(encoded string) ([]byte, error) {
	return NewDefaultCipher().DecodeNatural(encoded)
}

// Helper to capitalize first letter
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// Debug helpers (Cipher methods)
func (c *Cipher) DebugByte(b byte) string {
	return fmt.Sprintf("byte=%d (0x%02X) -> name=%s", b, b, c.names[b])
}

func (c *Cipher) DebugEncode(data []byte) {
	fmt.Println("=== Debug Encode ===")
	for i, b := range data {
		fmt.Printf("  [%d] %s\n", i, c.DebugByte(b))
	}
	fmt.Println("Encoded:", c.Encode(data))
}

// Debug helpers (backward compatibility)
func DebugByte(b byte) string {
	return NewDefaultCipher().DebugByte(b)
}

func DebugEncode(data []byte) {
	NewDefaultCipher().DebugEncode(data)
}

// Hex encoding for compatibility
func EncodeHex(data []byte) string {
	var parts []string
	for _, b := range data {
		parts = append(parts, strconv.FormatInt(int64(b), 16))
	}
	return strings.Join(parts, "")
}
