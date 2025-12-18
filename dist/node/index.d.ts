//#region src/node/lib.d.ts
interface Cipher {
  encode(data: Uint8Array): string;
  decode(encoded: string): Uint8Array;
  encodeString(s: string): string;
  decodeString(encoded: string): string;
  encodeNatural(data: Uint8Array): string;
  decodeNatural(encoded: string): Uint8Array;
  encrypt(plaintext: string): string;
  decrypt(encoded: string): string;
}
declare function createCipher$1(key: string): Cipher;
declare function createDefaultCipher$1(): Cipher;
declare function encode$1(data: Uint8Array): string;
declare function decode$1(encoded: string): Uint8Array;
declare function encodeString$1(s: string): string;
declare function decodeString$1(encoded: string): string;
declare function encodeNatural$1(data: Uint8Array): string;
declare function decodeNatural$1(encoded: string): Uint8Array;
declare function randomString$1(length: number): string;
//#endregion
//#region src/node/index.d.ts
declare const createCipher: typeof createCipher$1, createDefaultCipher: typeof createDefaultCipher$1, encode: typeof encode$1, decode: typeof decode$1, encodeString: typeof encodeString$1, decodeString: typeof decodeString$1, encodeNatural: typeof encodeNatural$1, decodeNatural: typeof decodeNatural$1, randomString: typeof randomString$1;
//#endregion
export { createCipher, createDefaultCipher, decode, decodeNatural, decodeString, encode, encodeNatural, encodeString, randomString };