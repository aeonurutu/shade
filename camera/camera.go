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

// Package camera implements a simple camera that should meet most users needs.

package camera

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
)

// Context contains the camera's state
type Context struct {
	// Pos of the camera
	Pos mgl32.Vec3
	// Width of the camera in pixels (projected onto the screen)
	Width float32
	// Height of the camera in pixels (projected onto the screen)
	Height float32
	// Offset to move the camera when rendering
	Offset mgl32.Vec2
	// Top edge position of camera (after offset is applied)
	Top float32
	// Bottom edge position of camera (after offset is applied)
	Bottom float32
	// Left edge position of camera (after offset is applied)
	Left float32
	// Right edge position of camera (after offset is applied)
	Right float32
	// TopStop prevents the camera's Top position from exceeding it.
	TopStop float32
	// BottomStop prevents the camera's Bottom position from exceeding it.
	BottomStop float32
	// LeftStop prevents the camera's Left position from exceeding it.
	LeftStop float32
	// RightStop prevents the camera's Right position from exceeding it.
	RightStop float32

	// ProjMatrix for the current camera's position
	ProjMatrix mgl32.Mat4
	// ViewMatrix for the current camera's position
	ViewMatrix mgl32.Mat4

	// viewMatrixLoc in GLSL program
	viewMatrixLoc int32
}

// New camera is returned.
func New() (*Context, error) {
	c := Context{
		Pos:    mgl32.Vec3{},
		Width:  640,
		Height: 480,
	}
	return &c, nil
}

// Bind the camera to OpenGL
func (c *Context) Bind(program uint32) {
	var left, right, top, bottom, near, far float32
	right = float32(c.Width)
	top = float32(c.Height)
	near = 0.1
	far = 100.0
	c.ProjMatrix = mgl32.Ortho(left, right, bottom, top, near, far)
	projUniform := gl.GetUniformLocation(program, gl.Str("ProjMatrix\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &c.ProjMatrix[0])

	c.viewMatrixLoc = gl.GetUniformLocation(program, gl.Str("ViewMatrix\x00"))
	c.Move(mgl32.Vec3{})
	gl.UniformMatrix4fv(c.viewMatrixLoc, 1, false, &c.ViewMatrix[0])
}

// Move the camera to pos, unless that position conflicts with stops, in which case
// the camera will move as close to pos as possible.
func (c *Context) Move(pos mgl32.Vec3) {
	pos[0] -= c.Offset[0]
	pos[1] -= c.Offset[1]

	if c.LeftStop != 0 && pos[0] < c.LeftStop {
		pos[0] = c.LeftStop
	} else if c.RightStop != 0 && pos[0] > c.RightStop {
		pos[0] = c.RightStop
	}
	if c.TopStop != 0 && pos[1] > c.TopStop {
		pos[1] = c.TopStop
	} else if c.BottomStop != 0 && pos[1] < c.BottomStop {
		pos[1] = c.BottomStop
	}

	c.Left = c.Pos[0]
	c.Right = c.Pos[0] + c.Width
	c.Bottom = c.Pos[1]
	c.Top = c.Pos[1] + c.Height

	c.Pos = pos

	var eye, center, up mgl32.Vec3
	eye = mgl32.Vec3{c.Pos[0], c.Pos[1], 7.0}
	center = mgl32.Vec3{c.Pos[0], c.Pos[1], -1.0}
	up = mgl32.Vec3{0.0, 1.0, 0.0}
	c.ViewMatrix = mgl32.LookAtV(eye, center, up)

	gl.UniformMatrix4fv(c.viewMatrixLoc, 1, false, &c.ViewMatrix[0])
}

// Follow the pos uing a simple linear interpolation algorithm.  How fast the camera
// snaps to the position can be adjusted with lerp.  Like Move, it also obeys stops.
func (c *Context) Follow(pos mgl32.Vec3, lerp float32) {
	pos[0] -= c.Offset[0]
	pos[1] -= c.Offset[1]

	if c.LeftStop != 0 && pos[0] < c.LeftStop {
		pos[0] = c.LeftStop
	} else if c.RightStop != 0 && pos[0] > c.RightStop {
		pos[0] = c.RightStop
	}
	if c.TopStop != 0 && pos[1] > c.TopStop {
		pos[1] = c.TopStop
	} else if c.BottomStop != 0 && pos[1] < c.BottomStop {
		pos[1] = c.BottomStop
	}

	c.Left = c.Pos[0]
	c.Right = c.Pos[0] + c.Width
	c.Bottom = c.Pos[1]
	c.Top = c.Pos[1] + c.Height

	// TODO: refactor to call Move
	c.Pos[0] = c.Pos[0] + (pos[0]-c.Pos[0])*lerp
	c.Pos[1] = c.Pos[1] + (pos[1]-c.Pos[1])*lerp
	c.Pos[2] = c.Pos[2] + (pos[2]-c.Pos[2])*lerp

	var eye, center, up mgl32.Vec3
	eye = mgl32.Vec3{c.Pos[0], c.Pos[1], 7.0}
	center = mgl32.Vec3{c.Pos[0], c.Pos[1], -1.0}
	up = mgl32.Vec3{0.0, 1.0, 0.0}
	c.ViewMatrix = mgl32.LookAtV(eye, center, up)

	gl.UniformMatrix4fv(c.viewMatrixLoc, 1, false, &c.ViewMatrix[0])
}

// TODO: refactor entity interface so Update can be removed for cameras.
func (c *Context) Update(dt float32, g *[]entity.Entity) {
}

// TODO: refactor entity interface so Draw can be removed for cameras.
func (c Context) Draw() {
}
