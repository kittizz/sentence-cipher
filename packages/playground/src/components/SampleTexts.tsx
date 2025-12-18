import { For, Show } from "solid-js";
import { Mode, EncodingType } from "../types";

interface Props {
  mode: Mode;
  encodingType: EncodingType;
  setInput: (text: string) => void;
}

const sampleTexts = ["Hello World", "Secret Message", "The quick brown fox", "🔐 Encrypted!"];

export const SampleTexts = (props: Props) => (
  <Show when={props.mode === "encode" && props.encodingType !== "binary"}>
    <div class="flex flex-wrap gap-2 mb-4">
      <span class="text-slate-500 text-sm">Try:</span>
      <For each={sampleTexts}>
        {(sample) => (
          <button
            onClick={() => props.setInput(sample)}
            class="text-sm px-3 py-1 rounded-full bg-slate-800/50 border border-slate-700 text-slate-400 hover:text-white hover:border-purple-500 transition-all"
          >
            {sample}
          </button>
        )}
      </For>
    </div>
  </Show>
);
