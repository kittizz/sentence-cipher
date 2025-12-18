import { useCipher } from "./hooks/useCipher";
import { Header } from "./components/Header";
import { ModeToggle } from "./components/ModeToggle";
import { EncodingTypeToggle } from "./components/EncodingTypeToggle";
import { KeySection } from "./components/KeySection";
import { SampleTexts } from "./components/SampleTexts";
import { InputSection } from "./components/InputSection";
import { OutputSection } from "./components/OutputSection";
import { InfoSection } from "./components/InfoSection";
import { Footer } from "./components/Footer";

const App = () => {
  const {
    mode,
    setMode,
    encodingType,
    setEncodingType,
    input,
    setInput,
    key,
    setKey,
    useKey,
    setUseKey,
    copied,
    error,
    fileData,
    fileName,
    fileType,
    decodedFileData,
    output,
    handleCopy,
    handleGenerateRandomKey,
    handleSwap,
    handleFileUpload,
    clearFile,
    handleDownloadDecoded,
    getFilePreviewUrl,
    getDecodedPreviewUrl,
  } = useCipher();

  return (
    <div class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <div class="container mx-auto px-4 py-8 max-w-4xl">
        <Header />

        <ModeToggle mode={mode()} setMode={setMode} />

        <EncodingTypeToggle encodingType={encodingType()} setEncodingType={setEncodingType} />

        <KeySection
          useKey={useKey()}
          setUseKey={setUseKey}
          key={key()}
          setKey={setKey}
          onGenerateRandomKey={handleGenerateRandomKey}
        />

        <SampleTexts mode={mode()} encodingType={encodingType()} setInput={setInput} />

        <div class="grid gap-4 md:gap-6">
          <InputSection
            mode={mode}
            encodingType={encodingType}
            input={input}
            setInput={setInput}
            fileData={fileData}
            fileName={fileName}
            fileType={fileType}
            handleFileUpload={handleFileUpload}
            clearFile={clearFile}
            getFilePreviewUrl={getFilePreviewUrl}
          />

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

          <OutputSection
            mode={mode}
            encodingType={encodingType}
            input={input}
            output={output}
            error={error}
            copied={copied}
            decodedFileData={decodedFileData}
            fileData={fileData}
            handleCopy={handleCopy}
            handleDownloadDecoded={handleDownloadDecoded}
            getDecodedPreviewUrl={getDecodedPreviewUrl}
          />
        </div>

        <InfoSection />
        <Footer />
      </div>
    </div>
  );
};

export default App;
