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

package shade

import (
	"testing"

	"github.com/aeonurutu/shade/entity"
	"github.com/aeonurutu/shade/scene"
	"github.com/go-gl/mathgl/mgl32"
)

// TestScene for unit tests.
type TestScene struct {
	projMatrix *mgl32.Mat4
	viewMatrix *mgl32.Mat4
	subs       []scene.Scene

	setupCalled    bool
	EntitiesCalled bool
}

// ... interface
func (ctx *TestScene) ProjMatrix() *mgl32.Mat4 {
	return ctx.projMatrix
}

// Camera interface
func (ctx *TestScene) ViewMatrix() *mgl32.Mat4 {
	return ctx.viewMatrix
}

func (ctx *TestScene) Setup() error {
	return nil
}

func (ctx *TestScene) Entities() []entity.Entity {
	return nil
}

func (ctx *TestScene) SubScenes() []scene.Scene {
	return ctx.subs
}

func (ctx TestScene) ShouldStop() bool {
	return true
}

func (ctx *TestScene) Cleanup() {
}

func TestEngine(t *testing.T) {
	eng := Engine{}

	s := TestScene{}

	// TODO(hurricanerix): Test Setup is called.
	// TODO(hurricanerix): Test Entieis() is called.
	// TODO(hurricanerix): Test Update() is called for each entity
	// TODO(hurricanerix): Test Draw() is called for each entity
	// TODO(hurricanerix): Test SubScenese is called
	// TODO(hurricanerix): Test each SubScene is processed.

	err := eng.Run(&s)
	if err != nil {
		t.Error("Foo")
	}

}
