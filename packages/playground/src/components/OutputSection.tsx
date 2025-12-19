import { Show, Accessor } from "solid-js";
import { Mode, EncodingType } from "../types";

interface Props {
  mode: Accessor<Mode>;
  encodingType: Accessor<EncodingType>;
  input: Accessor<string>;
  output: Accessor<string>;
  error: Accessor<string | null>;
  copied: Accessor<boolean>;
  decodedFileData: Accessor<Uint8Array | null>;
  fileData: Accessor<Uint8Array | null>;
  handleCopy: () => void;
  handleDownloadDecoded: () => void;
  getDecodedPreviewUrl: () => string | null;
}

export const OutputSection = (props: Props) => (
  <div class="bg-slate-800/30 backdrop-blur-sm rounded-2xl border border-slate-700/50 p-5">
    <div class="flex items-center justify-between mb-3">
      <label class="text-white font-medium">
        {props.mode() === "encode"
          ? "📤 Encoded Sentences"
          : props.encodingType() === "binary"
            ? "📤 Decoded File"
            : "📤 Decoded Text"}
      </label>
      <div class="flex gap-2">
        {/* Download button for binary decode */}
        <Show
          when={
            props.encodingType() === "binary" &&
            props.mode() === "decode" &&
            props.decodedFileData()
          }
        >
          <button
            onClick={props.handleDownloadDecoded}
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
          onClick={props.handleCopy}
          disabled={!props.output() || !!props.error()}
          class={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${
            props.copied()
              ? "bg-green-600 text-white"
              : "bg-purple-600 hover:bg-purple-500 text-white disabled:opacity-50 disabled:cursor-not-allowed"
          }`}
        >
          {props.copied() ? (
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

    <Show when={props.error()}>
      <div class="bg-red-900/30 border border-red-700/50 rounded-xl px-4 py-3 mb-3">
        <p class="text-red-400 text-sm">❌ Error: {props.error()}</p>
      </div>
    </Show>

    {/* Binary decode with image preview */}
    <Show
      when={
        props.encodingType() === "binary" && props.mode() === "decode" && props.decodedFileData()
      }
    >
      <div class="space-y-3">
        <div
          class={`w-full bg-slate-900/50 border rounded-xl px-4 py-3 font-mono text-sm whitespace-pre-wrap ${
            props.error() ? "border-red-600/50" : "border-purple-500/50 text-green-400"
          }`}
        >
          {props.output()}
        </div>

        {/* Image preview for decoded data */}
        <Show when={props.getDecodedPreviewUrl()}>
          <div class="bg-slate-900/50 border border-slate-700 rounded-xl p-4">
            <p class="text-slate-400 text-sm mb-3">🖼️ Image Preview:</p>
            <div class="rounded-lg overflow-hidden border border-slate-700">
              <img
                src={props.getDecodedPreviewUrl() || ""}
                alt="Decoded preview"
                class="max-h-64 mx-auto object-contain"
              />
            </div>
          </div>
        </Show>
      </div>
    </Show>

    {/* Regular output */}
    <Show
      when={
        !(props.encodingType() === "binary" && props.mode() === "decode" && props.decodedFileData())
      }
    >
      <div
        class={`w-full min-h-40 bg-slate-900/50 border rounded-xl px-4 py-3 font-mono text-sm overflow-auto whitespace-pre-wrap ${
          props.error()
            ? "border-red-600/50"
            : props.output()
              ? "border-purple-500/50 text-green-400"
              : "border-slate-600 text-slate-500"
        }`}
      >
        {props.output() || (
          <span class="text-slate-500 italic">
            {props.mode() === "encode"
              ? "Your encoded sentences will appear here..."
              : "Your decoded message will appear here..."}
          </span>
        )}
      </div>
    </Show>

    <Show when={props.output() && !props.error()}>
      <div class="flex items-center justify-between mt-3 text-slate-500 text-sm">
        <span>{props.output().length} characters</span>
        <Show
          when={
            props.mode() === "encode" && !(props.encodingType() === "binary" && props.fileData())
          }
        >
          <span>
            ~{(props.output().length / Math.max(props.input().length, 1)).toFixed(1)}x expansion
          </span>
        </Show>
        <Show
          when={props.mode() === "encode" && props.encodingType() === "binary" && props.fileData()}
        >
          <span>
            ~{(props.output().length / Math.max(props.fileData()!.length, 1)).toFixed(1)}x expansion
          </span>
        </Show>
      </div>
    </Show>
  </div>
);
