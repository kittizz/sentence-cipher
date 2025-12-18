#!/usr/bin/env node

import { execSync, spawnSync } from "child_process";
import { rmSync, existsSync } from "fs";
import { resolve } from "path";

const GO119 = "go1.19.13";
const GOPHERJS_VERSION = "v1.19.0-beta2";
const DIST_DIR = "dist";

function run(cmd: string, options: Record<string, unknown> = {}): void {
  console.log(`> ${cmd}`);
  try {
    execSync(cmd, { stdio: "inherit", ...options });
  } catch (error) {
    process.exit((error as { status?: number }).status || 1);
  }
}

function getGoEnv(varName: string): string {
  const result = spawnSync(GO119, ["env", varName], { encoding: "utf-8" });
  return result.stdout.trim();
}

const commands: Record<string, () => void> = {
  setup() {
    console.log("📦 Installing Go toolchain and dependencies...\n");
    run(`go install "golang.org/dl/${GO119}@latest"`);
    run(`${GO119} download`);
    run(`${GO119} install "github.com/gopherjs/gopherjs@${GOPHERJS_VERSION}"`);
    run("corepack enable");
    run("corepack prepare pnpm@latest --activate");
    run("pnpm install");
    console.log("\n✅ Setup complete!");
  },

  clean() {
    console.log(`🧹 Cleaning ${DIST_DIR}...`);
    if (existsSync(DIST_DIR)) {
      rmSync(DIST_DIR, { recursive: true, force: true });
    }
    console.log("✅ Clean complete!");
  },

  "build:node"() {
    console.log("🔨 Building Node.js library with GopherJS...\n");
    commands.clean();

    const goroot = getGoEnv("GOROOT");
    const gopath = getGoEnv("GOPATH");
    const gopherjsBin = resolve(gopath, "bin", "gopherjs");

    run(`"${gopherjsBin}" build -m -o "./src/node/lib.js" ./src/node`, {
      env: { ...process.env, GOPHERJS_GOROOT: goroot },
    });

    console.log("\n✅ Node.js build complete!");
  },

  "build:ts"() {
    console.log("🔷 Building TypeScript...\n");
    run("pnpm tsdown");
    console.log("\n✅ TypeScript build complete!");
  },

  build() {
    console.log("🚀 Building all...\n");
    commands["build:node"]();
    commands["build:ts"]();
    console.log("\n✅ All builds complete!");
  },

  test() {
    console.log("🧪 Running tests...\n");
    run("pnpm test");
  },
};

const command = process.argv[2];

if (!command || !commands[command]) {
  console.log(`
Usage: node build.ts <command>

Commands:
  setup       Install Go toolchain for GopherJS and dependencies
  build       Build all (node + ts)
  build:node  Build Node.js library with GopherJS
  build:ts    Build TypeScript library
  clean       Clean dist folder
  test        Run tests
`);
  process.exit(command ? 1 : 0);
}

commands[command]();
