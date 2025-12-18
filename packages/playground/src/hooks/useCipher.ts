import { createSignal, createMemo } from "solid-js";
import { createCipher, createDefaultCipher } from "sentence-cipher";
import { Mode, EncodingType } from "../types";
import { randomString } from "../utils/crypto";

export const useCipher = () => {
  const [mode, setMode] = createSignal<Mode>("encode");
  const [encodingType, setEncodingType] = createSignal<EncodingType>("string");
  const [input, setInput] = createSignal("");
  const [key, setKey] = createSignal("");
  const [useKey, setUseKey] = createSignal(false);
  const [copied, setCopied] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);
  const [fileData, setFileData] = createSignal<Uint8Array | null>(null);
  const [fileName, setFileName] = createSignal<string>("");
  const [fileType, setFileType] = createSignal<string>("");
  const [decodedFileData, setDecodedFileData] = createSignal<Uint8Array | null>(null);

  const cipher = createMemo(() => {
    if (useKey() && key().trim()) {
      return createCipher(key().trim());
    }
    return createDefaultCipher();
  });

  const output = createMemo(() => {
    const currentType = encodingType();
    const currentMode = mode();

    // Binary mode with file
    if (currentType === "binary" && currentMode === "encode" && fileData()) {
      setError(null);
      try {
        const c = cipher();
        return c.encode(fileData()!);
      } catch (e) {
        const errorMessage = e instanceof Error ? e.message : "Unknown error";
        setError(errorMessage);
        return "";
      }
    }

    // Binary mode decode
    if (currentType === "binary" && currentMode === "decode") {
      const inputText = input().trim();
      if (!inputText) {
        setDecodedFileData(null);
        return "";
      }
      setError(null);
      try {
        const c = cipher();
        const decoded = c.decode(inputText);
        setDecodedFileData(decoded);
        return `[Binary data: ${decoded.length} bytes]`;
      } catch (e) {
        const errorMessage = e instanceof Error ? e.message : "Unknown error";
        setError(errorMessage);
        setDecodedFileData(null);
        return "";
      }
    }

    const inputText = input().trim();
    if (!inputText) return "";

    setError(null);

    try {
      const c = cipher();

      if (currentMode === "encode") {
        if (currentType === "string") {
          return c.encodeString(inputText);
        } else if (currentType === "natural") {
          return c.encodeNatural(new TextEncoder().encode(inputText));
        } else {
          return c.encode(new TextEncoder().encode(inputText));
        }
      } else {
        if (currentType === "string") {
          return c.decodeString(inputText);
        } else if (currentType === "natural") {
          const decoded = c.decodeNatural(inputText);
          return new TextDecoder().decode(decoded);
        } else {
          const decoded = c.decode(inputText);
          return new TextDecoder().decode(decoded);
        }
      }
    } catch (e) {
      const errorMessage = e instanceof Error ? e.message : "Unknown error";
      setError(errorMessage);
      return "";
    }
  });

  const handleCopy = async () => {
    const text = output();
    if (!text) return;

    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Fallback
      const textarea = document.createElement("textarea");
      textarea.value = text;
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand("copy");
      document.body.removeChild(textarea);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleGenerateRandomKey = () => {
    setKey(randomString(16));
    setUseKey(true);
  };

  const handleSwap = () => {
    const currentOutput = output();
    if (currentOutput && !error()) {
      // Don't swap for binary file mode
      if (encodingType() === "binary" && fileData()) {
        setInput(currentOutput);
        setFileData(null);
        setFileName("");
        setFileType("");
        setMode("decode");
      } else {
        setInput(currentOutput);
        setMode(mode() === "encode" ? "decode" : "encode");
      }
    }
  };

  const handleFileUpload = (e: Event) => {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;

    setFileName(file.name);
    setFileType(file.type);

    const reader = new FileReader();
    reader.onload = () => {
      const arrayBuffer = reader.result as ArrayBuffer;
      setFileData(new Uint8Array(arrayBuffer));
      setInput(""); // Clear text input when file is uploaded
    };
    reader.readAsArrayBuffer(file);
  };

  const clearFile = () => {
    setFileData(null);
    setFileName("");
    setFileType("");
  };

  const handleDownloadDecoded = () => {
    const data = decodedFileData();
    if (!data) return;

    const blob = new Blob([new Uint8Array(data)]);
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "decoded_file";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const getFilePreviewUrl = () => {
    const data = fileData();
    const type = fileType();
    if (!data || !type.startsWith("image/")) return null;
    const blob = new Blob([new Uint8Array(data)], { type });
    return URL.createObjectURL(blob);
  };

  const getDecodedPreviewUrl = () => {
    const data = decodedFileData();
    if (!data) return null;
    // Try to detect if it's an image by checking magic bytes
    const isJpeg = data[0] === 0xff && data[1] === 0xd8;
    const isPng = data[0] === 0x89 && data[1] === 0x50 && data[2] === 0x4e && data[3] === 0x47;
    const isGif = data[0] === 0x47 && data[1] === 0x49 && data[2] === 0x46;
    const isWebp = data[8] === 0x57 && data[9] === 0x45 && data[10] === 0x42 && data[11] === 0x50;

    let mimeType = "";
    if (isJpeg) mimeType = "image/jpeg";
    else if (isPng) mimeType = "image/png";
    else if (isGif) mimeType = "image/gif";
    else if (isWebp) mimeType = "image/webp";

    if (!mimeType) return null;
    const blob = new Blob([new Uint8Array(data)], { type: mimeType });
    return URL.createObjectURL(blob);
  };

  return {
    mode,
    setMode,
    encodingType,
    setEncodingType,
    input,
    setInput,
    key,
    setKey,
    useKey,
    setUseKey,
    copied,
    error,
    fileData,
    setFileData,
    fileName,
    fileType,
    decodedFileData,
    output,
    handleCopy,
    handleGenerateRandomKey,
    handleSwap,
    handleFileUpload,
    clearFile,
    handleDownloadDecoded,
    getFilePreviewUrl,
    getDecodedPreviewUrl,
  };
};
