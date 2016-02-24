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
// Package game manages the main game loop.

package game

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/camera"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/examples/ex1-pong/ball"
	"github.com/hurricanerix/shade/examples/ex1-pong/player"
	"github.com/hurricanerix/shade/fonts"
	"github.com/hurricanerix/shade/sprite"
	"github.com/hurricanerix/shade/time/clock"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Config TODO doc
type Config struct {
	DevMode bool
}

// Context TODO doc
type Context struct {
	Screen *display.Context
}

// New TODO doc
func New(screen *display.Context) (Context, error) {
	return Context{
		Screen: screen,
	}, nil
}

// Main TODO doc
func (c *Context) Main(screen *display.Context, config Config) {
	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Bind(c.Screen.Program)

	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	paddleSprite, err := loadSpriteAsset("assets/paddle.png", "", 1, 3)
	if err != nil {
		panic(err)
	}
	paddleSprite.Bind(screen.Program)

	ballSprite, err := loadSpriteAsset("assets/ball.png", "", 1, 1)
	if err != nil {
		panic(err)
	}
	ballSprite.Bind(screen.Program)

	var objects []entity.Entity
	player1 := player.New(1, cam.Left+15, screen.Height/4, paddleSprite)
	objects = append(objects, player1)
	player2 := player.New(2, cam.Right-15, screen.Height/4, paddleSprite)
	objects = append(objects, player2)
	ball := ball.New(mgl32.Vec3{screen.Width / 2, screen.Height / 2, 0.0}, mgl32.Vec3{0, 1, 0}, ballSprite)
	objects = append(objects, ball)

	font, err := fonts.SimpleASCII()
	if err != nil {
		panic(err)
	}
	font.Bind(screen.Program)

	// font color
	efxFont := sprite.Effects{
		EnableLighting: false,
		Scale:          mgl32.Vec3{2.0, 2.0, 1.0},
		Tint:           mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
	}

	for running := true; running; {

		screen.Fill(0, 0, 0)

		dt := clock.Tick(30)

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
			player1.Handle(event)
			player2.Handle(event)
		}

		for _, e := range objects {
			if u, ok := e.(entity.Updater); ok {
				u.Update(dt, &objects)
			}
			if d, ok := e.(entity.Drawer); ok {
				d.Draw()
			}
		}

		if config.DevMode {
			msg := "Dev Mode!\n"
			msg += fmt.Sprintf("Player1: %v\n", player1.Pos())
			msg += fmt.Sprintf("Player2: %v\n", player2.Pos())
			msg += fmt.Sprintf("Ball: %v\n", ball.Pos())
			font.DrawText(mgl32.Vec3{cam.Left + 20, cam.Top - 40, 0}, &efxFont, msg)
		}
		msg := "Score: 0"
		font.DrawText(mgl32.Vec3{cam.Left, 0, 0}, &efxFont, msg)
		w, _ := font.SizeText(&efxFont, msg)
		font.DrawText(mgl32.Vec3{cam.Right - w, 0, 0}, &efxFont, msg)
		screen.Flip()
		events.Poll()
	}
}

func loadSpriteAsset(colorName, normalName string, framesWide, framesHigh int) (*sprite.Context, error) {
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
