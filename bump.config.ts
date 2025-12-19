import { defineConfig } from "bumpp";
import fs from "node:fs";
import { execSync } from "node:child_process";

export default defineConfig({
  execute: (config) => {
    const version = config.state.newVersion;
    console.log(`Updating Go version to ${version}...`);

    const goFilePath = "cipher.go";
    const goFileContent = fs.readFileSync(goFilePath, "utf-8");
    const versionRegex = /const Version = ".*"/;

    if (versionRegex.test(goFileContent)) {
      const newContent = goFileContent.replace(versionRegex, `const Version = "${version}"`);
      if (newContent !== goFileContent) {
        fs.writeFileSync(goFilePath, newContent);
        console.log(`Updated cipher.go to version ${version}`);
      }
    } else {
      console.error("Could not find Version constant in cipher.go");
      process.exit(1);
    }

    console.log("Running build...");
    execSync("pnpm --filter exports build", { stdio: "inherit" });
    execSync("git add .", { stdio: "inherit" });
  },
});
