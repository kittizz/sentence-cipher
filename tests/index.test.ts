import { describe, it, expect } from "vitest";
import {
  encode,
  decode,
  encodeString,
  decodeString,
  encodeNatural,
  decodeNatural,
  createCipher,
  createDefaultCipher,
  randomString,
} from "../dist/node";

// Helper to convert string to Uint8Array
function toBytes(str: string): Uint8Array {
  return new TextEncoder().encode(str);
}

describe("encode/decode (bytes)", () => {
  const tests = [
    { name: "empty", input: new Uint8Array([]) },
    { name: "single byte zero", input: new Uint8Array([0x00]) },
    { name: "single byte max", input: new Uint8Array([0xff]) },
    { name: "two bytes Hi", input: new Uint8Array([0x48, 0x69]) },
    { name: "hello", input: toBytes("Hello") },
    { name: "hello world", input: toBytes("Hello, World!") },
    { name: "numbers", input: toBytes("12345") },
    { name: "special chars", input: toBytes("!@#$%") },
    { name: "thai text", input: toBytes("สวัสดี") },
    { name: "mixed", input: toBytes("Hello สวัสดี 123") },
  ];

  for (const tt of tests) {
    it(`encode/decode ${tt.name}`, () => {
      const encoded = encode(tt.input);
      const decoded = decode(encoded) as Uint8Array;
      expect(decoded).toEqual(tt.input);
    });
  }
});

describe("encodeString/decodeString", () => {
  const tests = [
    { name: "empty", input: "" },
    { name: "simple", input: "Hello" },
    { name: "with spaces", input: "Hello World" },
    { name: "with punctuation", input: "Hello, World!" },
    { name: "thai", input: "สวัสดีครับ" },
    { name: "emoji", input: "Hello 👋" },
    { name: "secret message", input: "Meet at 9pm" },
    { name: "long text", input: "The quick brown fox jumps over the lazy dog" },
  ];

  for (const tt of tests) {
    it(`encodeString/decodeString ${tt.name}`, () => {
      const encoded = encodeString(tt.input);
      const decoded = decodeString(encoded);
      expect(decoded).toBe(tt.input);
    });
  }
});

describe("encodeNatural/decodeNatural (bytes)", () => {
  const tests = [
    { name: "empty", input: new Uint8Array([]) },
    { name: "hello", input: toBytes("Hello") },
    { name: "secret", input: toBytes("Secret message") },
    { name: "thai", input: toBytes("สวัสดี") },
  ];

  for (const tt of tests) {
    it(`encodeNatural/decodeNatural ${tt.name}`, () => {
      const encoded = encodeNatural(tt.input);
      const decoded = decodeNatural(encoded) as Uint8Array;
      expect(decoded).toEqual(tt.input);
    });
  }
});

describe("encoded looks like English", () => {
  it("should produce valid sentence structure", () => {
    const input = toBytes("Secret");
    const encoded = encode(input);

    expect(encoded.length).toBeGreaterThan(0);
    expect(encoded.endsWith(".")).toBe(true);
    expect(encoded).toContain(" ");
  });
});

describe("createCipher", () => {
  it("should create cipher with key", () => {
    const cipher = createCipher("test-key");
    const input = "Hello World";
    const encoded = cipher.encodeString(input);
    const decoded = cipher.decodeString(encoded);
    expect(decoded).toBe(input);
  });

  it("should encrypt/decrypt with cipher", () => {
    const cipher = createCipher("secret-key");
    const plaintext = "Secret message";
    const encrypted = cipher.encrypt(plaintext);
    const decrypted = cipher.decrypt(encrypted);
    expect(decrypted).toBe(plaintext);
  });
});

describe("createDefaultCipher", () => {
  it("should create default cipher with encodeString/decodeString", () => {
    const cipher = createDefaultCipher();
    const input = "Hello";
    const encoded = cipher.encodeString(input);
    const decoded = cipher.decodeString(encoded);
    expect(decoded).toBe(input);
  });
});

describe("randomString", () => {
  it("should generate random string with specified length", () => {
    const length = 16;
    const result = randomString(length);
    expect(result.length).toBe(length);
  });

  it("should generate different strings each time", () => {
    const str1 = randomString(32);
    const str2 = randomString(32);
    expect(str1).not.toBe(str2);
  });
});

describe("all byte values", () => {
  it("should encode/decode all possible byte values (0-255)", () => {
    for (let i = 0; i < 256; i++) {
      const input = new Uint8Array([i]);
      const encoded = encode(input);
      const decoded = decode(encoded) as Uint8Array;
      expect(decoded).toEqual(input);
    }
  });
});

describe("two byte combinations", () => {
  const samples = [
    [0x00, 0x00],
    [0xff, 0xff],
    [0x00, 0xff],
    [0xff, 0x00],
    [0x48, 0x69], // "Hi"
    [0xab, 0xcd],
    [0x12, 0x34],
  ];

  for (const [b1, b2] of samples) {
    it(`should encode/decode [${b1.toString(16)}, ${b2.toString(16)}]`, () => {
      const input = new Uint8Array([b1, b2]);
      const encoded = encode(input);
      const decoded = decode(encoded) as Uint8Array;
      expect(decoded).toEqual(input);
    });
  }
});

describe("randomString", () => {
  it("should generate random string with specified length", () => {
    const length = 16;
    const result = randomString(length);
    expect(result.length).toBe(length);
  });
});
