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
// Package fonts TODO doc

package fonts

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/aeonurutu/shade/core/util/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Context TODO doc
type Context struct {
	Sprite     *sprite.Context
	LocMap     map[int32]mgl32.Vec2
	UnknownLoc mgl32.Vec2
}

// New TODO doc
func New(s *sprite.Context, m map[int32]mgl32.Vec2, u mgl32.Vec2) (*Context, error) {
	c := Context{
		Sprite:     s,
		LocMap:     m,
		UnknownLoc: u,
	}
	return &c, nil
}

func SimpleASCII() (*Context, error) {
	path := fmt.Sprintf("%s/src/github.com/aeonurutu/shade/assets/font.png", os.Getenv("GOPATH"))
	i, err := sprite.Load(path)
	if err != nil {
		return nil, err
	}

	s, err := sprite.New(i, nil, 32, 3)
	if err != nil {
		return nil, err
	}

	m := make(map[int32]mgl32.Vec2, s.Width*s.Height)
	for y := float32(0); y < 3; y++ {
		for x := float32(0); x < 32; x++ {
			m[int32((y+1)*32+x)] = mgl32.Vec2{x, y}
		}
	}

	u := mgl32.Vec2{31, 1}

	f, err := New(s, m, u)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Bind TODO doc
func (c *Context) Bind(program uint32) {
	c.Sprite.Bind(program)
}

// DrawText TODO doc
//func (c Context) DrawText(x, y, sx, sy float32, color *mgl32.Vec4, msg string) {
func (c *Context) DrawText(pos mgl32.Vec3, e *sprite.Effects, msg string) {
	startX := pos[0]
	if e == nil {
		e = &sprite.Effects{
			Scale: mgl32.Vec3{1.0, 1.0, 1.0},
		}
	}

	for _, r := range msg {
		if l, ok := c.LocMap[r]; ok {
			c.Sprite.DrawFrame(l, pos, e)
			pos[0] += float32(c.Sprite.Width) * e.Scale[0]
		} else if r == 10 {
			pos[0] = startX
			pos[1] -= float32(c.Sprite.Height) * e.Scale[1]
		} else {
			c.Sprite.DrawFrame(c.UnknownLoc, pos, e)
			pos[0] += float32(c.Sprite.Width) * e.Scale[0]
		}
	}
}

// SizeText TODO doc
func (c Context) SizeText(e *sprite.Effects, msg string) (float32, float32) {
	if e == nil {
		e = &sprite.Effects{
			Scale: mgl32.Vec3{1.0, 1.0, 1.0},
		}
	}
	var lx float32 = 0.0
	var cx float32 = 0.0
	var cy float32 = float32(c.Sprite.Height) * e.Scale[1]
	for _, r := range msg {
		if r == 10 { // code for newline
			cx = 0
			cy += float32(c.Sprite.Height) * e.Scale[0]
		} else {
			cx += float32(c.Sprite.Width) * e.Scale[0]
		}
		if cx > lx {
			lx = cx
		}
	}
	return lx, cy
}
