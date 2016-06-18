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
*/
package shade

import (
	"github.com/aeonurutu/shade/entity"
	"github.com/go-gl/mathgl/mgl32"
)

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
	Title string
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

		if scene.ShouldStop() {
			running = false
		}
	}

	return nil
}
