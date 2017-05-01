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

// package project TODO: doc
package project

import (
	"fmt"

	"github.com/hurricanerix/shade/core/scene"
)

// Project interface
type Project interface {
	Name() string

	Scene(string) (scene.Scene, error)
	AddScene(scene.Scene) error
	RemoveScene(string)

	State() ProjectState
	SetState(s ProjectState)
}

// ProjectState structure for serialization.
type ProjectState struct {
	Name      string             `json:"name"`
	SceneList []scene.SceneState `json:"scene_list"`
}

// New Project is created and returned.
func New(name string) Project {
	return &CoreProject{
		name:      name,
		sceneList: make(map[string]scene.Scene),
	}
}

// CoreProject is a built in implementation of the Project interface.
type CoreProject struct {
	name      string
	sceneList map[string]scene.Scene
}

// Name of the project.
func (prj CoreProject) Name() string {
	return prj.name
}

func (prj CoreProject) Scene(name string) (scene.Scene, error) {
	if val, ok := prj.sceneList[name]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("scene with name \"%s\" does not exist", name)
}

func (prj *CoreProject) AddScene(scn scene.Scene) error {
	name := scn.Name()
	if _, ok := prj.sceneList[name]; ok {
		return fmt.Errorf("scene with name \"%s\" already exists", name)
	}
	prj.sceneList[name] = scn
	return nil
}

func (prj *CoreProject) RemoveScene(name string) {
	delete(prj.sceneList, name)
}

// State of the project internally.
func (prj CoreProject) State() ProjectState {
	ste := ProjectState{
		Name:      prj.Name(),
		SceneList: make([]scene.SceneState, len(prj.sceneList)),
	}
	i := 0
	for _, s := range prj.sceneList {
		ste.SceneList[i] = s.State()
	}
	return ste
}

// SetState of the project.
func (prj *CoreProject) SetState(ste ProjectState) {
	prj.name = ste.Name
}
