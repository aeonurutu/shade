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
	"github.com/convexbit/shade/fonts"
	"github.com/convexbit/shade/sprite"
)

const windowWidth = 640
const windowHeight = 480

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	screen, err := display.SetMode("05-font", windowWidth, windowHeight)
	if err != nil {
		log.Fatalln("failed to set display mode:", err)
	}

	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Bind(screen.Program)

	font, err := fonts.SimpleASCII()
	if err != nil {
		panic(err)
	}
	font.Bind(screen.Program)

	var msg string
	var w, h float32

	for running := true; running; {
		// TODO move this somewhere else (maybe a Clear method of display
		screen.Fill(200.0/256.0, 200/256.0, 200/256.0)
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

		e := sprite.Effects{
			Scale: mgl32.Vec3{3.0, 3.0, 1.0},
		}
		msg = "Bottom Left"
		pos := mgl32.Vec3{0, 0, 0}
		font.DrawText(pos, &e, msg)

		msg = "Bottom Right"
		w, _ = font.SizeText(&e, msg)
		pos = mgl32.Vec3{screen.Width - w, 0, 0}
		font.DrawText(pos, &e, msg)

		msg = "Top Left"
		_, h = font.SizeText(&e, msg)
		pos = mgl32.Vec3{0, screen.Height - h, 0}
		font.DrawText(pos, &e, msg)

		msg = "Top Right"
		w, h = font.SizeText(&e, msg)
		pos = mgl32.Vec3{screen.Width - w, screen.Height - h, 0}
		font.DrawText(pos, &e, msg)

		msg = "Center\nMulti-Line\nText\nWith\nColor"
		w, h = font.SizeText(&e, msg)
		pos = mgl32.Vec3{screen.Width/2 - w/2, screen.Height/2 + h/2, 0}
		e.Tint = mgl32.Vec4{1, 0, 0, 0}
		font.DrawText(pos, &e, msg)

		screen.Flip()
		events.Poll()
	}

}
