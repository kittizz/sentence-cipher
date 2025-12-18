import { Mode } from "../types";

interface Props {
  mode: Mode;
  setMode: (mode: Mode) => void;
}

export const ModeToggle = (props: Props) => (
  <div class="flex justify-center mb-8">
    <div class="bg-slate-800/50 backdrop-blur-sm rounded-xl p-1 border border-slate-700/50">
      <button
        onClick={() => props.setMode("encode")}
        class={`px-6 py-2.5 rounded-lg font-medium transition-all duration-200 ${
          props.mode === "encode"
            ? "bg-gradient-to-r from-purple-600 to-pink-600 text-white shadow-lg shadow-purple-500/25"
            : "text-slate-400 hover:text-white"
        }`}
      >
        🔒 Encode
      </button>
      <button
        onClick={() => props.setMode("decode")}
        class={`px-6 py-2.5 rounded-lg font-medium transition-all duration-200 ${
          props.mode === "decode"
            ? "bg-gradient-to-r from-purple-600 to-pink-600 text-white shadow-lg shadow-purple-500/25"
            : "text-slate-400 hover:text-white"
        }`}
      >
        🔓 Decode
      </button>
    </div>
  </div>
);
