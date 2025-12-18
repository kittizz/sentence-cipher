export const Header = () => (
  <header class="text-center mb-10">
    <h1 class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 bg-clip-text text-transparent mb-3">
      Sentence Cipher Steganography Playground
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
);
