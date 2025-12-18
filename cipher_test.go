package sentencecipher

import (
	"bytes"
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
			encoded := Encode(tt.input)
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
			encoded := EncodeString(tt.input)
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
			encoded := EncodeNatural(tt.input)
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
	encoded := Encode(input)

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
	encoded := EncodeNatural(input)

	// Natural encoding should use varied templates
	t.Logf("Natural encoded: %q", encoded)

	// Should contain some variety markers
	hasVariety := bytes.Contains([]byte(encoded), []byte("always")) ||
		bytes.Contains([]byte(encoded), []byte("today")) ||
		bytes.Contains([]byte(encoded), []byte("Today,")) ||
		bytes.Contains([]byte(encoded), []byte("really"))

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
		encoded := Encode(input)
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
		encoded := Encode(input)
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

func BenchmarkEncode(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(input)
	}
}

func BenchmarkDecode(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	encoded := Encode(input)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(encoded)
	}
}

func BenchmarkEncodeNatural(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeNatural(input)
	}
}

func BenchmarkDecodeNatural(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	encoded := EncodeNatural(input)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DecodeNatural(encoded)
	}
}

// Example tests
func ExampleEncodeString() {
	encoded := EncodeString("Hi")
	_ = encoded // Encoded looks like English sentences
}

func ExampleDecodeString() {
	encoded := EncodeString("Hi")
	decoded, _ := DecodeString(encoded)
	_ = decoded // "Hi"
}
