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

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestCircle1(t *testing.T) {
	var circleCenterX float32 = 10
	var circleCenterY float32 = 15
	var circleRadius float32 = 20
	center := mgl32.Vec2{circleCenterX, circleCenterY}
	circle := NewCircle(center, circleRadius)
	if circle.Type != "circle" {
		t.Error(
			"Creating a circle did not return an object with Type 'circle'",
			"Type was", circle.Type)
	}
	if len(circle.Data) != 3 {
		t.Error(
			"Expected len(circle.Data) to be", 3,
			"But it was", len(circle.Data))
	}
	if circle.Data[0] != circleCenterX ||
		circle.Data[1] != circleCenterY ||
		circle.Data[2] != circleRadius {
		t.Error(
			"Expected Data to be", circleCenterX-1, circleCenterY, circleRadius,
			"but found",
			"Type was", circle.Data[0], circle.Data[1], circle.Data[2])
	}
}

func TestRect(t *testing.T) {
	rect := NewRect(5, 10, 15, 20)
	if rect.Type != "rect" {
		t.Error(
			"Creating a rectangle did not return an object with Type 'rect'",
			"Type is", rect.Type)
	}
	if rect.Data[0] != 5 || rect.Data[1] != 10 ||
		rect.Data[2] != 15 || rect.Data[3] != 20 {
		t.Error(
			"Expected Data to be", 5, 10, 15, 20,
			"but found",
			"Type was", rect.Data[0], rect.Data[1], rect.Data[2], rect.Data[3])
	}
}
