// Copyright 2016-2017 Richard Hawkins
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

package manager

import (
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

var needsInit bool = true
var sdlCount int = 0

// init TODO: doc
func init() {
	runtime.LockOSThread()
}

// initSDL TODO: doc
func initSDL(flags uint32) error {
	// TODO: track quits and call defer sdl.Quit() when no more need for SDL
	var err error

	if needsInit {
		err = sdl.Init(flags)
		if err == nil {
			needsInit = false
			sdlCount++
		}
		return err
	}
	err = sdl.InitSubSystem(flags)
	if err == nil {
		sdlCount++
	}
	return err
}

// quitSDL decrement reference count, if there are no remaining
// SDL based managers left, quit SDL.
func quitSDL() {
	sdlCount--
	if sdlCount <= 0 {
		sdl.Quit()
	}
}
