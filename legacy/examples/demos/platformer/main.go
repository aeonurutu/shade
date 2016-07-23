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
	"flag"
	"log"
	"runtime"

	"github.com/aeonurutu/shade/core/display"
	"github.com/aeonurutu/shade/core/splash"

	"github.com/aeonurutu/shade/examples/demos/platformer/game"
)

var (
	dev      bool
	nosplash bool
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	flag.BoolVar(&dev, "dev", false, "dev mode.")
	flag.BoolVar(&nosplash, "nosplash", false, "don't show splash screen.")
}

func main() {
	flag.Parse()

	screen, err := display.SetMode("ex2-platform", 640, 480)
	if err != nil {
		log.Fatalln("failed to set display mode:", err)
	}

	g, err := game.New(screen)
	if err != nil {
		log.Fatalln("failed to create game:", err)
	}

	if !nosplash {
		// Please see shade/splash/splash.go for details on
		// creating a splash screen
		splash.Main(screen)
	}

	config := game.Config{DevMode: dev}

	g.Main(screen, config)
}
