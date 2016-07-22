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
// Package shapes TODO doc

package shapes

import "github.com/go-gl/mathgl/mgl32"

type Shape struct {
	Type string
	Data []float32
}

// NewRect TODO doc
func NewRect(left, right, bottom, top float32) *Shape {
	r := Shape{
		Type: "rect",
		Data: []float32{left, right, bottom, top},
	}
	return &r
}

// NewCircle TODO doc
func NewCircle(center mgl32.Vec2, radius float32) *Shape {
	r := Shape{
		Type: "circle",
		Data: []float32{center[0], center[1], radius},
	}
	return &r
}
