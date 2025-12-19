import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Assume script is in /scripts, so project root is one level up
const rootDir = path.resolve(__dirname, '..');
const packageJsonPath = path.resolve(rootDir, 'package.json');
const goFilePath = path.resolve(rootDir, 'cipher.go');

console.log('Updating Go version...');

try {
  const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf-8'));
  const version = packageJson.version;

  let goFileContent = fs.readFileSync(goFilePath, 'utf-8');
  // Look for: const Version = "1.0.0"
  const versionRegex = /const Version = ".*"/;

  if (versionRegex.test(goFileContent)) {
    const newContent = goFileContent.replace(versionRegex, `const Version = "${version}"`);
    if (newContent !== goFileContent) {
      fs.writeFileSync(goFilePath, newContent);
      console.log(`Updated cipher.go to version ${version}`);
    } else {
        console.log(`cipher.go version is already ${version}`);
    }
  } else {
    console.error('Could not find Version constant in cipher.go');
    process.exit(1);
  }
} catch (error) {
  console.error('Error updating version:', error);
  process.exit(1);
}

