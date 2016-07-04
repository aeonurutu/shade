// Copyright 2016 Richard Hawkins, Alan Erwin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package player manages a player's state

package player

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/aeonurutu/shade/core/entity"
	"github.com/aeonurutu/shade/core/events"
	"github.com/aeonurutu/shade/core/shapes"
	"github.com/aeonurutu/shade/core/util/sprite"
)

const NumToWin = 5
const TopY = 450
const BottomY = 75

// Player state
type Player struct {
	pos        mgl32.Vec3
	Score      int
	Sprite     *sprite.Context
	Shape      shapes.Shape
	PlayerNum  int // player 1 or player 2
	PaddleSize int
	upKey      bool
	downKey    bool
	keyUp      glfw.Key
	keyDown    glfw.Key
	velocity   int
}

func New(playerNum int, x, y float32, s *sprite.Context) *Player {
	// create initial paddle
	p := Player{
		pos:        mgl32.Vec3{x, y, 0.0},
		Sprite:     s,
		PlayerNum:  playerNum,
		PaddleSize: 8,
		velocity:   6,
	}

	// set shape for collision
	p.SetShape()

	// assign keys to player
	if p.PlayerNum == 1 {
		p.keyUp = glfw.KeyQ
		p.keyDown = glfw.KeyA
	} else {
		p.keyUp = glfw.KeyP
		p.keyDown = glfw.KeyL
	}

	fmt.Println(fmt.Sprintf("Player %d created.", playerNum))
	return &p
}

func (p *Player) SetShape() {
	p.Shape = *shapes.NewRect(0, float32(p.Sprite.Width), 0, float32(p.Sprite.Height*(p.PaddleSize-2)))
}

func (p Player) Pos() mgl32.Vec3 {
	return p.pos
}

// Bind TODO doc
func (p *Player) Bind(program uint32) error {
	return p.Sprite.Bind(program)
}

func (p Player) Bounds() shapes.Shape {
	return p.Shape
}

func (p *Player) Handle(event events.Event) {
	// TODO: move this to SDK to handle things like holding Left & Right at the same time correctly

	if (event.Type == events.KeyDown || event.Type == events.KeyRepeat) && event.Key == p.keyUp {
		p.upKey = true
	}
	if (event.Type == events.KeyDown || event.Type == events.KeyRepeat) && event.Key == p.keyDown {
		p.downKey = true
	}

	if event.Type == events.KeyUp && event.Key == p.keyUp {
		p.upKey = false
	}
	if event.Type == events.KeyUp && event.Key == p.keyDown {
		p.downKey = false
	}
}

// Update(dt?, group?)
func (p *Player) Update(dt float32, group *[]entity.Entity) {
	posY := p.pos[1]
	if p.upKey && posY <= TopY {
		p.pos[1] += float32(p.velocity)
	}
	if p.downKey && posY >= BottomY {
		p.pos[1] -= float32(p.velocity)
	}
}

func (p Player) Draw() {
	posX := p.pos[0]
	posY := p.pos[1]

	// DrawFrame(frame to render, position, effect); postion 0,0 is bottom left corner of screen
	// draw top of paddle
	p.Sprite.DrawFrame(mgl32.Vec2{0, 0}, mgl32.Vec3{posX, posY, 0}, nil)

	// draw middle part(s) of paddle
	for i := 0; i < p.PaddleSize-2; i++ {
		// position of paddle middle parts are offset by player posY minus (PaddleSize * i + 1, i.e. index of loop + 1)
		midPosY := posY - float32(p.PaddleSize*(i+1))
		p.Sprite.DrawFrame(mgl32.Vec2{0, 1}, mgl32.Vec3{posX, midPosY, 0}, nil)
	}

	// draw bottom of paddle
	p.Sprite.DrawFrame(mgl32.Vec2{0, 2}, mgl32.Vec3{posX, posY - float32(p.PaddleSize*(NumToWin+2)), 0}, nil)
}
