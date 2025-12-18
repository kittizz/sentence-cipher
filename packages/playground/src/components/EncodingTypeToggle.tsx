import { For } from "solid-js";
import { EncodingType } from "../types";

interface Props {
  encodingType: EncodingType;
  setEncodingType: (type: EncodingType) => void;
}

const types: { value: EncodingType; label: string; icon: string }[] = [
  { value: "string", label: "String", icon: "📝" },
  { value: "natural", label: "Natural", icon: "🌿" },
  { value: "binary", label: "Binary", icon: "💾" },
];

export const EncodingTypeToggle = (props: Props) => (
  <div class="flex justify-center gap-3 mb-6">
    <For each={types}>
      {(type) => (
        <button
          onClick={() => props.setEncodingType(type.value)}
          class={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 border ${
            props.encodingType === type.value
              ? "bg-slate-700 border-purple-500 text-white"
              : "bg-slate-800/50 border-slate-700 text-slate-400 hover:border-slate-600"
          }`}
        >
          {type.icon} {type.label}
        </button>
      )}
    </For>
  </div>
);
