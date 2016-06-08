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
	"fmt"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/fonts"
	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Ball TODO doc
type Ball struct {
	pos    mgl32.Vec3
	Sprite *sprite.Context
	Font   fonts.Context
	Shape  *shapes.Shape
}

// New TODO doc
func New(x, y float32, s *sprite.Context, f fonts.Context) Ball {
	// TODO should take a group in as a argument
	b := Ball{
		pos:    mgl32.Vec3{x, y, 1.0},
		Sprite: s,
		Font:   f,
		Shape:  shapes.NewCircle(mgl32.Vec2{float32(s.Width) / 2, float32(s.Height) / 2}, float32(s.Width)/2),
	}

	return b
}

// Bind TODO doc
func (b *Ball) Bind(program uint32) error {
	return b.Sprite.Bind(program)
}

func (b Ball) Bounds() shapes.Shape {
	return *b.Shape
}

func (b Ball) Pos() mgl32.Vec3 {
	return b.pos
}

// Update TODO doc
func (b *Ball) Update(dt float32, g []entity.Collider) {
}

// Draw TODO doc
func (b Ball) Draw() {
	efx := sprite.Effects{
		Scale: mgl32.Vec3{2.0, 2.0, 1.0},
	}

	msg := fmt.Sprintf("Pos: (%.0f,%.0f)\n", b.pos[0], b.pos[1])
	msg += fmt.Sprintf("Data: [\n")
	msg += fmt.Sprintf("  Center: (%.0f, %.0f)\n", b.Shape.Data[0], b.Shape.Data[1])
	msg += fmt.Sprintf("  Radius: %.0f\n", b.Shape.Data[2])
	msg += fmt.Sprintf("]\n")
	_, h := b.Font.SizeText(&efx, msg)
	//b.Font.DrawText(mgl32.Vec3{b.pos[0] - w, b.pos[1] - 16, 0}, &efx, msg)
	b.Font.DrawText(mgl32.Vec3{b.pos[0], b.pos[1] + h + 7, 0}, &efx, msg)

	b.Sprite.Draw(b.pos, nil)
}
