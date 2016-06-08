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

//go:generate ../../bindata.sh "../.."

import (
	_ "image/png"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/camera"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/shade"
	"github.com/hurricanerix/shade/splash"
	"github.com/hurricanerix/shade/sprite"
)

const windowWidth = 640
const windowHeight = 480

func init() {
	// GLFW event handling must run on the main OS thread
	// TODO(hurricanerix): We might be able to only do this in files that
	//                     actually handle input.
	runtime.LockOSThread()
}

func main() {
	shade.Run()

	// TODO(hurricanerix): Everything below this should be moved to a "handler"
	//                     that is exectured with shade.Run()
	screen, err := display.SetMode("01-basic", windowWidth, windowHeight)
	if err != nil {
		log.Fatalln("failed to set display mode:", err)
	}

	splash.Main(screen)

	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Bind(screen.Program)

	s, err := loadSprite("test-pattern128x128.png", 1, 1)
	if err != nil {
		panic(err)
	}
	s.Bind(screen.Program)

	for running := true; running; {
		screen.Fill(0.0, 0.0, 0.0)
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

		s.Draw(mgl32.Vec3{windowWidth/2 - float32(s.Width)/2, windowHeight/2 - float32(s.Height)/2, 0.0}, nil)

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
