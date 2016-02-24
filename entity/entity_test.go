// Copyright 2016 Richard Hawkins
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

// Package entity provies interfaces for game objects.
package entity

import (
	"math"
	"github.com/go-gl/mathgl/mgl32"
	"testing"
	"fmt"
)

func TestGetDir(t *testing.T) {
	data := make(map[string][]mgl32.Vec3)
	// TODO: Double check the expected values are what you expcted. Especially
	// TODO: on the mixed cases.
	t0 := "a.x < b.x, a.y < b.y"
	data[t0] = append(data[t0], mgl32.Vec3{ 0.0,  0.0, 0.0}) // a
	data[t0] = append(data[t0], mgl32.Vec3{ 5.0,  5.0, 0.0}) // b
	data[t0] = append(data[t0], mgl32.Vec3{-0.7, -0.7, 1.0}) // expected

	t1 := "a.x > b.x, a.y > b.y"
	data[t1] = append(data[t1], mgl32.Vec3{ 5.0,  5.0, 0.0}) // a
	data[t1] = append(data[t1], mgl32.Vec3{ 0.0,  0.0, 0.0}) // b
	data[t1] = append(data[t1], mgl32.Vec3{ 0.7,  0.7, 1.0}) // expected

	t2 := "a.x = b.x, a.y < b.y"
	data[t2] = append(data[t2], mgl32.Vec3{ 0.0,  0.0, 0.0}) // a
	data[t2] = append(data[t2], mgl32.Vec3{ 0.0,  5.0, 0.0}) // b
	data[t2] = append(data[t2], mgl32.Vec3{-0.7, -0.7, 1.0}) // expected

	t3 := "a.x = b.x, a.y > b.y"
	data[t3] = append(data[t3], mgl32.Vec3{ 0.0,  5.0, 0.0}) // a
	data[t3] = append(data[t3], mgl32.Vec3{ 0.0,  0.0, 0.0}) // b
	data[t3] = append(data[t3], mgl32.Vec3{-0.7, -0.7, 1.0}) // expected

	t4 := "a.x = b.x, a.y = b.y"
	data[t4] = append(data[t4], mgl32.Vec3{ 5.0,  5.0, 0.0}) // a
	data[t4] = append(data[t4], mgl32.Vec3{ 5.0,  5.0, 0.0}) // b
	data[t4] = append(data[t4], mgl32.Vec3{ 0.0,  0.0, 1.0}) // expected

	t5 := "a.x < b.x, a.y = b.y"
	data[t5] = append(data[t5], mgl32.Vec3{ 0.0,  5.0, 0.0}) // a
	data[t5] = append(data[t5], mgl32.Vec3{ 5.0,  5.0, 0.0}) // b
	data[t5] = append(data[t5], mgl32.Vec3{-0.7, -0.7, 1.0}) // expected

	t6 := "a.x > b.x, a.y = b.y"
	data[t6] = append(data[t6], mgl32.Vec3{ 5.0,  5.0, 0.0}) // a
	data[t6] = append(data[t6], mgl32.Vec3{ 0.0,  0.0, 0.0}) // b
	data[t6] = append(data[t6], mgl32.Vec3{ 0.7,  0.7, 1.0}) // expected

	// a lower and right of b
	for k := range data {
		aPos := data[k][0]
		bPos := data[k][1]
		expectedDir := data[k][2]
		dir := getDir(aPos, bPos)
		if (!(AboutTheSame(expectedDir[0], dir[0]) ||
				AboutTheSame(expectedDir[1], dir[1]) ||
				AboutTheSame(expectedDir[2], dir[2]))) {
			fmt.Println("Case", k, ":",
				"Expected getDir(a,b) results to be",
				expectedDir[0], expectedDir[1], expectedDir[2],
				"but found",
				dir[0], dir[1], dir[2])
			t.Error(
				"Expected getDir(a,b) results to be",
				expectedDir[0], expectedDir[1], expectedDir[2],
				"but found",
				dir[0], dir[1], dir[2])
		}
	}
}

//func TestRectTest(t *testing.T) {
//}

func AboutTheSame(a32, b32 float32) bool {
	tolerance := float64(0.01)
	a := float64(a32)
	b := float64(b32)
	if ((a < 0) && (b >= 0)) {
		return false
	}
	if (math.Abs(b - a) <= tolerance) {
		return true
	}
	return false
}
