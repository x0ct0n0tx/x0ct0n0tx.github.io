# A little ping pong
Making a pong game in Go compiled to WebAssembly, hosted on GitHub Pages. No backend, no dependencies.

## Controls
| Key | Action |
|-----|--------|
| ↑ / ↓ Arrow keys | Move paddle |
| W / S | Move paddle (alt) |
| Space / Enter | Start / restart |

First to **8 points** wins.

## How is it put together?
The game logic is written in Go and compiled to a `.wasm` binary that runs directly in the browser. Go talks to the browser's Canvas API via the `syscall/js` package to handle keyboard input and draw every frame. A small JavaScript shim (`wasm_exec.js`) ships with Go and acts as the bridge between the WASM binary and the browser. GitHub Actions automatically builds and deploys everything to GitHub Pages on every push to main.

## What does it need?
- Go 1.24+
- A local HTTP server (browsers block WASM over `file://`)

## Local test
Run in shell:
```bash
GOOS=js GOARCH=wasm go build -o main.wasm .
cp "$(find $(go env GOROOT) -name wasm_exec.js)" .
python3 -m http.server 8080
```
Then open http://localhost:8080 in your browser.
