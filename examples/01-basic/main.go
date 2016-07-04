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

//go:generate go generate github.com/aeonurutu/shade

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/aeonurutu/shade"
	"github.com/aeonurutu/shade/core/camera"
	"github.com/aeonurutu/shade/core/display"
	"github.com/aeonurutu/shade/core/entity"
	"github.com/aeonurutu/shade/core/scene"
)

func main() {
	// Configure the engine
	eng := shade.New("01-basic")

	// Configure your app
	scene := &MyScene{}

	// Start the app
	if err := eng.Run(scene); err != nil {
		panic(err)
	}
}

// MyScene renders a single sprite.
type MyScene struct {
	camera.Camera2D // needed for default ViewMatrix() func
	display.Context // needed for default ProjMatrix() func
	scene.Single    // needed for default SubScenes() func

	image []io.Reader
	count int
}

// Setup MyScene
func (ctx *MyScene) Setup() error {
	fmt.Println("MyScene.Setup()")
	f, err := os.Open("test-pattern128x128.png")
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	ctx.image = append(ctx.image, r)
	return nil
}

// Entities for MyScene
func (ctx *MyScene) Entities() []entity.Entity {
	fmt.Println("MyScene.Entities()")
	fmt.Println(ctx.image)
	ctx.count++
	// TODO(hurricanerix): return sprite
	return nil
}

// ShouldStop MyScene when
func (ctx MyScene) ShouldStop() bool {
	// TODO(hurricanerix): This should really be triggered by user input, but
	//  will work for now.
	if ctx.count > 10 {
		return true
	}
	return false
}

// Cleanup MyScene
func (ctx *MyScene) Cleanup() {
	fmt.Println("MyScene.Cleanup()")
	// TODO(hurricanerix): unload sprite
}
