export const InfoSection = () => (
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
        <p class="text-slate-400">Adds variety with time prefixes like "Today", "Later", "Then"</p>
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
);
