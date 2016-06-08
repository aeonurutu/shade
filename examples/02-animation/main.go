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
// Package app manages the main game loop.

package main

import (
	_ "image/png"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/convexbit/shade/camera"
	"github.com/convexbit/shade/display"
	"github.com/convexbit/shade/events"
	"github.com/convexbit/shade/sprite"
	"github.com/convexbit/shade/time/clock"
)

const windowWidth = 640
const windowHeight = 480

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	screen, err := display.SetMode("02-animation", windowWidth, windowHeight)
	if err != nil {
		log.Fatalln("failed to set display mode:", err)
	}

	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Bind(screen.Program)

	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	a, err := loadSprite("animation.png", 3, 1)
	if err != nil {
		panic(err)
	}
	a.Bind(screen.Program)
	var aframe float32 = 0.0

	for running := true; running; {
		dt := clock.Tick(30)
		screen.Fill(200.0/256.0, 200/256.0, 200/256.0)

		// TODO move this somewhere else (maybe a Clear method of display
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, event := range events.Get() {
			if event.Type == events.KeyUp && event.Key == glfw.KeyEscape {
				// Send window close event
				screen.Close()
			}
			if event.Type == events.WindowClose {
				// Handle window close
				running = false
			}
		}

		aframe += 0.003 * dt
		if int(aframe) > 2 {
			aframe = 0
		}

		frame := mgl32.Vec2{
			float32(int(aframe)), // Truncate to correct frame
			0}
		pos := mgl32.Vec3{
			windowWidth/2 - float32(a.Width)/2,
			windowHeight/2 - float32(a.Height)/2,
			0}
		a.DrawFrame(frame, pos, nil)

		screen.Flip()
		events.Poll()
	}

}

func loadSprite(path string, framesWide, framesHigh int) (*sprite.Context, error) {
	i, err := sprite.Load(path)
	if err != nil {
		return nil, err
	}
	s, err := sprite.New(i, nil, framesWide, framesHigh)
	if err != nil {
		return nil, err
	}

	return s, nil
}
