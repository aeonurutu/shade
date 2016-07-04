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

	"github.com/aeonurutu/shade/core/camera"
	"github.com/aeonurutu/shade/core/display"
	"github.com/aeonurutu/shade/core/entity"
	"github.com/aeonurutu/shade/core/events"
	"github.com/aeonurutu/shade/core/shapes"
	"github.com/aeonurutu/shade/core/time/clock"
	"github.com/aeonurutu/shade/core/util/fonts"
	"github.com/aeonurutu/shade/core/util/sprite"

	"github.com/aeonurutu/shade/examples/03-basic-collisions/ball"
	"github.com/aeonurutu/shade/examples/03-basic-collisions/block"
	"github.com/aeonurutu/shade/examples/03-basic-collisions/player"
)

const windowWidth = 640
const windowHeight = 480

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	// Setup
	screen, err := display.SetMode("03-basic-collisions", windowWidth, windowHeight)
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

	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	objects := []entity.Entity{}

	blockSprite, err := loadSprite("assets/block32x32.png", "", 2, 1)
	if err != nil {
		panic(err)
	}
	blockSprite.Bind(screen.Program)
	objects = append(objects, block.New(0, float32(windowWidth)/6, float32(windowHeight)/2, blockSprite, *font))

	ballSprite, err := loadSprite("assets/ball.png", "", 1, 1)
	if err != nil {
		panic(err)
	}
	ballSprite.Bind(screen.Program)
	objects = append(objects, ball.New(float32(windowWidth)/2, float32(windowHeight)/2, ballSprite, *font))

	//shapes.NewCircle(mgl32.Vec2{float32(s.Width) / 2, float32(s.Height) / 2}, float32(s.Width)/2),
	tmpSprites := []sprite.Context{*blockSprite, *ballSprite}
	tmpShapes := []shapes.Shape{
		*shapes.NewRect(0, float32(blockSprite.Width), 0, float32(blockSprite.Height)),
		*shapes.NewCircle(mgl32.Vec2{float32(ballSprite.Width) / 2, float32(ballSprite.Height) / 2}, float32(ballSprite.Width)/2),
	}
	pl := player.New(0, 0, tmpSprites, tmpShapes, *font)
	objects = append(objects, &pl)

	for running := true; running; {
		dt := clock.Tick(30)

		screen.Fill(0.3, 0.3, 0.6)

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

			pl.Handle(event)
		}

		for _, e := range objects {
			if u, ok := e.(entity.Updater); ok {
				u.Update(dt, &objects)
			}
			if d, ok := e.(entity.Drawer); ok {
				d.Draw()
			}
		}

		screen.Flip()
		events.Poll()

	}
}

func loadSprite(colorName, normalName string, framesWide, framesHigh int) (*sprite.Context, error) {
	c, err := sprite.LoadAsset(colorName)
	if err != nil {
		return nil, err
	}

	n, err := sprite.LoadAsset(normalName)
	if err != nil {
		return nil, err
	}

	s, err := sprite.New(c, n, framesWide, framesHigh)
	if err != nil {
		return nil, err
	}

	return s, nil
}
