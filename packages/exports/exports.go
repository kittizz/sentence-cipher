//go:build js

package main

import (
	"github.com/gopherjs/gopherjs/js"
	sentencecipher "github.com/kittizz/sentence-cipher"
)

var (
	global   = js.Global
	exports  = js.Module.Get("exports")
	errorObj = global.Get("Error")
	object   = global.Get("Object")
)

func throwError(err error) {
	panic(errorObj.New(err.Error()))
}

// createCipher สร้าง Cipher ด้วย key และ return object ที่มี methods พร้อมใช้
func createCipher(key string) *js.Object {
	cipher, err := sentencecipher.NewCipher(key)
	if err != nil {
		throwError(err)
	}
	return createCipherObject(cipher)
}

// createDefaultCipher สร้าง Cipher ด้วย default word lists (ไม่มี key)
func createDefaultCipher() *js.Object {
	cipher := sentencecipher.NewDefaultCipher()
	return createCipherObject(cipher)
}

// createCipherObject สร้าง JS object จาก Cipher
func createCipherObject(cipher *sentencecipher.Cipher) *js.Object {
	obj := object.New()

	// Basic encode/decode
	obj.Set("encode", func(data []byte) string {
		result, err := cipher.Encode(data)
		if err != nil {
			throwError(err)
		}
		return result
	})
	obj.Set("decode", func(encoded string) []byte {
		result, err := cipher.Decode(encoded)
		if err != nil {
			throwError(err)
		}
		return result
	})

	// String encode/decode
	obj.Set("encodeString", func(s string) string {
		result, err := cipher.EncodeString(s)
		if err != nil {
			throwError(err)
		}
		return result
	})
	obj.Set("decodeString", func(encoded string) string {
		result, err := cipher.DecodeString(encoded)
		if err != nil {
			throwError(err)
		}
		return result
	})

	// Natural encode/decode
	obj.Set("encodeNatural", func(data []byte) string {
		result, err := cipher.EncodeNatural(data)
		if err != nil {
			throwError(err)
		}
		return result
	})
	obj.Set("decodeNatural", func(encoded string) []byte {
		result, err := cipher.DecodeNatural(encoded)
		if err != nil {
			throwError(err)
		}
		return result
	})

	// Alias methods (encrypt/decrypt = encodeString/decodeString)
	obj.Set("encrypt", func(plaintext string) string {
		result, err := cipher.EncodeString(plaintext)
		if err != nil {
			throwError(err)
		}
		return result
	})
	obj.Set("decrypt", func(encoded string) string {
		result, err := cipher.DecodeString(encoded)
		if err != nil {
			throwError(err)
		}
		return result
	})

	return obj
}

// Package-level functions (backward compatibility with default cipher)
func encode(data []byte) string {
	result, err := sentencecipher.Encode(data)
	if err != nil {
		throwError(err)
	}
	return result
}

func decode(encoded string) []byte {
	result, err := sentencecipher.Decode(encoded)
	if err != nil {
		throwError(err)
	}
	return result
}

func encodeString(s string) string {
	result, err := sentencecipher.EncodeString(s)
	if err != nil {
		throwError(err)
	}
	return result
}

func decodeString(encoded string) string {
	result, err := sentencecipher.DecodeString(encoded)
	if err != nil {
		throwError(err)
	}
	return result
}

func encodeNatural(data []byte) string {
	result, err := sentencecipher.EncodeNatural(data)
	if err != nil {
		throwError(err)
	}
	return result
}

func getVersion() string {
	return sentencecipher.Version
}

func decodeNatural(encoded string) []byte {
	result, err := sentencecipher.DecodeNatural(encoded)
	if err != nil {
		throwError(err)
	}
	return result
}

func main() {

	// Cipher constructors
	exports.Set("createCipher", createCipher)
	exports.Set("createDefaultCipher", createDefaultCipher)

	// Package-level functions (use default cipher)
	exports.Set("encode", encode)
	exports.Set("decode", decode)
	exports.Set("encodeString", encodeString)
	exports.Set("decodeString", decodeString)
	exports.Set("encodeNatural", encodeNatural)
	exports.Set("decodeNatural", decodeNatural)
	exports.Set("getVersion", getVersion)

}
