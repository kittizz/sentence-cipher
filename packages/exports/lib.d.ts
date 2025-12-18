export interface Cipher {
  encode(data: Uint8Array): string;
  decode(encoded: string): Uint8Array;
  encodeString(s: string): string;
  decodeString(encoded: string): string;
  encodeNatural(data: Uint8Array): string;
  decodeNatural(encoded: string): Uint8Array;
  encrypt(plaintext: string): string;
  decrypt(encoded: string): string;
}

export function createCipher(key: string): Cipher;
export function createDefaultCipher(): Cipher;

export function encode(data: Uint8Array): string;
export function decode(encoded: string): Uint8Array;
export function encodeString(s: string): string;
export function decodeString(encoded: string): string;
export function encodeNatural(data: Uint8Array): string;
export function decodeNatural(encoded: string): Uint8Array;
