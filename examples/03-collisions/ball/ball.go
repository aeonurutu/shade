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
// Package ball TODO doc

package ball

import (
	"math"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/aeonurutu/shade/entity"
	"github.com/aeonurutu/shade/shapes"
	"github.com/aeonurutu/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Ball TODO doc
type Ball struct {
	pos    mgl32.Vec3
	Sprite *sprite.Context
	Shape  shapes.Shape
	dx     float32
	dy     float32
}

// New TODO doc
func New(x, y, speed, angle float32, s *sprite.Context) Ball {
	// TODO should take a group in as a argument
	b := Ball{
		pos:    mgl32.Vec3{x, y, 1.0},
		Sprite: s,
		Shape:  *shapes.NewCircle(mgl32.Vec2{float32(s.Width) / 2, float32(s.Width) / 2}, float32(s.Width)/2),
	}
	b.dx = float32(math.Cos(float64(angle))) * speed
	b.dy = float32(math.Sin(float64(angle))) * speed
	return b
}

// Bind TODO doc
func (b *Ball) Bind(program uint32) error {
	return b.Sprite.Bind(program)
}

func (b Ball) Bounds() shapes.Shape {
	return b.Shape
}

func (b Ball) Pos() mgl32.Vec3 {
	return b.pos
}

// Update TODO doc
func (b *Ball) Update(dt float32, group *[]entity.Entity) {
	lastPos := mgl32.Vec3{b.pos[0], b.pos[1], b.pos[2]}
	switchDx := false
	switchDy := false

	b.pos[0] += b.dx * dt
	b.pos[1] += b.dy * dt

	var cgroup []entity.Collider
	for i := range *group {
		if c, ok := (*group)[i].(entity.Collider); ok {
			cgroup = append(cgroup, c)
		}
	}

	for _, c := range entity.Collide(b, &cgroup, false) {

		eb := c.Hit.Bounds()
		//ep := c.Hit.Pos()

		if math.Abs(float64(c.Dir[0])) > math.Abs(float64(c.Dir[1])) {
			switchDx = true
			if eb.Type == "circle" {
				c.Hit.(*Ball).dx *= -1
			}
		} else if math.Abs(float64(c.Dir[1])) > math.Abs(float64(c.Dir[0])) {
			switchDy = true
			if eb.Type == "circle" {
				c.Hit.(*Ball).dx *= -1
				c.Hit.(*Ball).dy *= -1
			}
		} else {
			switchDx = true
			switchDy = true
			if eb.Type == "circle" {
				c.Hit.(*Ball).dx *= -1
				c.Hit.(*Ball).dy *= -1
			}
		}
	}
	if switchDx {
		b.pos[0] = lastPos[0]
		b.dx *= -1
	}
	if switchDy {
		b.pos[1] = lastPos[1]
		b.dy *= -1
	}
}

// Draw TODO doc
func (b Ball) Draw() {
	b.Sprite.Draw(b.pos, nil)
}
