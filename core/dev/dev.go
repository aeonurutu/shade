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
// Package display TODO doc

/*
Package dev implements a developer window for relaying info to
developers/testers.
*/
package dev

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/aeonurutu/shade/core/gen"
	"github.com/aeonurutu/shade/core/shader"
	"github.com/aeonurutu/shade/core/util/archive"
)

const ( // Program IDs
	trianglesProgID = iota
	numPrograms     = iota
)

const ( // VAO Names
	trianglesName = iota
	numVAOs       = iota
)

const ( // Buffer Names
	arrayBufferName = iota
	numBuffers      = iota
)

const ( // Attrib Locations
	mcVertexLoc = 0
)

var (
	programs    [numPrograms]uint32
	vaos        [numVAOs]uint32
	numVertices [numVAOs]int32
	buffers     [numBuffers]uint32
)

// Context of the dev window.
type Context struct {
	Window *glfw.Window
}

// New dev context for a dev window is returned.
// TODO(hurricanerix): In the future it might be nice to be able to specify
// that you don't want a new window.  Instead dev info should be overlayed
// to the screen.
func New() *Context {
	var err error

	c := Context{}

	c.Window = createDevWindow()
	c.Window.MakeContextCurrent()

	vertSrc, err := archive.Get("./shaders/dev.vert")
	if err != nil {
		panic(err)
	}

	fragSrc, err := archive.Get("./shaders/dev.frag")
	if err != nil {
		panic(err)
	}

	// Load the GLSL program
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Source: bytes.NewReader(vertSrc)},
		shader.Info{Type: gl.FRAGMENT_SHADER, Source: bytes.NewReader(fragSrc)},
	}

	programs[trianglesProgID], err = shader.Load(&shaders)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(programs[trianglesProgID])

	// Setup model to be rendered
	vertices := []float32{
		-0.50, -0.50, // Triangle 1
		0.85, -0.50,
		-0.50, 0.85,
		0.50, -0.85, // Triangle 2
		0.50, 0.50,
		-0.85, 0.50,
	}

	numVertices[trianglesName] = int32(len(vertices))

	gl.GenVertexArrays(numVAOs, &vaos[0])
	gl.BindVertexArray(vaos[trianglesName])

	sizeVertices := len(vertices) * int(unsafe.Sizeof(vertices[0]))
	gl.GenBuffers(numBuffers, &buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[arrayBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertices, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(mcVertexLoc, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(mcVertexLoc)

	return &c
}

func (ctx *Context) Update() {

}

func (ctx *Context) Draw() {
	if ctx.Window != nil {
		ctx.Window.MakeContextCurrent()
	}
	gl.UseProgram(programs[trianglesProgID])
	// Clear buffer
	gl.ClearColor(1.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	// Render
	gl.BindVertexArray(vaos[trianglesName])
	gl.DrawArrays(gl.TRIANGLES, 0, numVertices[trianglesName])

	// TODO(hurricanerix): We actually want to print the following to the dev window.
	fmt.Printf("Shade Version: %s\n", gen.Version)
	fmt.Printf("Built from: %s/commit/%s\n", gen.GitURL, gen.Hash)

	// Swap Buffers
	gl.Flush()
	if ctx.Window != nil {
		ctx.Window.SwapBuffers()
	}
}

func createDevWindow() *glfw.Window {
	var devWindow *glfw.Window
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	devWindow, err := glfw.CreateWindow(256, 256, "dev", nil, nil)
	if err != nil {
		panic(fmt.Errorf("failed to create window: %s", err))
	}
	devWindow.SetPos(0, 0)
	if err = gl.Init(); err != nil {
		panic(fmt.Errorf("unable to initialize Glow ... exiting: %s", err))
	}
	devWindow.MakeContextCurrent()
	fmt.Println("NewWindow: ", "dev")
	fmt.Println("OpenGL vendor", gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println("OpenGL renderer", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("GLSL version", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	return devWindow
}
