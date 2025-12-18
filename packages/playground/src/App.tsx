import { createSignal, createMemo, Show, For } from "solid-js";
import { createCipher, createDefaultCipher } from "sentence-cipher-exports";

const randomString = (length: number): string => {
  const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  let result = "";
  const randomValues = new Uint32Array(length);
  crypto.getRandomValues(randomValues);
  for (let i = 0; i < length; i++) {
    result += chars[randomValues[i] % chars.length];
  }
  return result;
};

type Mode = "encode" | "decode";
type EncodingType = "string" | "natural" | "binary";

const App = () => {
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

  const sampleTexts = ["Hello World", "Secret Message", "The quick brown fox", "🔐 Encrypted!"];

  return (
    <div class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <div class="container mx-auto px-4 py-8 max-w-4xl">
        {/* Header */}
        <header class="text-center mb-10">
          <h1 class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 bg-clip-text text-transparent mb-3">
            Sentence Cipher
          </h1>
          <p class="text-slate-400 text-lg">
            Hide your secrets in plain sight with natural-looking sentences
          </p>
          <a
            href="https://www.npmjs.com/package/sentence-cipher"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-2 mt-3 text-purple-400 hover:text-purple-300 transition-colors"
          >
            <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
              <path d="M0 7.334v8h6.666v1.332H12v-1.332h12v-8H0zm6.666 6.664H5.334v-4H3.999v4H1.335V8.667h5.331v5.331zm4 0v1.336H8.001V8.667h5.334v5.332h-2.669v-.001zm12.001 0h-1.33v-4h-1.336v4h-1.335v-4h-1.33v4h-2.671V8.667h8.002v5.331zM10.665 10H12v2.667h-1.335V10z" />
            </svg>
            sentence-cipher
          </a>
        </header>

        {/* Mode Toggle */}
        <div class="flex justify-center mb-8">
          <div class="bg-slate-800/50 backdrop-blur-sm rounded-xl p-1 border border-slate-700/50">
            <button
              onClick={() => setMode("encode")}
              class={`px-6 py-2.5 rounded-lg font-medium transition-all duration-200 ${
                mode() === "encode"
                  ? "bg-gradient-to-r from-purple-600 to-pink-600 text-white shadow-lg shadow-purple-500/25"
                  : "text-slate-400 hover:text-white"
              }`}
            >
              🔒 Encode
            </button>
            <button
              onClick={() => setMode("decode")}
              class={`px-6 py-2.5 rounded-lg font-medium transition-all duration-200 ${
                mode() === "decode"
                  ? "bg-gradient-to-r from-purple-600 to-pink-600 text-white shadow-lg shadow-purple-500/25"
                  : "text-slate-400 hover:text-white"
              }`}
            >
              🔓 Decode
            </button>
          </div>
        </div>

        {/* Encoding Type */}
        <div class="flex justify-center gap-3 mb-6">
          <For
            each={
              [
                { value: "string", label: "String", icon: "📝" },
                { value: "natural", label: "Natural", icon: "🌿" },
                { value: "binary", label: "Binary", icon: "💾" },
              ] as { value: EncodingType; label: string; icon: string }[]
            }
          >
            {(type) => (
              <button
                onClick={() => setEncodingType(type.value)}
                class={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 border ${
                  encodingType() === type.value
                    ? "bg-slate-700 border-purple-500 text-white"
                    : "bg-slate-800/50 border-slate-700 text-slate-400 hover:border-slate-600"
                }`}
              >
                {type.icon} {type.label}
              </button>
            )}
          </For>
        </div>

        {/* Key Section */}
        <div class="bg-slate-800/30 backdrop-blur-sm rounded-2xl border border-slate-700/50 p-5 mb-6">
          <div class="flex items-center justify-between mb-3">
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={useKey()}
                onChange={(e) => setUseKey(e.currentTarget.checked)}
                class="w-4 h-4 rounded border-slate-600 bg-slate-700 text-purple-600 focus:ring-purple-500"
              />
              <span class="text-white font-medium">🔑 Use Encryption Key</span>
            </label>
            <button
              onClick={handleGenerateRandomKey}
              class="text-sm text-purple-400 hover:text-purple-300 transition-colors"
            >
              Generate Random Key
            </button>
          </div>
          <Show when={useKey()}>
            <input
              type="text"
              value={key()}
              onInput={(e) => setKey(e.currentTarget.value)}
              placeholder="Enter your secret key..."
              class="w-full bg-slate-900/50 border border-slate-600 rounded-lg px-4 py-2.5 text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500 transition-all"
            />
            <p class="text-slate-500 text-xs mt-2">
              ⚠️ Remember your key! You'll need it to decode the message.
            </p>
          </Show>
        </div>

        {/* Sample Texts */}
        <Show when={mode() === "encode" && encodingType() !== "binary"}>
          <div class="flex flex-wrap gap-2 mb-4">
            <span class="text-slate-500 text-sm">Try:</span>
            <For each={sampleTexts}>
              {(sample) => (
                <button
                  onClick={() => setInput(sample)}
                  class="text-sm px-3 py-1 rounded-full bg-slate-800/50 border border-slate-700 text-slate-400 hover:text-white hover:border-purple-500 transition-all"
                >
                  {sample}
                </button>
              )}
            </For>
          </div>
        </Show>

        {/* Main Content */}
        <div class="grid gap-4 md:gap-6">
          {/* Input */}
          <div class="bg-slate-800/30 backdrop-blur-sm rounded-2xl border border-slate-700/50 p-5">
            <div class="flex items-center justify-between mb-3">
              <label class="text-white font-medium">
                {mode() === "encode"
                  ? encodingType() === "binary"
                    ? "📥 Input File/Image"
                    : "📥 Input Text"
                  : "📥 Encoded Sentences"}
              </label>
              <Show when={!(encodingType() === "binary" && mode() === "encode")}>
                <span class="text-slate-500 text-sm">{input().length} characters</span>
              </Show>
            </div>

            {/* Binary mode file upload (encode only) */}
            <Show when={encodingType() === "binary" && mode() === "encode"}>
              <div class="space-y-4">
                {/* File upload area */}
                <Show when={!fileData()}>
                  <label class="flex flex-col items-center justify-center w-full h-40 border-2 border-dashed border-slate-600 rounded-xl cursor-pointer bg-slate-900/50 hover:border-purple-500 hover:bg-slate-800/50 transition-all">
                    <div class="flex flex-col items-center justify-center pt-5 pb-6">
                      <svg
                        class="w-10 h-10 mb-3 text-slate-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                        />
                      </svg>
                      <p class="mb-2 text-sm text-slate-400">
                        <span class="font-semibold">Click to upload</span> or drag and drop
                      </p>
                      <p class="text-xs text-slate-500">Any file type supported</p>
                    </div>
                    <input type="file" class="hidden" onChange={handleFileUpload} />
                  </label>
                </Show>

                {/* File preview */}
                <Show when={fileData()}>
                  <div class="bg-slate-900/50 border border-slate-600 rounded-xl p-4">
                    <div class="flex items-center justify-between mb-3">
                      <div class="flex items-center gap-3">
                        <div class="w-10 h-10 bg-purple-600/20 rounded-lg flex items-center justify-center">
                          <svg
                            class="w-5 h-5 text-purple-400"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                            />
                          </svg>
                        </div>
                        <div>
                          <p class="text-white text-sm font-medium">{fileName()}</p>
                          <p class="text-slate-500 text-xs">
                            {fileData()!.length.toLocaleString()} bytes
                          </p>
                        </div>
                      </div>
                      <button
                        onClick={clearFile}
                        class="text-slate-400 hover:text-red-400 transition-colors"
                      >
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M6 18L18 6M6 6l12 12"
                          />
                        </svg>
                      </button>
                    </div>

                    {/* Image preview */}
                    <Show when={fileType().startsWith("image/")}>
                      <div class="mt-3 rounded-lg overflow-hidden border border-slate-700">
                        <img
                          src={getFilePreviewUrl() || ""}
                          alt="Preview"
                          class="max-h-48 mx-auto object-contain"
                        />
                      </div>
                    </Show>
                  </div>
                </Show>
              </div>
            </Show>

            {/* Regular text input */}
            <Show when={!(encodingType() === "binary" && mode() === "encode")}>
              <textarea
                value={input()}
                onInput={(e) => setInput(e.currentTarget.value)}
                placeholder={
                  mode() === "encode"
                    ? "Enter your secret message..."
                    : "Paste encoded sentences here..."
                }
                class="w-full h-40 bg-slate-900/50 border border-slate-600 rounded-xl px-4 py-3 text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500 transition-all resize-none font-mono text-sm"
              />
            </Show>
          </div>

          {/* Swap Button */}
          <div class="flex justify-center">
            <button
              onClick={handleSwap}
              disabled={!output() || !!error()}
              class="p-3 rounded-full bg-slate-800 border border-slate-700 text-slate-400 hover:text-white hover:border-purple-500 hover:bg-slate-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              title="Swap input and output"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
                />
              </svg>
            </button>
          </div>

          {/* Output */}
          <div class="bg-slate-800/30 backdrop-blur-sm rounded-2xl border border-slate-700/50 p-5">
            <div class="flex items-center justify-between mb-3">
              <label class="text-white font-medium">
                {mode() === "encode"
                  ? "📤 Encoded Sentences"
                  : encodingType() === "binary"
                    ? "📤 Decoded File"
                    : "📤 Decoded Text"}
              </label>
              <div class="flex gap-2">
                {/* Download button for binary decode */}
                <Show
                  when={encodingType() === "binary" && mode() === "decode" && decodedFileData()}
                >
                  <button
                    onClick={handleDownloadDecoded}
                    class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-all bg-green-600 hover:bg-green-500 text-white"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                      />
                    </svg>
                    Download
                  </button>
                </Show>
                <button
                  onClick={handleCopy}
                  disabled={!output() || !!error()}
                  class={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${
                    copied()
                      ? "bg-green-600 text-white"
                      : "bg-purple-600 hover:bg-purple-500 text-white disabled:opacity-50 disabled:cursor-not-allowed"
                  }`}
                >
                  {copied() ? (
                    <>
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M5 13l4 4L19 7"
                        />
                      </svg>
                      Copied!
                    </>
                  ) : (
                    <>
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                        />
                      </svg>
                      Copy
                    </>
                  )}
                </button>
              </div>
            </div>

            <Show when={error()}>
              <div class="bg-red-900/30 border border-red-700/50 rounded-xl px-4 py-3 mb-3">
                <p class="text-red-400 text-sm">❌ Error: {error()}</p>
              </div>
            </Show>

            {/* Binary decode with image preview */}
            <Show when={encodingType() === "binary" && mode() === "decode" && decodedFileData()}>
              <div class="space-y-3">
                <div
                  class={`w-full bg-slate-900/50 border rounded-xl px-4 py-3 font-mono text-sm ${
                    error() ? "border-red-600/50" : "border-purple-500/50 text-green-400"
                  }`}
                >
                  {output()}
                </div>

                {/* Image preview for decoded data */}
                <Show when={getDecodedPreviewUrl()}>
                  <div class="bg-slate-900/50 border border-slate-700 rounded-xl p-4">
                    <p class="text-slate-400 text-sm mb-3">🖼️ Image Preview:</p>
                    <div class="rounded-lg overflow-hidden border border-slate-700">
                      <img
                        src={getDecodedPreviewUrl() || ""}
                        alt="Decoded preview"
                        class="max-h-64 mx-auto object-contain"
                      />
                    </div>
                  </div>
                </Show>
              </div>
            </Show>

            {/* Regular output */}
            <Show when={!(encodingType() === "binary" && mode() === "decode" && decodedFileData())}>
              <div
                class={`w-full min-h-40 bg-slate-900/50 border rounded-xl px-4 py-3 font-mono text-sm overflow-auto ${
                  error()
                    ? "border-red-600/50"
                    : output()
                      ? "border-purple-500/50 text-green-400"
                      : "border-slate-600 text-slate-500"
                }`}
              >
                {output() || (
                  <span class="text-slate-500 italic">
                    {mode() === "encode"
                      ? "Your encoded sentences will appear here..."
                      : "Your decoded message will appear here..."}
                  </span>
                )}
              </div>
            </Show>

            <Show when={output() && !error()}>
              <div class="flex items-center justify-between mt-3 text-slate-500 text-sm">
                <span>{output().length} characters</span>
                <Show when={mode() === "encode" && !(encodingType() === "binary" && fileData())}>
                  <span>
                    ~{(output().length / Math.max(input().length, 1)).toFixed(1)}x expansion
                  </span>
                </Show>
                <Show when={mode() === "encode" && encodingType() === "binary" && fileData()}>
                  <span>
                    ~{(output().length / Math.max(fileData()!.length, 1)).toFixed(1)}x expansion
                  </span>
                </Show>
              </div>
            </Show>
          </div>
        </div>

        {/* Info Section */}
        <div class="mt-10 bg-slate-800/20 backdrop-blur-sm rounded-2xl border border-slate-700/30 p-6">
          <h2 class="text-xl font-semibold text-white mb-4">📚 How It Works</h2>
          <div class="grid md:grid-cols-3 gap-4 text-sm">
            <div class="bg-slate-800/30 rounded-xl p-4 border border-slate-700/50">
              <h3 class="text-purple-400 font-medium mb-2">📝 String Mode</h3>
              <p class="text-slate-400">
                Standard encoding - transforms text into office-themed sentences
              </p>
            </div>
            <div class="bg-slate-800/30 rounded-xl p-4 border border-slate-700/50">
              <h3 class="text-green-400 font-medium mb-2">🌿 Natural Mode</h3>
              <p class="text-slate-400">
                Adds variety with time prefixes like "Today", "Later", "Then"
              </p>
            </div>
            <div class="bg-slate-800/30 rounded-xl p-4 border border-slate-700/50">
              <h3 class="text-blue-400 font-medium mb-2">💾 Binary Mode</h3>
              <p class="text-slate-400">Raw byte encoding for arbitrary binary data</p>
            </div>
          </div>

          <div class="mt-6 p-4 bg-slate-900/50 rounded-xl border border-slate-700/30">
            <h3 class="text-white font-medium mb-2">💡 Example</h3>
            <div class="grid md:grid-cols-2 gap-4 text-sm">
              <div>
                <span class="text-slate-500">Input:</span>
                <code class="block text-purple-400 mt-1">"Hello"</code>
              </div>
              <div>
                <span class="text-slate-500">Output:</span>
                <code class="block text-green-400 mt-1">
                  "ruth trains isabella prints. carl cleans daily."
                </code>
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <footer class="mt-10 text-center text-slate-500 text-sm">
          <p>
            Built with{" "}
            <a
              href="https://www.solidjs.com/"
              target="_blank"
              rel="noopener noreferrer"
              class="text-purple-400 hover:text-purple-300"
            >
              SolidJS
            </a>{" "}
            •{" "}
            <a
              href="https://github.com/kittizz/sentence-cipher"
              target="_blank"
              rel="noopener noreferrer"
              class="text-purple-400 hover:text-purple-300"
            >
              GitHub
            </a>
          </p>
        </footer>
      </div>
    </div>
  );
};

export default App;
