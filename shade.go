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

package shade

import (
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/hurricanerix/shade/core/window"
)

// Flags only allowed if compiled to allow developer mode.
var ldDevBuild string
var (
	devBuild bool
	fps      float64
)

// Flags always available.
var (
	version bool
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	var err error // To keep err from shadowing devBuild
	devBuild, err = strconv.ParseBool(ldDevBuild)
	if err != nil {
		devBuild = false
	}

	if devBuild {
		// Parse dev build specific flags.
		flag.Float64Var(&fps, "fps", -1, "frames per second")
	}
}

// EntryPoint interface applications use.
type EntryPoint interface {
	Init() error
	ProcessInput()
	Update()
	Render(time.Duration)
	Terminate()
}

// Engine handles main loop.
type Engine struct {
	Name   string        // Name of application.
	FPS    float64       // FPS targeted to run at.
	App    EntryPoint    // App for the engine to run.
	Window window.Window // Window to render to.
}

// New instance of Engine.
func New(name string) *Engine {
	e := Engine{
		Name: name,
		FPS:  29.97,
	}
	return &e
}

// SetFPS to target when rendering.
func (e *Engine) SetFPS(fps float64) {
	e.FPS = fps
}

// SetEntryPoint for application.
func (e *Engine) SetEntryPoint(app EntryPoint) {
	e.App = app
}

// Run application.
func (e *Engine) Run() error {
	flag.Parse()

	if devBuild {
		if fps != -1 {
			e.FPS = fps
		}

		fmt.Printf("Running at %3.2f FPS\n", e.FPS)
	}

	var err error
	e.Window, err = window.New(e.Name)
	if err != nil {
		return err
	}

	if err := e.App.Init(); err != nil {
		return err
	}

	var lag time.Duration
	var dpu time.Duration // Duration Per Update
	dpu, err = time.ParseDuration(fmt.Sprintf("%fs", 1.0/e.FPS))
	if err != nil {
		return err
	}

	previous := getPlayerTime()
	for !e.Window.ShouldClose() {
		current := getPlayerTime()
		elapsed := current.Sub(previous)
		previous = current
		lag += elapsed

		glfw.PollEvents()
		e.ProcessInput()
		e.App.ProcessInput()

		for lag >= dpu {
			e.Update()
			e.App.Update()
			lag -= dpu
		}

		e.Render(lag / dpu)
		e.App.Render(lag / dpu)
		gl.Flush()
		e.Window.SwapBuffers()
	}

	e.App.Terminate()

	return nil
}

func getPlayerTime() time.Time {
	return time.Now()
}

func (e *Engine) ProcessInput() {
}

func (e *Engine) Update() {
}

func (e *Engine) Render(d time.Duration) {
}
