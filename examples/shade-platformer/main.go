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

// package main TODO: doc
package main

import (
	"github.com/hurricanerix/shade"
	"github.com/hurricanerix/shade/core/project"
	"github.com/hurricanerix/shade/core/scene"
)

// Run with:
// go run -ldflags="-X github.com/hurricanerix/shade.ldDevBuild=true" main.go --debug-server
// to enable dev mode.

func main() {

	prj := project.New("Platformer")

	scn := scene.New("main")
	if err := prj.AddScene(scn); err != nil {
		panic(err)
	}

	eng, err := shade.AttachEngine(prj)
	if err != nil {
		panic(err)
	}

	if err := eng.Run(); err != nil {
		panic(err)
	}
}
