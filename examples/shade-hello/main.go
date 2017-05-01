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

// shade-hello is a bare bones "hello, world" style app.
// The app will load a single image and render it in the center of the window.
//
// Run the example with the following command:
//
// 	go run \
//	-ldflags="-X github.com/hurricanerix/shade/core/engine.ldDevBuild=true" \
//	$GOPATH/src/github.com/hurricanerix/shade/examples/shade-hello/main.go \
//	--debug-server
//
// Setting ldDevBuild in the ldflgas builds the app as a dev build.  Dev builds
// provide additional CLI flags which are not avaialble otherwise.
//
// One such flag is the --deubg-server, providing this flag launchs a HTTP server
// while the app is running.  Going to http://127.0.0.1:10888/ in a web browser
// will load Shade Engine's debug interface.
package main

import (
	"go/build"
	"log"
	"os"

	"github.com/hurricanerix/shade/core/engine"
	"github.com/hurricanerix/shade/core/project"
	"github.com/hurricanerix/shade/core/scene"
	"github.com/hurricanerix/shade/core/sprite"
)

func main() {
	// All applications contain a single project.
	prj := project.New("Hello, World")

	// All applications contain one or more scenes, for this example, we
	// will only create a single scene named "main".
	scn := scene.New("main")
	if err := prj.AddScene(scn); err != nil {
		panic(err)
	}

	// We need something to render in our scene, so a sprite is created
	// and added to the scene.
	hello := sprite.New("hello.png")
	// TODO: center the sprite in the window.
	scn.AddEntity(hello)

	// Last we run our project using the engine.
	if err := engine.Run(prj); err != nil {
		panic(err)
	}
}

// The code below was borrowed from the go-gl examples
// https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	dir, err := importPathToDir("github.com/hurricanerix/shade/examples/shade-hello")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
