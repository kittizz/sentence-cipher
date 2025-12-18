import { defineConfig, type UserConfig } from "tsdown";

const configs: UserConfig[] = [
  {
    entry: "index.ts",
    format: "esm",
    outDir: "dist/node",
    clean: true,
    minify: false,
    dts: true,
    treeshake: true,
    platform: "neutral",
    shims: true,
    sourcemap: false,
  },
];

export default defineConfig(configs);
