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

//go:generate go generate github.com/aeonurutu/shade

import (
	"runtime"

	"github.com/aeonurutu/shade"
	"github.com/aeonurutu/shade/examples/demos/pong/game"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	// Configure the engine
	eng := shade.New("Pong")

	// Configure your app
	scene := game.New()

	// Start the app
	if err := eng.Run(scene); err != nil {
		panic(err)
	}
}
