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

// Package sprite manages images
package sprite

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/png" // register PNG decode
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/archive"
	"github.com/hurricanerix/shade/light"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Sprite TODO doc
type Sprite interface {
	Bind(prog uint32) error
	//Update(float32, []entity.Entity)
	Draw(pos mgl32.Vec3, efx *Effects)
}

// Context TODO doc
type Context struct {
	ColorMap        image.Image
	NormalMap       image.Image
	Width           int
	Height          int
	framesX         int
	framesY         int
	vao             uint32
	vbo             uint32
	texLoc          uint32
	normalLoc       uint32
	model           mgl32.Mat4
	modelMatrix     int32
	tex             mgl32.Mat3
	texMatrix       int32
	addColorLoc     int32
	addColor        int32
	aColorLoc       int32
	aColor          mgl32.Vec4
	subColorLoc     int32
	subColor        int32
	sColorLoc       int32
	sColor          mgl32.Vec4
	AmbientColorLoc int32
	AmbientColor    mgl32.Vec4
	LightPosLoc     int32
	LightColorLoc   int32
	LightPowerLoc   int32
	Light           light.Positional
}

// Load
func Load(path string) (image.Image, error) {
	if path == "" {
		return nil, nil
	}

	imgFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %v", path, err)
	}
	defer imgFile.Close()
	i, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("could not decode file %s: %v", path, err)
	}
	return i, nil
}

func LoadAsset(name string) (image.Image, error) {
	if name == "" {
		return nil, nil
	}

	imgFile, err := archive.Get(name)
	if err != nil {
		return nil, fmt.Errorf("could not load asset %s: %v", name, err)
	}
	i, _, err := image.Decode(bytes.NewReader(imgFile))
	if err != nil {
		return nil, fmt.Errorf("could not decode file %s: %v", name, err)
	}
	return i, nil
}

// New TODO doc
func New(colorMap, normalMap image.Image, framesX, framesY int) (*Context, error) {
	c := Context{
		ColorMap:     colorMap,
		NormalMap:    normalMap,
		framesX:      framesX,
		framesY:      framesY,
		AmbientColor: mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
	}

	if colorMap != nil {
		rgba := image.NewRGBA(colorMap.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			return nil, fmt.Errorf("color map has unsupported stride")
		}
		c.Width = int(float32(rgba.Rect.Size().X) / float32(framesX))
		c.Height = int(float32(rgba.Rect.Size().Y) / float32(framesY))
	}

	if normalMap != nil {
		rgba := image.NewRGBA(normalMap.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			return nil, fmt.Errorf("normal map has unsupported stride")
		}
		if colorMap == nil {
			c.Width = int(float32(rgba.Rect.Size().X) / float32(framesX))
			c.Height = int(float32(rgba.Rect.Size().Y) / float32(framesY))
		}
	}

	return &c, nil
}

// Bind TODO doc
func (c *Context) Bind(program uint32) error {
	if c.ColorMap != nil {
		rgba := image.NewRGBA(c.ColorMap.Bounds())

		draw.Draw(rgba, rgba.Bounds(), c.ColorMap, image.Point{0, 0}, draw.Src)

		gl.GenTextures(1, &c.texLoc)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, c.texLoc)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_BASE_LEVEL, 0)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix))

		if c.NormalMap != nil {
			rgba = image.NewRGBA(c.NormalMap.Bounds())

			draw.Draw(rgba, rgba.Bounds(), c.NormalMap, image.Point{0, 0}, draw.Src)
		}

		gl.GenTextures(1, &c.normalLoc)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, c.normalLoc)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_BASE_LEVEL, 0)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix))
	}

	gl.UseProgram(program)

	colorMap := gl.GetUniformLocation(program, gl.Str("ColorMap\x00"))
	gl.Uniform1i(colorMap, 0)

	normalMapLoc := gl.GetUniformLocation(program, gl.Str("NormalMap\x00"))
	gl.Uniform1i(normalMapLoc, 1)

	c.modelMatrix = gl.GetUniformLocation(program, gl.Str("ModelMatrix\x00"))
	gl.UniformMatrix4fv(c.modelMatrix, 1, false, &c.model[0])

	c.texMatrix = gl.GetUniformLocation(program, gl.Str("TexMatrix\x00"))
	gl.UniformMatrix3fv(c.texMatrix, 1, false, &c.tex[0])

	// Add color
	c.addColorLoc = gl.GetUniformLocation(program, gl.Str("AddColor\x00"))
	gl.Uniform1i(c.addColor, c.addColor)
	c.aColorLoc = gl.GetUniformLocation(program, gl.Str("AColor\x00"))
	gl.UniformMatrix3fv(c.aColorLoc, 1, false, &c.aColor[0])

	// Sub color
	c.subColorLoc = gl.GetUniformLocation(program, gl.Str("SubColor\x00"))
	gl.Uniform1i(c.subColor, c.subColor)
	c.sColorLoc = gl.GetUniformLocation(program, gl.Str("SColor\x00"))
	gl.UniformMatrix3fv(c.sColorLoc, 1, false, &c.sColor[0])

	c.AmbientColorLoc = gl.GetUniformLocation(program, gl.Str("AmbientColor\x00"))
	gl.Uniform4fv(c.AmbientColorLoc, 1, &c.AmbientColor[0])

	c.LightPosLoc = gl.GetUniformLocation(program, gl.Str("LightPos\x00"))
	gl.Uniform3fv(c.LightPosLoc, 1, &c.Light.Pos[0])

	c.LightColorLoc = gl.GetUniformLocation(program, gl.Str("LightColor\x00"))
	gl.Uniform4fv(c.LightColorLoc, 1, &c.Light.Color[0])

	c.LightPowerLoc = gl.GetUniformLocation(program, gl.Str("LightPower\x00"))
	gl.Uniform1f(c.LightPowerLoc, c.Light.Power)

	if c.vao == 0 {
		gl.GenVertexArrays(1, &c.vao)
		gl.BindVertexArray(c.vao)
	}

	if c.vbo == 0 {
		gl.GenBuffers(1, &c.vbo)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	mcVertex := uint32(gl.GetAttribLocation(program, gl.Str("MCVertex\x00")))
	gl.EnableVertexAttribArray(mcVertex)
	gl.VertexAttribPointer(mcVertex, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(0))

	mcNormal := uint32(gl.GetAttribLocation(program, gl.Str("MCNormal\x00")))
	gl.EnableVertexAttribArray(mcNormal)
	gl.VertexAttribPointer(mcNormal, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(3*4))

	mcTangent := uint32(gl.GetAttribLocation(program, gl.Str("MCTangent\x00")))
	gl.EnableVertexAttribArray(mcTangent)
	gl.VertexAttribPointer(mcTangent, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(6*4))

	texCoord0 := uint32(gl.GetAttribLocation(program, gl.Str("TexCoord0\x00")))
	gl.EnableVertexAttribArray(texCoord0)
	gl.VertexAttribPointer(texCoord0, 2, gl.FLOAT, false, 11*4, gl.PtrOffset(9*4))

	return nil
}

// Draw TODO doc
func (c *Context) Draw(pos mgl32.Vec3, e *Effects) {
	c.DrawFrame(mgl32.Vec2{0, 0}, pos, e)
}

type Effects struct {
	Scale          mgl32.Vec3
	Tint           mgl32.Vec4
	EnableLighting bool
	AmbientColor   mgl32.Vec4
	Light          light.Positional
}

// DrawFrame TODO doc
//func (c *Context) DrawFrame(fx, fy int, sx, sy, px, py float32, addColor, subColor, ambientColor *mgl32.Vec4, light *light.Positional) {
func (c *Context) DrawFrame(frame mgl32.Vec2, pos mgl32.Vec3, e *Effects) {
	if e == nil {
		// Default effects
		e = &Effects{
			Scale: mgl32.Vec3{1.0, 1.0, 1.0},
		}
	}
	c.model = mgl32.Ident4()
	c.model = c.model.Mul4(mgl32.Translate3D(float32(c.Width*int(e.Scale[0]))/2.0, float32(c.Height*int(e.Scale[1]))/2.0, 0.0))
	c.model = c.model.Mul4(mgl32.Translate3D(pos[0], pos[1], pos[2]))
	c.model = c.model.Mul4(mgl32.Scale3D(float32(c.Width)*e.Scale[0], float32(c.Height)*e.Scale[1], 0.0))
	gl.UniformMatrix4fv(c.modelMatrix, 1, false, &c.model[0])

	c.tex = mgl32.Ident3()
	c.tex = c.tex.Mul3(mgl32.Scale2D(1.0/float32(c.framesX), 1.0/float32(c.framesY)))
	c.tex = c.tex.Mul3(mgl32.Translate2D(frame[0], frame[1]))
	gl.UniformMatrix3fv(c.texMatrix, 1, false, &c.tex[0])

	// TODO change addColor to MaxColor in shader
	ac := int32(1)
	c.aColor = e.Tint
	gl.Uniform4fv(c.aColorLoc, 1, &c.aColor[0])
	c.addColor = ac
	gl.Uniform1i(c.addColorLoc, c.addColor)

	/*
		//MinColor       mgl32.Vec4
		sc := int32(0)
		if subColor != nil {
			sc = 1
			c.sColor = *subColor
			gl.Uniform4fv(c.sColorLoc, 1, &c.sColor[0])
		}
		c.subColor = sc
		gl.Uniform1i(c.subColorLoc, c.subColor)
	*/
	if e.EnableLighting == true {
		c.AmbientColor = e.AmbientColor
		gl.Uniform4fv(c.AmbientColorLoc, 1, &c.AmbientColor[0])

		c.Light = e.Light
		gl.Uniform3fv(c.LightPosLoc, 1, &c.Light.Pos[0])
		gl.Uniform4fv(c.LightColorLoc, 1, &c.Light.Color[0])
		gl.Uniform1f(c.LightPowerLoc, c.Light.Power)
	}

	gl.BindVertexArray(c.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, c.texLoc)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, c.normalLoc)

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}

// Update TODO doc
func (c *Context) Update(dt float32) {
}

// Pos(X, Y, Z), Normal(X, Y, Z), Tangent(X, Y, Z), TextureCo(S, T)
var vertices = []float32{
	-0.5, -0.5, -0.5, 0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 0.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, 1.0,
	0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 0.0, 1.0,
	0.5, 0.5, -0.5, 0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, 0.0,
}
