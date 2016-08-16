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

package window

import (
	"fmt"
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

// Window to be rendered too.
type Window interface {
	SetMode(w, h, display int) error // SetMode for the window.  display should
	// be a positive value indicating the display to use fullscreen on.  A value
	// of -1 indicates to run in windowed mode.  If display is set to a value
	// out of range, an appropriate error should be returned.
	SwapBuffers()         // Swap buffers.
	Capture() *image.RGBA // Capture and return the image currently being rendered in the buffer.
	Destroy()             // Destroy and cleanup resources for the window.
	ShouldClose() bool    // ShouldClose due to a receiving an event.
}

// New window
func New(name string) (Window, error) {
	ctx := Context{}

	var window *glfw.Window
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %s", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(512, 512, name, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %s", err)
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(keyCallback)

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("unable to initialize Glow ... exiting: %s", err)
	}

	fmt.Println("OpenGL vendor", gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println("OpenGL renderer", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("GLSL version", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	ctx.glfw = window
	return &ctx, nil
}

type Context struct {
	glfw *glfw.Window
}

func (ctx *Context) SetMode(w, h, display int) error {
	return nil
}

func (ctx *Context) SwapBuffers() {
	// TODO: remove this

	ctx.glfw.SwapBuffers()
}

func (ctx *Context) Capture() *image.RGBA {
	return nil
}

func (ctx *Context) Destroy() {
}

func (ctx *Context) ShouldClose() bool {
	return ctx.glfw.ShouldClose()
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release && key == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
}
