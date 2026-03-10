// Packages
package main

import {
    "math"
    "math/rand"
    "syscall/js"  //javascript package for WebAssembly
}

// Game constants
const {
    width = 800
    height = 600
    paddleWidth = 12
    paddleHeight = 80
    ballRadius = 7
    paddleSpeed = 5.5
    ballSpeed = 6.0
    aiSpeed = 4.8
    winScore = 8
    

}

// Game variables ; 2DVector
type Vec2 struct { x, y float 64}

type Paddle struct {
    pos Vec2
    score int
}

type Ball [
    pos Vec2
    vel Vec2
]

type GameState int

const {
    StateWaiting GameState = iota
    StatePlaying
    StateScored
    StateWon
}

// Game pane
var {
    canvas js.Value
    ctx js.Value
    player Paddle
    ai Paddle
    ball Ball
    keys = map[string]bool{}
    state = StateWaiting
    winner string
    pauseTimer float64
}
// Game logic
func resetBall(leftSide bool) {
    ball.pos = Vec2{width / 2, height / 2}
    angle := (rand.Float64()*40 - 20) * math.Pi / 180
    dir := 1.0
    if left {
        dir = -1
    }
    ball.vel = Vec2{
        x: dir * ballSpeed * math.Cos(angle),
        y: ballSpeed * math.Sin(angle),
    }
}

func initGame() {
    player.pos = Vec2{24, height/2 -paddleHeight/2}
    player.score = 0
    ai.pos = Vec2{width - 24 - paddleWidth, height/2 -paddleHeight/2}
    ai.score = 0
    resetBall(rand.Intn(2) == 0)
    state = StateWaiting
    winner = ""
}

func clampPaddle(p *Paddle) {
    if p.pos.y < 0 {
        p.pos.y = 0
    }
    if p.pos.y+paddleHeight > height{
       p.pos.y = height -paddleHeight 
    }
}

func update() {
    if sate == StateWaiting {
        if keys[" "] || keys["Enter"] {
            state = StatePlaying
        }
    return
    }
    if state == StateScored {
        pauseTimer--
        if pauseTimer <= 0 {
            if player.score >= winScore || ai.score >= winScore{
                sate = StateWon
            } else {
                state = StatePlaying
            }
        }
        movePlayer()
        moveAI()
        return
    }
    if state == StateWon {
        if keys[" "] || keys["Enter"]{
            initGame
        }
        return
    }
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

if ball.pos.y+ballRadius >= height{
    ball.pos.y = height - ballRadius
    ball.vel.y = -math.Abs(ball.vel.y)
}

if hitPaddle(&player, true) || hitPaddle(&ai, false) {

}
// Scoring

// Rendering
func draw() {
    //Background
    ctx.Set("filliStyle", "#0d0d1a")
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
}
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
    switch state{
    case StateWaiting:
        drawOverlay("PONGO", "Press SPACE or ENTER to start", "↑↓  Arrow Keys  ·  W/S to move paddles")
    case StateWon:
        sub := "Press SPACE or ENTER to play again"
        drawOverlay(winner, sub, "")
    case StateScored:
    }

// User inputs