import { Show } from "solid-js";

interface Props {
  useKey: boolean;
  setUseKey: (useKey: boolean) => void;
  key: string;
  setKey: (key: string) => void;
  onGenerateRandomKey: () => void;
}

export const KeySection = (props: Props) => (
  <div class="bg-slate-800/30 backdrop-blur-sm rounded-2xl border border-slate-700/50 p-5 mb-6">
    <div class="flex items-center justify-between mb-3">
      <label class="flex items-center gap-2 cursor-pointer">
        <input
          type="checkbox"
          checked={props.useKey}
          onChange={(e) => props.setUseKey(e.currentTarget.checked)}
          class="w-4 h-4 rounded border-slate-600 bg-slate-700 text-purple-600 focus:ring-purple-500"
        />
        <span class="text-white font-medium">🔑 Use Encryption Key</span>
      </label>
      <button
        onClick={props.onGenerateRandomKey}
        class="text-sm text-purple-400 hover:text-purple-300 transition-colors"
      >
        Generate Random Key
      </button>
    </div>
    <Show when={props.useKey}>
      <input
        type="text"
        value={props.key}
        onInput={(e) => props.setKey(e.currentTarget.value)}
        placeholder="Enter your secret key..."
        class="w-full bg-slate-900/50 border border-slate-600 rounded-lg px-4 py-2.5 text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500 transition-all"
      />
      <p class="text-slate-500 text-xs mt-2">
        ⚠️ Remember your key! You'll need it to decode the message.
      </p>
    </Show>
  </div>
);
