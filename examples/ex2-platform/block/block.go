// Copyright 2016 Richard Hawkins
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
// Package block TODO doc

package block

import (
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Block struct {
	pos    mgl32.Vec3
	Sprite *sprite.Context
	Shape  *shapes.Shape
	Index  float32
}

// New TODO doc
func New(x, y float32, i float32, s *sprite.Context) *Block {
	b := Block{
		pos:    mgl32.Vec3{x, y, 1.0},
		Sprite: s,
		Shape:  shapes.NewRect(0, float32(s.Width), 0, float32(s.Height)),
		Index:  i,
	}
	return &b
}

func (b Block) Bounds() shapes.Shape {
	return *b.Shape
}

func (b Block) Pos() mgl32.Vec3 {
	return b.pos
}

// Bind TODO doc
func (b *Block) Bind(program uint32) error {
	return b.Sprite.Bind(program)
}

// Draw TODO doc
func (b Block) Draw() {
	//e *sprite.Effects) {
	//b.Sprite.Draw(b.Pos, e)
	b.Sprite.DrawFrame(mgl32.Vec2{b.Index, 0}, b.pos, nil)
}
