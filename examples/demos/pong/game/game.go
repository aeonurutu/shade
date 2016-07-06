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
// Package game manages the main game loop.

package game

import (
	"bytes"
	"fmt"
	"image"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/aeonurutu/shade/core/camera"
	"github.com/aeonurutu/shade/core/display"
	"github.com/aeonurutu/shade/core/entity"
	"github.com/aeonurutu/shade/core/events"
	"github.com/aeonurutu/shade/core/scene"
	"github.com/aeonurutu/shade/core/time/clock"
	"github.com/aeonurutu/shade/core/util/archive"
	"github.com/aeonurutu/shade/core/util/fonts"
	"github.com/aeonurutu/shade/core/util/sprite"

	"github.com/aeonurutu/shade/examples/demos/pong/ball"
	"github.com/aeonurutu/shade/examples/demos/pong/player"
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
	Screen          *display.Context
	camera.Camera2D // needed for default ViewMatrix() func
	display.Context // needed for default ProjMatrix() func
	scene.Single    // needed for default SubScenes() func
}

// New TODO doc
func New() *Context {
	ctx := Context{
	//Screen: screen,
	}
	return &ctx
}

// Setup Game
func (ctx *Context) Setup() error {
	fmt.Println("MyScene.Setup()")
	return nil
}

// Entities for Game
func (ctx *Context) Entities() []entity.Entity {
	return nil
}

// ShouldStop Game when
func (ctx Context) ShouldStop() bool {
	return false
}

// Cleanup Game
func (ctx *Context) Cleanup() {
	fmt.Println("MyScene.Cleanup()")
}

// Main TODO doc
func (c *Context) Main() {
	var screen *display.Context
	var config *Config

	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Bind(c.Screen.Program)

	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	// font object for text display on screen
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

	// load paddle and ball sprites
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

	// setup players and ball objects
	var objects []entity.Entity
	player1 := player.New(1, cam.Left+15, screen.Height/4, paddleSprite)
	objects = append(objects, player1)
	player2 := player.New(2, cam.Right-15, screen.Height/4, paddleSprite)
	objects = append(objects, player2)
	ball := ball.New(mgl32.Vec3{screen.Width / 2, screen.Height / 2, 0.0}, mgl32.Vec3{0, 1, 0}, ballSprite, player1, player2)
	objects = append(objects, ball)

	// has winner flag
	var hasWinner bool

	// game loop
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

		// draw scores
		msgScore1 := fmt.Sprintf("Score: %d", player1.Score)
		font.DrawText(mgl32.Vec3{cam.Left, cam.Top - 15, 0}, &efxFont, msgScore1)

		msgScore2 := fmt.Sprintf("Score: %d", player2.Score)
		w, _ := font.SizeText(&efxFont, msgScore2)
		font.DrawText(mgl32.Vec3{cam.Right - w, cam.Top - 15, 0}, &efxFont, msgScore2)

		if config.DevMode {
			msg := "Dev Mode!\n"
			msg += fmt.Sprintf("Player1: %v\n", player1.Pos())
			msg += fmt.Sprintf("Player2: %v\n", player2.Pos())
			msg += fmt.Sprintf("Ball: %v\n", ball.Pos())
			font.DrawText(mgl32.Vec3{cam.Left + 20, cam.Top - 40, 0}, &efxFont, msg)
		}

		// print winner
		if player1.Score >= player.NumToWin {
			font.DrawText(mgl32.Vec3{cam.Left + 150, cam.Top - 50, 0}, &efxFont, "Player 1 is the Winner!!!!")
			hasWinner = true
		} else if player2.Score >= player.NumToWin {
			font.DrawText(mgl32.Vec3{cam.Left + 150, cam.Top - 50, 0}, &efxFont, "Player 2 is the Winner!!!!")
			hasWinner = true
		}

		// screen Flip and events polling
		screen.Flip()
		events.Poll()

		if hasWinner {
			// wait and reset score
			duration := time.Duration(5) * time.Second
			time.Sleep(duration)
			player1.Score = 0
			player2.Score = 0
			hasWinner = false
		}
	}
}

func loadSpriteAsset(colorName, normalName string, framesWide, framesHigh int) (*sprite.Context, error) {
	c, err := archive.Get("example_assets.tar", colorName)
	if err != nil {
		return nil, err
	}
	ic, _, err := image.Decode(bytes.NewReader(c))

	n, err := archive.Get("example_assets.tar", normalName)
	if err != nil {
		return nil, err
	}
	in, _, err := image.Decode(bytes.NewReader(n))

	s, err := sprite.New(ic, in, framesWide, framesHigh)
	if err != nil {
		return nil, err
	}

	return s, nil
}