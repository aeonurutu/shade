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
	"bufio"
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/aeonurutu/shade/core/camera"
	"github.com/aeonurutu/shade/core/display"
	"github.com/aeonurutu/shade/core/entity"
	"github.com/aeonurutu/shade/core/events"
	"github.com/aeonurutu/shade/core/time/clock"
	"github.com/aeonurutu/shade/core/util/archive"
	"github.com/aeonurutu/shade/core/util/fonts"
	"github.com/aeonurutu/shade/core/util/sprite"

	"github.com/aeonurutu/shade/examples/demos/platformer/block"
	"github.com/aeonurutu/shade/examples/demos/platformer/player"
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

type Scene struct {
	Sprites []sprite.Sprite
	Player  *player.Player
	Objects []entity.Entity
	//Walls   []entity.Collider
}

// Main TODO doc
func (c *Context) Main(screen *display.Context, config Config) {
	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Offset = mgl32.Vec2{200, 100}
	cam.TopStop = 64 * 6.5 // TODO: should be 64 x 14
	cam.RightStop = 64 * 54
	cam.LeftStop = 1
	cam.BottomStop = 1
	cam.Bind(c.Screen.Program)

	scene, err := loadMap("map.data")
	if err != nil {
		panic(err)
	}
	cam.Move(scene.Player.Pos())

	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	for _, s := range scene.Sprites {
		s.Bind(screen.Program)
	}

	font, err := fonts.SimpleASCII()
	if err != nil {
		panic(err)
	}
	font.Bind(screen.Program)

	for running := true; running; {

		screen.Fill(1, 0, 0)

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
			for _, e := range scene.Objects {
				if h, ok := e.(events.Handler); ok {
					h.Handle(event)
				}
			}
		}

		//cam.Move(scene.Player.Pos())
		cam.Follow(scene.Player.Pos(), 0.1)

		for _, e := range scene.Objects {
			if u, ok := e.(entity.Updater); ok {
				u.Update(dt/1000, &scene.Objects)
			}
			if d, ok := e.(entity.Drawer); ok {
				d.Draw()
			}
		}

		//scene.Player.Update(dt/1000.0, scene.Walls)

		/*
			effect := sprite.Effects{
				Scale:          mgl32.Vec3{1.0, 1.0, 1.0},
				EnableLighting: true,
				AmbientColor:   mgl32.Vec4{0.3, 0.3, 0.3, 1.0},
				Light:          *scene.Player.Light}
		*/

		/*
			for _, w := range scene.Walls {
				e, ok := w.(entity.Entity)
				if ok {
					println("OKDRAW")
					e.Draw()
				}
			}
			scene.Player.Draw()
		*/

		if config.DevMode {
			deveff := sprite.Effects{
				EnableLighting: false,
				Scale:          mgl32.Vec3{2.0, 2.0, 1.0},
				Tint:           mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
			}
			msg := "Dev Mode!\n"
			msg += fmt.Sprintf("Camera Pos: %.0f, %.0f\n", cam.Pos[0], cam.Pos[1])
			msg += fmt.Sprintf("Player {\n")
			msg += fmt.Sprintf("  Pos: %v\n", scene.Player.Pos())
			msg += fmt.Sprintf("  Facing: %.0f\n", scene.Player.Facing)
			msg += fmt.Sprintf("  Resting: %t\n", scene.Player.Resting)
			msg += fmt.Sprintf("  Walking: %t\n", scene.Player.Walking)
			msg += fmt.Sprintf("  Light: {\n")
			msg += fmt.Sprintf("    Pos: %.0f, %.0f\n", scene.Player.Light.Pos[0], scene.Player.Light.Pos[1])
			msg += fmt.Sprintf("  }\n")
			msg += fmt.Sprintf("}\n")
			font.DrawText(mgl32.Vec3{cam.Left + 20, cam.Top - 40, 0}, &deveff, msg)
		}
		screen.Flip()
		events.Poll()

	}
}

// sprites, player, collidable
func loadMap(path string) (*Scene, error) {
	scene := Scene{}

	playerSprite, err := loadSpriteAsset("./gopher.png", "./gopher.normal.png", 5, 4)
	if err != nil {
		return &scene, err
	}
	scene.Sprites = append(scene.Sprites, playerSprite)
	blockSprite, err := loadSpriteAsset("./blocks.png", "./blocks.normal.png", 8, 1)
	if err != nil {
		return &scene, err
	}
	scene.Sprites = append(scene.Sprites, blockSprite)

	f, err := os.Open(path)
	if err != nil {
		return &scene, err
	}

	s := bufio.NewScanner(f)
	count := 0
	lines := []string{}
	for s.Scan() {
		count += 1
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		return &scene, err
	}

	x := float32(0)
	y := float32(0)
	for i := count - 1; i >= 0; i -= 1 {
		for _, c := range lines[i] {
			switch c {
			case '#':
				scene.Objects = append(scene.Objects, block.New(float32(x), float32(y), 0, blockSprite))
			case '-':
				scene.Objects = append(scene.Objects, block.New(float32(x), float32(y), 1, blockSprite))
			case 'D':
				scene.Objects = append(scene.Objects, block.New(float32(x), float32(y), 2, blockSprite))
			case 'S':
				scene.Player = player.New(x, y, playerSprite)
				scene.Objects = append(scene.Objects, scene.Player)
			}
			x += float32(blockSprite.Width)
		}
		x = 0
		y += float32(blockSprite.Height)
	}

	return &scene, nil
}

func loadSpriteAsset(colorName, normalName string, framesWide, framesHigh int) (*sprite.Context, error) {

	c, err := archive.Get("example_assets.tar", colorName)
	if err != nil {
		return nil, err
	}
	ic, _, err := image.Decode(bytes.NewReader(c))

	n, err := archive.Get("example_assets.tar", colorName)
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
