// Packages
package main

import (
	"math"
	"math/rand"
	"syscall/js" //javascript package for WebAssembly
)

// Game constants
const (
	width        = 800
	height       = 600
	paddleWidth  = 12
	paddleHeight = 80
	ballRadius   = 7
	paddleSpeed  = 5.5
	ballSpeed    = 6.0
	aiSpeed      = 4.8
	winScore     = 8
)

// Game variables ; 2DVector
type Vec2 struct{ x, y float64 }

type Paddle struct {
	pos   Vec2
	score int
}

type Ball struct {
	pos Vec2
	vel Vec2
}

type GameState int

const (
	StateWaiting GameState = iota
	StatePlaying
	StateScored
	StateWon
)

// Game pane
var (
	canvas     js.Value
	ctx        js.Value
	player     Paddle
	ai         Paddle
	ball       Ball
	keys       = map[string]bool{}
	state      = StateWaiting
	winner     string
	pauseTimer float64
)

// Game logic
func resetBall(leftSide bool) {
	ball.pos = Vec2{width / 2, height / 2}
	angle := (rand.Float64()*40 - 20) * math.Pi / 180
	dir := 1.0
	if leftSide {
		dir = -1
	}
	ball.vel = Vec2{
		x: dir * ballSpeed * math.Cos(angle),
		y: ballSpeed * math.Sin(angle),
	}
}

func initGame() {
	player.pos = Vec2{24, height/2 - paddleHeight/2}
	player.score = 0
	ai.pos = Vec2{width - 24 - paddleWidth, height/2 - paddleHeight/2}
	ai.score = 0
	resetBall(rand.Intn(2) == 0)
	state = StateWaiting
	winner = ""
}

func clampPaddle(p *Paddle) {
	if p.pos.y < 0 {
		p.pos.y = 0
	}
	if p.pos.y+paddleHeight > height {
		p.pos.y = height - paddleHeight
	}
}

func update() {
	if state == StateWaiting {
		if keys[" "] || keys["Enter"] {
			state = StatePlaying
		}
		return
	}
	if state == StateScored {
		pauseTimer--
		if pauseTimer <= 0 {
			if player.score >= winScore || ai.score >= winScore {
				state = StateWon
			} else {
				state = StatePlaying
			}
		}
		movePlayer()
		moveAI()
		return
	}
	if state == StateWon {
		if keys[" "] || keys["Enter"] {
			initGame()
		}
		return
	}

	// Game state
	movePlayer()
	moveAI()

	ball.pos.x += ball.vel.x
	ball.pos.y += ball.vel.y

	if ball.pos.y-ballRadius <= 0 {
		ball.pos.y = ballRadius
		ball.vel.y = math.Abs(ball.vel.y)
	}

	if ball.pos.y+ballRadius >= height {
		ball.pos.y = height - ballRadius
		ball.vel.y = -math.Abs(ball.vel.y)
	}

	if hitPaddle(&player, true) || hitPaddle(&ai, false) {

	}
	// Scoring
	if ball.pos.x+ballRadius < 0 {
		ai.score++
		scored(true)
	} else if ball.pos.x-ballRadius > width {
		player.score++
		scored(false)
	}
}

// User inputs
func movePlayer() {
	if keys["ArrowUp"] || keys["w"] || keys["W"] {
		player.pos.y -= paddleSpeed
	}
	if keys["ArrowDown"] || keys["s"] || keys["S"] {
		player.pos.y += paddleSpeed
	}
	clampPaddle(&player)
}

func moveAI() {
	center := ai.pos.y + paddleHeight/2
	if center < ball.pos.y-4 {
		ai.pos.y += aiSpeed
	} else if center > ball.pos.y+4 {
		ai.pos.y -= aiSpeed
	}
	clampPaddle(&ai)
}

func hitPaddle(p *Paddle, isPlayer bool) bool {
	if isPlayer && ball.vel.x >= 0 {
		return false
	}
	if !isPlayer && ball.vel.x <= 0 {
		return false
	}

	bx, by := ball.pos.x, ball.pos.y
	px, py := p.pos.x, p.pos.y

	nearX := math.Max(px, math.Min(bx, px+paddleWidth))
	nearY := math.Max(py, math.Min(by, py+paddleHeight))
	dx, dy := bx-nearX, by-nearY

	if dx*dx+dy*dy > ballRadius*ballRadius {
		return false
	}

	ball.vel.x = -ball.vel.x

	rel := (by - (py + paddleHeight/2)) / (paddleHeight / 2)
	ball.vel.y = rel * ballSpeed * 1.5

	speed := math.Sqrt(ball.vel.x*ball.vel.x + ball.vel.y*ball.vel.y)
	if speed < ballSpeed {
		ball.vel.x *= ballSpeed / speed
		ball.vel.y *= ballSpeed / speed
	}

	if isPlayer {
		ball.pos.x = px + paddleWidth + ballRadius + 1
	} else {
		ball.pos.x = px - ballRadius - 1
	}
	return true
}

func scored(towardPlayer bool) {
	state = StateScored
	pauseTimer = 90

	if player.score >= winScore {
		winner = "YOU WIN!"
	} else if ai.score >= winScore {
		winner = "AI WINS!"
	} else {
		resetBall(towardPlayer)
	}
}

// Rendering
func draw() {
	//Background
	ctx.Set("fillStyle", "#0d0d1a")
	ctx.Call("fillRect", 0, 0, width, height)
	//Lines
	ctx.Set("strokeStyle", "rgba(255,255,255,0.15)")
	ctx.Set("lineWidth", 3)
	ctx.Call("setLineDash", js.ValueOf([]interface{}{12, 12}))
	ctx.Call("beginPath")
	ctx.Call("moveTo", width/2, 0)
	ctx.Call("lineTo", width/2, height)
	ctx.Call("stroke")
	ctx.Call("setLineDash", js.ValueOf([]interface{}{}))

	// Scores
	ctx.Set("font", "bold 56px 'Courier New', monospace")
	ctx.Set("textAlign", "center")
	ctx.Set("fillStyle", "rgba(255,255,255,0.9)")
	ctx.Call("fillText", player.score, width/4, 72)
	ctx.Call("fillText", ai.score, 3*width/4, 72)

	//Labels
	ctx.Set("font", "12px 'Courier New', monospace")
	ctx.Set("fillStyle", "rgba(255,255,255,0.4)")
	ctx.Call("fillText", "YOU", width/4, 95)
	ctx.Call("fillText", "AI", 3*width/4, 95)

	//Paddles
	drawPaddle(player.pos, "#e94560")
	drawPaddle(ai.pos, "#16c79a")

	if state != StateWaiting {
		drawBall()
	}

	//User messages
	switch state {
	case StateWaiting:
		drawOverlay("PONGO", "Press SPACE or ENTER to start", "↑↓  Arrow Keys  ·  W/S to move paddles")
	case StateWon:
		sub := "Press SPACE or ENTER to play again"
		drawOverlay(winner, sub, "")
	case StateScored:
	}

	if state == StatePlaying {
		ctx.Set("font", "11px 'Courier', monospace")
		ctx.Set("fillStyle", "rgba(255, 255, 255, 0.2)")
		ctx.Set("textAlign", "center")
		ctx.Call("fillText", "↑↓ / W S to move · first to "+itoa(winScore)+" wins", width/2, height-12)
	}
}

func drawPaddle(pos Vec2, color string) {
	r := 4.0
	x, y := pos.x, pos.y
	w, h := float64(paddleWidth), float64(paddleHeight)

	ctx.Set("fillStyle", color)
	ctx.Call("beginPath")
	ctx.Call("moveTo", x+r, y)
	ctx.Call("lineTo", x+w-r, y)
	ctx.Call("arcTo", x+w, y, x+w, y+r, r)
	ctx.Call("lineTo", x+w, y+h-r)
	ctx.Call("arcTo", x+w, y+h, x+w-r, y+h, r)
	ctx.Call("lineTo", x+r, y+h)
	ctx.Call("arcTo", x, y+h, x, y+h-r, r)
	ctx.Call("lineTo", x, y+r)
	ctx.Call("arcTo", x, y, x+r, y, r)
	ctx.Call("closePath")
	ctx.Call("fill")

	// Glow via shadow
	ctx.Set("shadowColor", color)
	ctx.Set("shadowBlur", 14)
	ctx.Call("fill")
	ctx.Set("shadowBlur", 0)
}

func drawBall() {
	ctx.Set("fillStyle", "#ffffff")
	ctx.Set("shadowColor", "#ffffff")
	ctx.Set("shadowBlur", 18)
	ctx.Call("beginPath")
	ctx.Call("arc", ball.pos.x, ball.pos.y, ballRadius, 0, math.Pi*2)
	ctx.Call("fill")
	ctx.Set("shadowBlur", 0)
}

func drawOverlay(title, sub1, sub2 string) {
	ctx.Set("fillStyle", "rgba(0,0,0,0.55)")
	ctx.Call("fillRect", 0, 0, width, height)

	ctx.Set("textAlign", "center")
	ctx.Set("fillStyle", "#ffffff")
	ctx.Set("font", "bold 64px 'Courier New', monospace")
	ctx.Call("fillText", title, width/2, height/2-30)

	ctx.Set("font", "18px 'Courier New', monospace")
	ctx.Set("fillStyle", "rgba(255,255,255,0.7)")
	ctx.Call("fillText", sub1, width/2, height/2+22)

	if sub2 != "" {
		ctx.Set("font", "13px 'Courier New', monospace")
		ctx.Set("fillStyle", "rgba(255,255,255,0.4)")
		ctx.Call("fillText", sub2, width/2, height/2+52)
	}
}

// itoa is a tiny int-to-string helper to avoid importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

// ── Entry point ───────────────────────────────────────────────────────────────

func gameLoop(_ js.Value, _ []js.Value) interface{} {
	update()
	draw()
	js.Global().Call("requestAnimationFrame", js.FuncOf(gameLoop))
	return nil
}

func main() {
	doc := js.Global().Get("document")
	canvas = doc.Call("getElementById", "gameCanvas")
	canvas.Set("width", width)
	canvas.Set("height", height)
	ctx = canvas.Call("getContext", "2d")

	initGame()

	// Keyboard handlers
	js.Global().Call("addEventListener", "keydown", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		key := args[0].Get("key").String()
		keys[key] = true
		if key == "ArrowUp" || key == "ArrowDown" || key == " " {
			args[0].Call("preventDefault")
		}
		return nil
	}))
	js.Global().Call("addEventListener", "keyup", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		keys[args[0].Get("key").String()] = false
		return nil
	}))

	// Start loop
	js.Global().Call("requestAnimationFrame", js.FuncOf(gameLoop))

	// Block forever so the Go runtime stays alive
	select {}
}
