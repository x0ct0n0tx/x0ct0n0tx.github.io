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

// Game pane

// Game variables ; 2DVector
type Vec2 struct { x, y float 64}

type Paddle struct {
    pos Vec2
    scrore int
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
func resetBall

func initGame

func clampPaddle

func update


// Game state

// Scoring

// Rendering

// User inputs