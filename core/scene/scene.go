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

// package scene TODO: doc
package scene

import "github.com/hurricanerix/shade/core/entity"

//import "github.com/hurricanerix/shade/core/entity"

type Scene interface {
	Name() string
	AddEntity(e entity.Entity)
	State() SceneState
	SetState(SceneState)
}

type SceneState struct {
	Name string `json:"name"`
}

func New(name string) Scene {
	return &CoreScene{
		name: name,
	}
}

type CoreScene struct {
	name string
}

func (scn CoreScene) Name() string {
	return scn.name
}

func (scn *CoreScene) AddEntity(e entity.Entity) {

}

func (scn CoreScene) State() SceneState {
	return SceneState{
		Name: scn.name,
	}
}

func (scn *CoreScene) SetState(ste SceneState) {
	scn.name = ste.Name
}
