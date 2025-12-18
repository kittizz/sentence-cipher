import { Show, Accessor } from "solid-js";
import { Mode, EncodingType } from "../types";

interface Props {
  mode: Accessor<Mode>;
  encodingType: Accessor<EncodingType>;
  input: Accessor<string>;
  setInput: (val: string) => void;
  fileData: Accessor<Uint8Array | null>;
  fileName: Accessor<string>;
  fileType: Accessor<string>;
  handleFileUpload: (e: Event) => void;
  clearFile: () => void;
  getFilePreviewUrl: () => string | null;
}

export const InputSection = (props: Props) => (
  <div class="bg-slate-800/30 backdrop-blur-sm rounded-2xl border border-slate-700/50 p-5">
    <div class="flex items-center justify-between mb-3">
      <label class="text-white font-medium">
        {props.mode() === "encode"
          ? props.encodingType() === "binary"
            ? "📥 Input File/Image"
            : "📥 Input Text"
          : "📥 Encoded Sentences"}
      </label>
      <Show when={!(props.encodingType() === "binary" && props.mode() === "encode")}>
        <span class="text-slate-500 text-sm">{props.input().length} characters</span>
      </Show>
    </div>

    {/* Binary mode file upload (encode only) */}
    <Show when={props.encodingType() === "binary" && props.mode() === "encode"}>
      <div class="space-y-4">
        {/* File upload area */}
        <Show when={!props.fileData()}>
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
            <input type="file" class="hidden" onChange={props.handleFileUpload} />
          </label>
        </Show>

        {/* File preview */}
        <Show when={props.fileData()}>
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
                  <p class="text-white text-sm font-medium">{props.fileName()}</p>
                  <p class="text-slate-500 text-xs">
                    {props.fileData()!.length.toLocaleString()} bytes
                  </p>
                </div>
              </div>
              <button
                onClick={props.clearFile}
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
            <Show when={props.fileType().startsWith("image/")}>
              <div class="mt-3 rounded-lg overflow-hidden border border-slate-700">
                <img
                  src={props.getFilePreviewUrl() || ""}
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
    <Show when={!(props.encodingType() === "binary" && props.mode() === "encode")}>
      <textarea
        value={props.input()}
        onInput={(e) => props.setInput(e.currentTarget.value)}
        placeholder={
          props.mode() === "encode"
            ? "Enter your secret message..."
            : "Paste encoded sentences here..."
        }
        class="w-full h-40 bg-slate-900/50 border border-slate-600 rounded-xl px-4 py-3 text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500 transition-all resize-none font-mono text-sm"
      />
    </Show>
  </div>
);
