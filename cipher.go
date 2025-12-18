package sentencecipher

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// Cipher holds the shuffled word lists based on a key
type Cipher struct {
	names   []string
	verbs   []string
	objects []string
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
	}, nil
}

// NewDefaultCipher creates a Cipher with default (unshuffled) word lists
func NewDefaultCipher() *Cipher {
	return &Cipher{
		names:   copySlice(defaultNames),
		verbs:   copySlice(defaultVerbs),
		objects: copySlice(defaultObjects),
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

			subject := c.names[b1]
			verb := c.verbs[b2]
			// IO derived for natural flow
			ioIdx := (int(b1) + int(b2)) % 256
			indirectObj := c.names[ioIdx]
			obj := c.objects[b3]

			sentences = append(sentences, fmt.Sprintf("%s %s %s %s.", subject, verb, indirectObj, obj))
			i += 3
		} else if remaining == 2 {
			// Pattern: S + V + O (encodes 2 bytes)
			b1 := data[i]
			b2 := data[i+1]

			subject := c.names[b1]
			verb := c.verbs[b2]

			sentences = append(sentences, fmt.Sprintf("%s %s daily.", subject, verb))
			i += 2
		} else {
			// Pattern: S + works (encodes 1 byte)
			b := data[i]
			subject := c.names[b]

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

		bytes, err := c.decodeSentence(words)
		if err != nil {
			return nil, err
		}
		result = append(result, bytes...)
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

// EncodeNatural creates more varied, natural-looking sentences (Cipher method)
func (c *Cipher) EncodeNatural(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	templates := []string{
		"%s %s %s %s.",
		"today %s %s %s %s.",
		"later %s %s %s %s.",
		"then %s %s %s %s.",
		"now %s %s %s %s.",
	}

	var sentences []string
	i := 0
	templateIdx := 0

	for i < len(data) {
		remaining := len(data) - i

		if remaining >= 3 {
			b1 := data[i]
			b2 := data[i+1]
			b3 := data[i+2]

			subject := c.names[b1]
			verb := c.verbs[b2]
			ioIdx := (int(b1) + int(b2)) % 256
			indirectObj := c.names[ioIdx]
			obj := c.objects[b3]

			template := templates[templateIdx%len(templates)]
			sentences = append(sentences, fmt.Sprintf(template, subject, verb, indirectObj, obj))
			templateIdx++
			i += 3
		} else if remaining == 2 {
			b1 := data[i]
			b2 := data[i+1]
			subject := c.names[b1]
			verb := c.verbs[b2]
			sentences = append(sentences, fmt.Sprintf("%s %s daily.", subject, verb))
			i += 2
		} else {
			b := data[i]
			subject := c.names[b]
			sentences = append(sentences, fmt.Sprintf("%s works.", subject))
			i++
		}
	}

	return strings.Join(sentences, " ")
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

	sentences := splitSentences(encoded)
	var result []byte

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		// Remove time prefixes
		sentence = strings.TrimPrefix(sentence, "today ")
		sentence = strings.TrimPrefix(sentence, "later ")
		sentence = strings.TrimPrefix(sentence, "then ")
		sentence = strings.TrimPrefix(sentence, "now ")
		sentence = strings.TrimSuffix(sentence, ".")

		words := strings.Fields(sentence)
		if len(words) == 0 {
			continue
		}

		bytes, err := c.decodeSentence(words)
		if err != nil {
			return nil, fmt.Errorf("error decoding '%s': %w", sentence, err)
		}
		result = append(result, bytes...)
	}

	return result, nil
}

// DecodeNatural decodes text created by EncodeNatural (backward compatibility)
func DecodeNatural(encoded string) ([]byte, error) {
	return NewDefaultCipher().DecodeNatural(encoded)
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
