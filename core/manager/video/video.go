// Copyright 2016-2017 Richard Hawkins
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

// package video TODO: doc
package video

import (
	"fmt"
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

// Context TODO: doc
type Context struct {
	window       *sdl.Window
	windowWidth  int
	windowHeight int
	context      sdl.GLContext
	QuitSDL      func()
}

// Init TODO: doc
func (ctx *Context) Init(width, height int) error {
	var err error

	// Set OpenGL core version 4.1 context
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	// var windowFlags uint32
	// if eng.Fullscreen {
	//   windowFlags |= sdl.WINDOW_FULLSCREEN
	// }

	ctx.window, err = sdl.CreateWindow("Cube", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		return err
	}

	ctx.windowWidth = width
	ctx.windowHeight = height

	ctx.context, err = sdl.GL_CreateContext(ctx.window)
	if err != nil {
		return err
	}
	sdl.GL_MakeCurrent(ctx.window, ctx.context)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		return err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	return nil
}

// WindowSize TODO: doc
func (ctx Context) WindowSize() (int, int) {
	return ctx.windowWidth, ctx.windowHeight
}

// Swap TODO: doc
func (ctx *Context) Swap() {
	sdl.GL_SwapWindow(ctx.window)
}

// Quit TODO: doc
func (ctx *Context) Quit() {
	ctx.window.Destroy()
	sdl.GL_DeleteContext(ctx.context)

	if ctx.QuitSDL == nil {
		return
	}
	ctx.QuitSDL()
}
