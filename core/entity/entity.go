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

// package entity TODO: doc
package entity

type Entity interface {
	//ID() string
}

type EntityState struct {
	ID string `json:"id"`
}

func New() Entity {
	return &CoreEntity{}
}

type CoreEntity struct {
	id string
}

// ID is something.
func (e CoreEntity) ID() string {
	return e.id
}

func (e CoreEntity) State() EntityState {
	return EntityState{
	//ID
	}
}

func (e *CoreEntity) SetState(newState EntityState) {
}
