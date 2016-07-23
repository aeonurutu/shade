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
// Package player TODO doc

package player

import (
	"math"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/aeonurutu/shade/core/entity"
	"github.com/aeonurutu/shade/core/events"
	"github.com/aeonurutu/shade/core/light"
	"github.com/aeonurutu/shade/core/shapes"
	"github.com/aeonurutu/shade/core/util/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Player struct {
	pos      mgl32.Vec3
	Shape    *shapes.Shape
	Sprite   *sprite.Context
	Light    *light.Positional
	Facing   float32
	Resting  bool
	Walking  bool
	dy       float32
	leftKey  bool
	rightKey bool
	jumpKey  bool
	whichLeg int
}

// New TODO doc
func New(x, y float32, s *sprite.Context) *Player {
	// TODO should take a group in as a argument
	p := Player{
		pos:    mgl32.Vec3{x, y, 1.0},
		Shape:  shapes.NewRect(0, 0, 24, 24),
		Sprite: s,
		Facing: 1,
	}
	light := light.Positional{
		Pos:   mgl32.Vec3{p.pos[0], float32(s.Height), 24.0},
		Color: mgl32.Vec4{0.7, 0.7, 1.0, 1.0},
		Power: 10000,
	}
	p.Light = &light
	return &p
}

func (p Player) Bounds() shapes.Shape {
	return *p.Shape
}

func (p Player) Pos() mgl32.Vec3 {
	return p.pos
}

// Handle TODO doc
func (p *Player) Handle(event events.Event) {
	// TODO: move this to SDK to handle things like holding Left & Right at the same time correctly

	if (event.Type == events.KeyDown || event.Type == events.KeyRepeat) && event.Key == glfw.KeyLeft {
		p.leftKey = true
	}
	if (event.Type == events.KeyDown || event.Type == events.KeyRepeat) && event.Key == glfw.KeyRight {
		p.rightKey = true
	}
	if (event.Type == events.KeyDown || event.Type == events.KeyRepeat) && event.Key == glfw.KeySpace {
		p.jumpKey = true
	}
	if event.Type == events.KeyUp && event.Key == glfw.KeyLeft {
		p.leftKey = false
	}
	if event.Type == events.KeyUp && event.Key == glfw.KeyRight {
		p.rightKey = false
	}
	if event.Type == events.KeyUp && event.Key == glfw.KeySpace {
		p.jumpKey = false
	}
}

// Bind TODO doc
func (p *Player) Bind(program uint32) error {
	return p.Sprite.Bind(program)
}

// Update TODO doc
func (p *Player) Update(dt float32, group *[]entity.Entity) {
	lastPos := mgl32.Vec3{p.pos[0], p.pos[1], p.pos[2]}
	p.Walking = false

	if p.leftKey {
		p.pos[0] -= 300.0 * dt
		p.Light.Pos[0] = p.pos[0]
		p.Facing = 0
		p.Walking = true
	}
	if p.rightKey {
		p.pos[0] += 300.0 * dt
		p.Facing = 1
		p.Light.Pos[0] = p.pos[0] + float32(p.Sprite.Width)
		p.Walking = true
	}
	if p.Resting && p.jumpKey {
		p.dy = 1500.0
	}
	p.dy = float32(math.Min(float64(1500.0), float64(p.dy-40.0)))

	p.pos[1] += p.dy * dt

	newPos := &p.pos
	p.Resting = false

	if p.pos[1] < 47 {
		p.Resting = true
		p.pos[1] = 48
		p.dy = 0.0
	}

	var cgroup []entity.Collider
	for i := range *group {
		if c, ok := (*group)[i].(entity.Collider); ok {
			cgroup = append(cgroup, c)
		}
	}
	collides := entity.Collide(p, &cgroup, false)
	for _, c := range collides {
		pos := c.Hit.Pos()
		s := c.Hit.Bounds().Data

		if (c.Dir[1] > -0.7 || c.Dir[1] < 0.7) && (c.Dir[0] < -0.7 || c.Dir[0] > 0.7) {
			newPos[0] = lastPos[0]
		}

		if c.Dir[0] > -0.7 || c.Dir[0] < 0.7 {
			if c.Dir[1] > 0.9 {
				// Hit top of tile
				newPos[1] = pos[1] - (s[3] + p.Shape.Data[3]) - 1
				p.dy = 0.0
			} else if c.Dir[1] < -0.9 {
				// Hit bottom of tile
				p.Resting = true
				newPos[1] = pos[1] + s[1] + 1
				p.dy = 0.0
			}
		}
	}
	p.Light.Pos[1] = p.pos[1] + float32(p.Sprite.Height)
}

// Draw TODO doc
func (p *Player) Draw() {
	if !p.Walking || !p.Resting {
		p.Sprite.DrawFrame(mgl32.Vec2{1, p.Facing}, p.pos, nil)
	} else {
		frame := float32(int(p.dy) % 2)
		switch {
		case p.whichLeg == 0:
			p.whichLeg = 2
		case p.whichLeg == 1:
			p.whichLeg = 4
		case p.whichLeg == 2:
			p.whichLeg = 3
		}
		p.Sprite.DrawFrame(mgl32.Vec2{frame + float32(p.whichLeg), p.Facing}, p.pos, nil)
	}
}
