package sentencecipher

import (
	"bytes"
	"compress/flate"
	"os"
	"testing"
)

func compressFlate(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func runComparison(t *testing.T, name string, data []byte) {
	t.Run(name, func(t *testing.T) {
		if len(data) == 0 {
			t.Skip("Empty data")
		}

		// Brotli compression (current implementation)
		brotliCompressed, err := compress(data)
		if err != nil {
			t.Fatalf("Brotli compression failed: %v", err)
		}

		// Flate compression
		flateCompressed, err := compressFlate(data)
		if err != nil {
			t.Fatalf("Flate compression failed: %v", err)
		}

		t.Logf("--- %s Data Analysis ---", name)
		t.Logf("Original size: %d bytes", len(data))
		t.Logf("Brotli compressed: %d bytes (%.2f%%)", len(brotliCompressed), float64(len(brotliCompressed))/float64(len(data))*100)
		t.Logf("Flate compressed:  %d bytes (%.2f%%)", len(flateCompressed), float64(len(flateCompressed))/float64(len(data))*100)

		if len(brotliCompressed) < len(flateCompressed) {
			diff := len(flateCompressed) - len(brotliCompressed)
			t.Logf(">> Brotli is smaller by %d bytes", diff)
		} else {
			diff := len(brotliCompressed) - len(flateCompressed)
			t.Logf(">> Flate is smaller by %d bytes", diff)
		}

		cipher := NewDefaultCipher()

		// Compare Standard Encoding (Cipher sentences)
		rawStd := cipher.encodeRaw(data)
		brotliStd := cipher.encodeRaw(brotliCompressed)
		flateStd := cipher.encodeRaw(flateCompressed)

		t.Logf("\n[Standard Encoding]")
		t.Logf("Raw (No Compress): %d chars", len(rawStd))
		t.Logf("Brotli encoded:    %d chars (%.2f%% of Raw)", len(brotliStd), float64(len(brotliStd))/float64(len(rawStd))*100)
		t.Logf("Flate encoded:     %d chars (%.2f%% of Raw)", len(flateStd), float64(len(flateStd))/float64(len(rawStd))*100)
		diffStd := len(flateStd) - len(brotliStd)
		t.Logf("Brotli vs Flate diff: %d chars", diffStd)

		// Compare Natural Encoding (Email style)
		rawNat := cipher.encodeNaturalRaw(data)
		brotliNat := cipher.encodeNaturalRaw(brotliCompressed)
		flateNat := cipher.encodeNaturalRaw(flateCompressed)

		t.Logf("\n[Natural Encoding]")
		t.Logf("Raw (No Compress): %d chars", len(rawNat))
		t.Logf("Brotli encoded:    %d chars (%.2f%% of Raw)", len(brotliNat), float64(len(brotliNat))/float64(len(rawNat))*100)
		t.Logf("Flate encoded:     %d chars (%.2f%% of Raw)", len(flateNat), float64(len(flateNat))/float64(len(rawNat))*100)
		diffNat := len(flateNat) - len(brotliNat)
		t.Logf("Brotli vs Flate diff: %d chars", diffNat)
		t.Logf("---------------------------")
	})
}

func TestCompressionComparison(t *testing.T) {
	// 1. Text Input
	input := `Dashwood contempt on mr unlocked resolved provided of of. Stanhill wondered it it welcomed oh. Hundred no prudent he however smiling at an offence. If earnestly extremity he he propriety something admitting convinced ye. Pleasant in to although as if differed horrible. Mirth his quick its set front enjoy hoped had there. Who connection imprudence middletons too but increasing celebrated principles joy. Herself too improve gay winding ask expense are compact. New all paid few hard pure she.

Insipidity the sufficient discretion imprudence resolution sir him decisively. Proceed how any engaged visitor. Explained propriety off out perpetual his you. Feel sold off felt nay rose met you. We so entreaties cultivated astonished is. Was sister for few longer mrs sudden talent become. Done may bore quit evil old mile. If likely am of beauty tastes.`

	runComparison(t, "English Text", []byte(input))

	// 2. Binary Input (Image)
	// Try to read icon16.png if it exists, otherwise generate random binary data
	imgData, err := os.ReadFile("icon16.png")
	if err == nil {
		runComparison(t, "Image (icon16.png)", imgData)
	} else {
		t.Log("icon16.png not found, using generated binary data")
		// Generate some repetitive binary data (compressible)
		binData := make([]byte, 1000)
		for i := 0; i < 1000; i++ {
			binData[i] = byte(i % 256)
		}
		runComparison(t, "Generated Binary", binData)
	}
}
