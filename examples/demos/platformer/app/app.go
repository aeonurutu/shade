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

package app

import (
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/aeonurutu/shade/core/shader"
)

const ( // Program IDs
	progID      = iota
	numPrograms = iota
)

const ( // VAO Names
	trianglesName = iota
	numVAOs       = iota
)

const ( // Buffer Names
	arrayBufferName   = iota
	elementBufferName = iota
	numBuffers        = iota
)

const ( // Attrib Locations
	mcVertexLoc = 0
	mcColor     = 1 // TODO: rename to mcColorLoc
)

var (
	programs    [numPrograms]uint32
	vaos        [numVAOs]uint32
	numVertices [numVAOs]int32
	buffers     [numBuffers]uint32
)
var ( // Uniform Locations
	modelMatrixLoc      int32
	projectionMatrixLoc int32
)

var ( // App Settings
	modelMatrix      mgl32.Mat4
	projectionMatrix mgl32.Mat4
	rotation         float32
	aspect           float32
)

type App struct {
}

func New() *App {
	a := App{}
	return &a
}

func (a *App) Init() error {
	var err error
	aspect = float32(512) / float32(512)

	// Load the GLSL program
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "assets/basic.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "assets/basic.frag"},
	}
	programs[progID], err = shader.Load(&shaders)
	if err != nil {
		return err
	}
	gl.UseProgram(programs[progID])

	// Setup model to be rendered
	modelMatrixLoc = gl.GetUniformLocation(programs[progID], gl.Str("modelMatrix\x00"))
	projectionMatrixLoc = gl.GetUniformLocation(programs[progID], gl.Str("projectionMatrix\x00"))
	vertexPositions := []float32{
		-1.0, -1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0,
		1.0, -1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
	}
	numVertices[trianglesName] = int32(len(vertexPositions))
	vertexColors := []float32{
		1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 1.0,
		1.0, 0.0, 0.0, 1.0,
		0.0, 1.0, 1.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0,
	}
	vertexIndices := []uint16{
		0, 1, 2, 3, 6, 7, 4, 5, // First strip
		0xFFFF,                 // <<-- This is the restart index
		2, 6, 0, 4, 1, 5, 3, 7, // Second strip
	}
	sizeVertexIndices := len(vertexIndices) * int(unsafe.Sizeof(vertexIndices[0]))
	gl.GenBuffers(numBuffers, &buffers[0])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffers[elementBufferName])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeVertexIndices, gl.Ptr(vertexIndices), gl.STATIC_DRAW)

	sizeVertexPositions := len(vertexPositions) * int(unsafe.Sizeof(vertexPositions[0]))
	sizeVertexColors := len(vertexColors) * int(unsafe.Sizeof(vertexColors[0]))

	gl.GenVertexArrays(numVAOs, &vaos[0])
	gl.BindVertexArray(vaos[trianglesName])

	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[arrayBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertexPositions+sizeVertexColors, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, sizeVertexPositions, gl.Ptr(vertexPositions))
	gl.BufferSubData(gl.ARRAY_BUFFER, sizeVertexPositions, sizeVertexColors, gl.Ptr(vertexColors))

	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(sizeVertexPositions))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.ClearColor(0.9, 0.9, 0.9, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	rotation = 0

	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	return nil
}

func (a *App) ProcessInput() {
}

func (a *App) Update() {
	rotation += 0.0001 // TODO: fix dt
	Y := mgl32.Vec3{0, 1, 0}
	Z := mgl32.Vec3{0, 0, 1}
	modelMatrix = mgl32.Translate3D(0, 0, -5).Mul4(mgl32.HomogRotate3D(rotation*360, Y)).Mul4(mgl32.HomogRotate3D(rotation*720, Z))
	projectionMatrix = mgl32.Frustum(-1, 1, -aspect, aspect, 1, 500)
}

func (a *App) Render(window *glfw.Window, d time.Duration) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Render
	gl.UniformMatrix4fv(modelMatrixLoc, 1, false, &modelMatrix[0])
	gl.UniformMatrix4fv(projectionMatrixLoc, 1, false, &projectionMatrix[0])

	// Set up for a glDrawElements call
	gl.BindVertexArray(vaos[trianglesName])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffers[elementBufferName])

	gl.DrawElements(gl.TRIANGLE_STRIP, 8, gl.UNSIGNED_SHORT, nil)
	gl.DrawElements(gl.TRIANGLE_STRIP, 8, gl.UNSIGNED_SHORT, gl.PtrOffset(9*2)) // (const GLvoid *)(9 * sizeof(GLushort))

	gl.Flush()
	window.SwapBuffers()
}

func (a *App) Terminate() {
}
