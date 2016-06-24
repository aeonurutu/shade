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
	"math/rand"
	"runtime"
	"time"

	"github.com/aeonurutu/shade/camera"
	"github.com/aeonurutu/shade/display"
	"github.com/aeonurutu/shade/entity"
	"github.com/aeonurutu/shade/events"
	"github.com/aeonurutu/shade/examples/03-collisions/ball"
	"github.com/aeonurutu/shade/examples/03-collisions/block"
	"github.com/aeonurutu/shade/sprite"
	"github.com/aeonurutu/shade/time/clock"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const windowWidth = 640
const windowHeight = 480

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	screen, err := display.SetMode("03-collisions", windowWidth, windowHeight)
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

	objects := []entity.Entity{}

	blockSprite, err := loadSprite("assets/block32x32.png", "", 2, 1)
	if err != nil {
		panic(err)
	}
	blockSprite.Bind(screen.Program)

	for x := 0; float32(x) < screen.Width; x += 32 {
		for y := 0; float32(y) < screen.Height; y += 32 {
			if x == 0 || x == 640-32 || y == 0 || y == 480-32 {
				objects = append(objects, block.New(float32(x), float32(y), blockSprite))
			}
		}
	}
	objects = append(objects, block.New(float32(blockSprite.Width)*4, float32(blockSprite.Height)*4, blockSprite))
	objects = append(objects, block.New(float32(blockSprite.Width)*4, windowHeight-float32(blockSprite.Height)*5, blockSprite))
	objects = append(objects, block.New(windowWidth-float32(blockSprite.Width)*5, float32(blockSprite.Height)*4, blockSprite))
	objects = append(objects, block.New(windowWidth-float32(blockSprite.Width)*5, windowHeight-float32(blockSprite.Height)*5, blockSprite))

	ballSprite, err := loadSprite("assets/ball.png", "", 1, 1)
	if err != nil {
		panic(err)
	}
	ballSprite.Bind(screen.Program)

	rand.Seed(time.Now().Unix())
	//rand.Seed(1)

	objects = append(objects, addBall(screen.Width/2, screen.Height/2, ballSprite))

	//	sprites.Bind(screen.Program)
	for running := true; running; {
		dt := clock.Tick(30)

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
			if (event.Type == events.KeyDown || event.Type == events.KeyRepeat) && event.Key == glfw.KeySpace {
				objects = append(objects, addBall(screen.Width/2, screen.Height/2, ballSprite))
			}
		}

		for _, e := range objects {
			if u, ok := e.(entity.Updater); ok {
				u.Update(dt/1000, &objects)
			}
			if d, ok := e.(entity.Drawer); ok {
				d.Draw()
			}
		}

		screen.Flip()
		events.Poll()
	}
}

func addBall(x, y float32, s *sprite.Context) *ball.Ball {
	speed := float32(rand.Intn(500) + 200)
	angle := float32(rand.Intn(360))
	ball := ball.New(x, y, speed, angle, s)
	return &ball
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
