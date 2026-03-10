# A little ping pong
Trying out making a pong game in Go as a small portfolio piece.

## Controls

## How is it put together?

## What does it need?

## Local test
Run in shell:

GOOS=js GOARCH=wasm go build -o main.wasm .
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
python3 -m http.server 8080

Then in your browser open: http://localhost:8080