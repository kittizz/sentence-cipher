package sentencecipher

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"empty", []byte{}},
		{"single byte", []byte{0x00}},
		{"single byte max", []byte{0xFF}},
		{"two bytes", []byte{0x48, 0x69}}, // "Hi"
		{"hello", []byte("Hello")},
		{"hello world", []byte("Hello, World!")},
		{"numbers", []byte("12345")},
		{"special chars", []byte("!@#$%")},
		{"thai text", []byte("สวัสดี")},
		{"mixed", []byte("Hello สวัสดี 123")},
		{"all byte values", func() []byte {
			b := make([]byte, 256)
			for i := 0; i < 256; i++ {
				b[i] = byte(i)
			}
			return b
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := Encode(tt.input)
			if err != nil {
				t.Fatalf("Encode error: %v", err)
			}
			decoded, err := Decode(encoded)
			if err != nil {
				t.Fatalf("Decode error: %v", err)
			}
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("Encode/Decode mismatch\noriginal: %v\ndecoded:  %v", tt.input, decoded)
			}
		})
	}
}

func TestEncodeDecodeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"simple", "Hello"},
		{"with spaces", "Hello World"},
		{"with punctuation", "Hello, World!"},
		{"thai", "สวัสดีครับ"},
		{"emoji", "Hello 👋"},
		{"secret message", "Meet at 9pm"},
		{"long text", "The quick brown fox jumps over the lazy dog"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := EncodeString(tt.input)
			if err != nil {
				t.Fatalf("EncodeString error: %v", err)
			}
			decoded, err := DecodeString(encoded)
			if err != nil {
				t.Fatalf("DecodeString error: %v", err)
			}
			if decoded != tt.input {
				t.Errorf("EncodeString/DecodeString mismatch\noriginal: %q\ndecoded:  %q", tt.input, decoded)
			}
		})
	}
}

func TestEncodeNaturalDecodeNatural(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"empty", []byte{}},
		{"hello", []byte("Hello")},
		{"secret", []byte("Secret message")},
		{"thai", []byte("สวัสดี")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := EncodeNatural(tt.input)
			if err != nil {
				t.Fatalf("EncodeNatural error: %v", err)
			}
			decoded, err := DecodeNatural(encoded)
			if err != nil {
				t.Fatalf("DecodeNatural error: %v", err)
			}
			t.Logf("Encoded: %q", encoded)
			t.Logf("Decoded: %q", decoded)
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("EncodeNatural/DecodeNatural mismatch\noriginal: %v\ndecoded:  %v", tt.input, decoded)
			}
		})
	}
}

func TestEncodedLooksLikeEnglish(t *testing.T) {
	input := []byte("Secret")
	encoded, err := Encode(input)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	// Check that encoded text contains expected sentence structure
	if len(encoded) == 0 {
		t.Error("Encoded text should not be empty")
	}

	// Should end with a period
	if encoded[len(encoded)-1] != '.' {
		t.Error("Encoded text should end with a period")
	}

	// Should contain spaces (word separation)
	if !bytes.Contains([]byte(encoded), []byte(" ")) {
		t.Error("Encoded text should contain spaces")
	}

	t.Logf("Input: %q", input)
	t.Logf("Encoded: %q", encoded)
}

func TestNaturalEncodingVariety(t *testing.T) {
	input := []byte("Hello World!")
	encoded, err := EncodeNatural(input)
	if err != nil {
		t.Fatalf("EncodeNatural error: %v", err)
	}

	// Natural encoding should use varied templates
	t.Logf("Natural encoded: %q", encoded)

	// Should contain some variety markers (Email structure)
	hasVariety := bytes.Contains([]byte(encoded), []byte("Subject:")) ||
		bytes.Contains([]byte(encoded), []byte("Hi")) ||
		bytes.Contains([]byte(encoded), []byte("Dear")) ||
		bytes.Contains([]byte(encoded), []byte("Regards,"))

	if !hasVariety {
		t.Log("Warning: Natural encoding might not have enough variety")
	}
}

func TestDecodeInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"gibberish", "asdfghjkl"},
		{"unknown subject", "Xyz loves something."},
		{"unknown verb", "Tom xyzs something."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decode(tt.input)
			if err == nil {
				t.Error("Expected error for invalid input, got nil")
			}
		})
	}
}

func TestAllByteValues(t *testing.T) {
	// Test that all possible byte values can be encoded and decoded
	for i := 0; i < 256; i++ {
		input := []byte{byte(i)}
		encoded, err := Encode(input)
		if err != nil {
			t.Errorf("Failed to encode byte %d: %v", i, err)
			continue
		}
		decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf("Failed to decode byte %d: %v", i, err)
			continue
		}
		if len(decoded) != 1 || decoded[0] != byte(i) {
			t.Errorf("Byte %d: expected %v, got %v", i, input, decoded)
		}
	}
}

func TestAllTwoByteCombinatons(t *testing.T) {
	// Test a sample of two-byte combinations
	samples := [][2]byte{
		{0x00, 0x00},
		{0xFF, 0xFF},
		{0x00, 0xFF},
		{0xFF, 0x00},
		{0x48, 0x69}, // "Hi"
		{0xAB, 0xCD},
		{0x12, 0x34},
	}

	for _, sample := range samples {
		input := []byte{sample[0], sample[1]}
		encoded, err := Encode(input)
		if err != nil {
			t.Errorf("Failed to encode %v: %v", input, err)
			continue
		}
		decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf("Failed to decode %v: %v", input, err)
			continue
		}
		if !bytes.Equal(decoded, input) {
			t.Errorf("Two bytes %v: expected %v, got %v", sample, input, decoded)
		}
	}
}

func TestBinaryData(t *testing.T) {
	// PNG header signature
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	tests := []struct {
		name  string
		input []byte
	}{
		{"png header", pngHeader},
		{"null bytes", []byte{0x00, 0x00, 0x00, 0x00}},
		{"high bytes", []byte{0xFF, 0xFE, 0xFD, 0xFC}},
		{"mixed binary", []byte{0x00, 0xFF, 0x7F, 0x80, 0x01, 0xFE}},
		{"random binary", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := Encode(tt.input)
			if err != nil {
				t.Fatalf("Encode error: %v", err)
			}
			decoded, err := Decode(encoded)
			if err != nil {
				t.Fatalf("Decode error: %v", err)
			}
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("Binary encode/decode mismatch\noriginal: %v\ndecoded:  %v", tt.input, decoded)
			}
		})
	}
}

func TestCipherWithKey(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		input []byte
	}{
		{"simple key", "secret", []byte("Hello")},
		{"complex key", "my-super-secret-key-123!", []byte("Secret message")},
		{"binary with key", "test-key", []byte{0x89, 0x50, 0x4E, 0x47}},
		{"all bytes with key", "key123", func() []byte {
			b := make([]byte, 256)
			for i := 0; i < 256; i++ {
				b[i] = byte(i)
			}
			return b
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cipher, err := NewCipher(tt.key)
			if err != nil {
				t.Fatalf("NewCipher error: %v", err)
			}

			encoded, err := cipher.Encode(tt.input)
			if err != nil {
				t.Fatalf("Encode error: %v", err)
			}
			decoded, err := cipher.Decode(encoded)
			if err != nil {
				t.Fatalf("Decode error: %v", err)
			}
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("Cipher encode/decode mismatch\noriginal: %v\ndecoded:  %v", tt.input, decoded)
			}
		})
	}
}

func TestCipherNaturalWithKey(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		input []byte
	}{
		{"natural with key", "my-key", []byte("Hello World")},
		{"binary natural", "binary-key", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cipher, err := NewCipher(tt.key)
			if err != nil {
				t.Fatalf("NewCipher error: %v", err)
			}

			encoded, err := cipher.EncodeNatural(tt.input)
			if err != nil {
				t.Fatalf("EncodeNatural error: %v", err)
			}
			decoded, err := cipher.DecodeNatural(encoded)
			if err != nil {
				t.Fatalf("DecodeNatural error: %v", err)
			}
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("Cipher natural encode/decode mismatch\noriginal: %v\ndecoded:  %v", tt.input, decoded)
			}
		})
	}
}

func TestDifferentKeysProduceDifferentOutput(t *testing.T) {
	input := []byte("Secret message")

	cipher1, _ := NewCipher("key1")
	cipher2, _ := NewCipher("key2")
	defaultCipher := NewDefaultCipher()

	encoded1, _ := cipher1.Encode(input)
	encoded2, _ := cipher2.Encode(input)
	encodedDefault, _ := defaultCipher.Encode(input)

	if encoded1 == encoded2 {
		t.Error("Different keys should produce different encoded output")
	}
	if encoded1 == encodedDefault {
		t.Error("Keyed cipher should produce different output than default")
	}
	if encoded2 == encodedDefault {
		t.Error("Keyed cipher should produce different output than default")
	}
}

func TestWrongKeyCannotDecode(t *testing.T) {
	input := []byte("Secret message")

	cipher1, _ := NewCipher("correct-key")
	cipher2, _ := NewCipher("wrong-key")

	encoded, err := cipher1.Encode(input)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	decoded, err := cipher2.Decode(encoded)

	// With wrong key, either error or different data
	if err == nil && bytes.Equal(decoded, input) {
		t.Error("Wrong key should not decode correctly")
	}
}

func TestAllThreeByteGroups(t *testing.T) {
	// Test all possible 3-byte combinations for first 16 values
	// (full 256^3 would be too slow)
	for b1 := 0; b1 < 16; b1++ {
		for b2 := 0; b2 < 16; b2++ {
			for b3 := 0; b3 < 16; b3++ {
				input := []byte{byte(b1 * 16), byte(b2 * 16), byte(b3 * 16)}
				encoded, err := Encode(input)
				if err != nil {
					t.Errorf("Failed to encode %v: %v", input, err)
					continue
				}
				decoded, err := Decode(encoded)
				if err != nil {
					t.Errorf("Failed to decode %v: %v", input, err)
					continue
				}
				if !bytes.Equal(decoded, input) {
					t.Errorf("Three bytes %v: got %v", input, decoded)
				}
			}
		}
	}
}

func TestLargeBinaryData(t *testing.T) {
	// Test with larger binary data (simulating a small file)
	sizes := []int{100, 500, 1000, 2000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			input := make([]byte, size)
			for i := 0; i < size; i++ {
				input[i] = byte(i % 256)
			}

			encoded, err := Encode(input)
			if err != nil {
				t.Fatalf("Encode error for size %d: %v", size, err)
			}
			decoded, err := Decode(encoded)
			if err != nil {
				t.Fatalf("Decode error for size %d: %v", size, err)
			}
			if !bytes.Equal(decoded, input) {
				t.Errorf("Large binary mismatch for size %d", size)
			}
		})
	}
}

func TestLargeBinaryWithKey(t *testing.T) {
	cipher, _ := NewCipher("test-key-for-large-data")

	input := make([]byte, 1000)
	for i := 0; i < 1000; i++ {
		input[i] = byte(i % 256)
	}

	encoded, err := cipher.Encode(input)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	decoded, err := cipher.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if !bytes.Equal(decoded, input) {
		t.Error("Large binary with key mismatch")
	}
}

func TestImageFile(t *testing.T) {
	// Read actual PNG file
	imgData, err := os.ReadFile("icon16.png")
	if err != nil {
		t.Skipf("Skipping image test: %v", err)
	}

	t.Run("normal_mode", func(t *testing.T) {
		encoded, err := Encode(imgData)
		if err != nil {
			t.Fatalf("Encode error: %v", err)
		}
		decoded, err := Decode(encoded)
		if err != nil {
			t.Fatalf("Decode error: %v", err)
		}
		if !bytes.Equal(decoded, imgData) {
			t.Errorf("Image encode/decode mismatch: original %d bytes, decoded %d bytes", len(imgData), len(decoded))
		}
	})

	t.Run("natural_mode", func(t *testing.T) {
		encoded, err := EncodeNatural(imgData)
		if err != nil {
			t.Fatalf("EncodeNatural error: %v", err)
		}
		decoded, err := DecodeNatural(encoded)
		if err != nil {
			t.Fatalf("DecodeNatural error: %v", err)
		}
		if !bytes.Equal(decoded, imgData) {
			t.Errorf("Image natural encode/decode mismatch: original %d bytes, decoded %d bytes", len(imgData), len(decoded))
		}
	})

	t.Run("with_key", func(t *testing.T) {
		cipher, _ := NewCipher("image-test-key")
		encoded, err := cipher.Encode(imgData)
		if err != nil {
			t.Fatalf("Encode error: %v", err)
		}
		decoded, err := cipher.Decode(encoded)
		if err != nil {
			t.Fatalf("Decode error: %v", err)
		}
		if !bytes.Equal(decoded, imgData) {
			t.Errorf("Image with key encode/decode mismatch: original %d bytes, decoded %d bytes", len(imgData), len(decoded))
		}
	})

	t.Run("natural_with_key", func(t *testing.T) {
		cipher, _ := NewCipher("image-natural-key")
		encoded, err := cipher.EncodeNatural(imgData)
		if err != nil {
			t.Fatalf("EncodeNatural error: %v", err)
		}
		decoded, err := cipher.DecodeNatural(encoded)
		if err != nil {
			t.Fatalf("DecodeNatural error: %v", err)
		}
		if !bytes.Equal(decoded, imgData) {
			t.Errorf("Image natural with key encode/decode mismatch: original %d bytes, decoded %d bytes", len(imgData), len(decoded))
		}
	})
}

func BenchmarkEncode(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Encode(input)
	}
}

func BenchmarkDecode(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	encoded, _ := Encode(input)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode(encoded)
	}
}

func BenchmarkEncodeNatural(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = EncodeNatural(input)
	}
}

func BenchmarkDecodeNatural(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	encoded, _ := EncodeNatural(input)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodeNatural(encoded)
	}
}

// Example tests
func ExampleEncodeString() {
	encoded, _ := EncodeString("Hi")
	_ = encoded // Encoded looks like English sentences
}

func ExampleDecodeString() {
	encoded, _ := EncodeString("Hi")
	decoded, _ := DecodeString(encoded)
	_ = decoded // "Hi"
}

func TestWordListIntegrity(t *testing.T) {
	lists := []struct {
		name string
		list []string
	}{
		{"defaultNames", defaultNames},
		{"defaultVerbs", defaultVerbs},
		{"defaultObjects", defaultObjects},
		{"techVerbs", techVerbs},
		{"techObjects", techObjects},
	}

	for _, l := range lists {
		t.Run(l.name, func(t *testing.T) {
			if len(l.list) != 256 {
				t.Errorf("List %s has length %d, want 256", l.name, len(l.list))
			}
			checkDuplicates(t, l.name, l.list)
		})
	}
}

func checkDuplicates(t *testing.T, name string, list []string) {
	seen := make(map[string]int)
	for i, w := range list {
		w = strings.ToLower(w)
		if idx, exists := seen[w]; exists {
			t.Errorf("Duplicate in %s: '%s' at index %d and %d", name, w, idx, i)
		}
		seen[w] = i
	}
}

func TestSubjectOverlap(t *testing.T) {
	// Check if any subject exists in both business and tech subjects
	// This could confuse theme detection
	business := businessSubjects
	tech := techSubjects

	for _, b := range business {
		for _, te := range tech {
			if strings.EqualFold(b, te) {
				t.Errorf("Subject overlap found: '%s' is in both Business and Tech lists", b)
			}
		}
	}
}

func TestNoSpacesInWords(t *testing.T) {
	lists := []struct {
		name string
		list []string
	}{
		{"defaultNames", defaultNames},
		{"defaultVerbs", defaultVerbs},
		{"defaultObjects", defaultObjects},
		{"techVerbs", techVerbs},
		{"techObjects", techObjects},
	}

	for _, l := range lists {
		t.Run(l.name, func(t *testing.T) {
			for i, w := range l.list {
				if strings.Contains(w, " ") {
					t.Errorf("Space found in %s[%d]: '%s'", l.name, i, w)
				}
			}
		})
	}
}

func TestNoPunctuationInWords(t *testing.T) {
	lists := []struct {
		name string
		list []string
	}{
		{"defaultNames", defaultNames},
		{"defaultVerbs", defaultVerbs},
		{"defaultObjects", defaultObjects},
		{"techVerbs", techVerbs},
		{"techObjects", techObjects},
	}

	for _, l := range lists {
		t.Run(l.name, func(t *testing.T) {
			for i, w := range l.list {
				if strings.ContainsAny(w, ".,") {
					t.Errorf("Punctuation found in %s[%d]: '%s'", l.name, i, w)
				}
			}
		})
	}
}
