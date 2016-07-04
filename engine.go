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

/*
Package shade is a 2.5D game engine.

A simple invocation of the engine can be done as follows:

  //go:generate go generate github.com/aeonurutu/shade

  package main

  import (
    "github.com/aeonurutu/shade"
    "mygame"
  )

  func main() {
    // Configure the engine
    eng := shade.New("AppName")

    // Configure your app
    scene := &mygame.MyScene{}

    // Start the app
    if err := eng.Run(scene); err != nil {
      panic(err)
    }
  }

To run the app in dev mode, use the following:

	$ go run -ldflags="-X github.com/aeonurutu/shade.allowDevMode=true" main.go -dev

When building your app, make sure to run go generate first to ensure that
Shade's generated files are up to date.

	$ go generate
	$ go build -ldflags="-X github.com/aeonurutu/shade.allowDevMode=true" main.go -o appname-dev
	$ # or
	$ go build -o appname main.go

*/
package shade

//go:generate ./scripts/gen.sh

import (
	"flag"
	"log"
	"runtime"
	"strconv"

	"github.com/aeonurutu/shade/core/dev"
	"github.com/aeonurutu/shade/core/display"
	"github.com/aeonurutu/shade/core/entity"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	allowDevMode string
	devFlag      bool
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	// allowDevMode should be set with ldflags
	adm, _ := strconv.ParseBool(allowDevMode)
	if adm {
		flag.BoolVar(&devFlag, "dev", false, "dev mode.")
	}
}

// Scene represents a View along with a collection of Entities and or SubScenes.
type Scene interface {
	Setup() error
	Cleanup()

	ViewMatrix() *mgl32.Mat4
	ProjMatrix() *mgl32.Mat4
	Entities() []entity.Entity

	SubScenes() []Scene

	ShouldStop() bool
}

// Engine contains state for the core systems of a 2.5D game.
type Engine struct {
	Title        string
	AllowDevMode bool
}

// New returns a pointer to an Engine.
func New(title string) *Engine {
	e := Engine{
		Title: title,
	}
	return &e
}

// Run a Scene with an existing Engine.
func (e *Engine) Run(scene Scene) error {
	var err error
	flag.Parse()

	screen, err := display.SetMode(e.Title, 512, 512)
	if err != nil {
		log.Fatalln("failed to set display mode:", err)
	}
	println(screen)

	var devDisplay *dev.Context
	if devFlag {
		devDisplay = dev.New()
	}

	running := true

	err = scene.Setup()
	if err != nil {
		return err
	}
	defer scene.Cleanup()

	for _, sub := range scene.SubScenes() {
		err = sub.Setup()
		defer sub.Cleanup()
		if err != nil {
			return err
		}
	}

	for running {
		for _, ent := range scene.Entities() {
			println(ent)
		}

		if devDisplay != nil {
			devDisplay.Update()
			devDisplay.Draw()
		}

		if scene.ShouldStop() {
			running = false
		}
	}

	return nil
}
