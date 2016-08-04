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

package main

import (
	"github.com/aeonurutu/shade"

	"github.com/aeonurutu/shade/examples/demos/platformer/app"
)

// Run with:
// go run -ldflags="-X github.com/aeonurutu/shade.ldDevBuild=true" main.go
// to enable dev mode.

func main() {
	a := app.New()

	e := shade.New("Platformer")
	e.SetFPS(59.94)
	e.SetEntryPoint(a)

	if err := e.Run(); err != nil {
		panic(err)
	}
}
