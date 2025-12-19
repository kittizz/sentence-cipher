import { createDefaultCipher } from "sentence-cipher";

export const InfoSection = () => {
  const cipher = createDefaultCipher();
  const sampleInput = "Hello";
  const sampleBytes = new TextEncoder().encode(sampleInput);

  const stringExample = cipher.encodeString(sampleInput);
  const naturalExample = cipher.encodeNatural(sampleBytes);
  const binaryExample = cipher.encode(new Uint8Array([72, 101, 108, 108, 111]));

  return (
    <div class="mt-10 bg-slate-800/20 backdrop-blur-sm rounded-2xl border border-slate-700/30 p-6">
      <h2 class="text-xl font-semibold text-white mb-4">📚 How It Works</h2>
      <div class="grid md:grid-cols-3 gap-4 text-sm mb-6">
        <div class="bg-slate-800/30 rounded-xl p-4 border border-slate-700/50">
          <h3 class="text-purple-400 font-medium mb-2">📝 String Mode</h3>
          <p class="text-slate-400">
            Standard encoding - transforms text into office-themed sentences. Each sentence encodes
            up to 3 bytes.
          </p>
        </div>
        <div class="bg-slate-800/30 rounded-xl p-4 border border-slate-700/50">
          <h3 class="text-green-400 font-medium mb-2">📧 Natural Mode</h3>
          <p class="text-slate-400">
            Advanced encoding that structures data into a natural-looking email format with
            subjects, openers, and connectors.
          </p>
        </div>
        <div class="bg-slate-800/30 rounded-xl p-4 border border-slate-700/50">
          <h3 class="text-blue-400 font-medium mb-2">💾 Binary Mode</h3>
          <p class="text-slate-400">
            Raw byte encoding for any file type (images, docs, etc.), maintaining data integrity
            across sentences.
          </p>
        </div>
      </div>

      <div class="space-y-4">
        {/* String Example */}
        <div class="p-4 bg-slate-900/50 rounded-xl border border-slate-700/30">
          <div class="flex items-center gap-2 mb-2">
            <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-purple-500/20 text-purple-400 uppercase tracking-wider">
              String Mode Example
            </span>
          </div>
          <div class="grid md:grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-slate-500">Input:</span>
              <code class="block text-slate-300 mt-1">"{sampleInput}"</code>
            </div>
            <div>
              <span class="text-slate-500">Output:</span>
              <code class="block text-purple-400 mt-1 whitespace-pre-wrap">{stringExample}</code>
            </div>
          </div>
        </div>

        {/* Natural Example */}
        <div class="p-4 bg-slate-900/50 rounded-xl border border-slate-700/30">
          <div class="flex items-center gap-2 mb-2">
            <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-green-500/20 text-green-400 uppercase tracking-wider">
              Natural Mode Example
            </span>
          </div>
          <div class="grid md:grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-slate-500">Input:</span>
              <code class="block text-slate-300 mt-1">"{sampleInput}"</code>
            </div>
            <div>
              <span class="text-slate-500">Output:</span>
              <code class="block text-green-400 mt-1 whitespace-pre-wrap text-[11px] leading-relaxed">
                {naturalExample}
              </code>
            </div>
          </div>
        </div>

        {/* Binary Example */}
        <div class="p-4 bg-slate-900/50 rounded-xl border border-slate-700/30">
          <div class="flex items-center gap-2 mb-2">
            <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-blue-500/20 text-blue-400 uppercase tracking-wider">
              Binary Mode Example
            </span>
          </div>
          <div class="grid md:grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-slate-500">Input (Bytes):</span>
              <code class="block text-slate-300 mt-1">[72, 101, 108, 108, 111]</code>
            </div>
            <div>
              <span class="text-slate-500">Output:</span>
              <code class="block text-blue-400 mt-1 whitespace-pre-wrap">{binaryExample}</code>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
