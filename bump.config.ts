import { defineConfig } from "bumpp";
import { execSync } from "node:child_process";

export default defineConfig({
  printCommits: true,
  files: ["package.json", "cipher.go"],
  all: true,
  execute: (config) => {
    console.log("Running build...");
    execSync("pnpm --filter exports build", { stdio: "inherit" });
    execSync("git add .", { stdio: "inherit" });
  },
});
