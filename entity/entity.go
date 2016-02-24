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
	"github.com/hurricanerix/shade/shapes"
)

// Entity ...
type Entity interface{}

// Updater ...
type Updater interface {
	Update(dt float32, group *[]Entity)
}

// Drawer ...
type Drawer interface {
	Pos() mgl32.Vec3
	Draw()
}

// Collider ...
type Collider interface {
	Pos() mgl32.Vec3
	Bounds() shapes.Shape
}

type Collision struct {
	Hit Collider
	Dir mgl32.Vec3
}

// Collide target with all enttities in group, returning all hits.  If cleanup is true
// hits are also removed from the group.
func Collide(target Collider, group *[]Collider, cleanup bool) (hits []Collision) {
	if target == nil || group == nil {
		return hits
	}

	tBounds := target.Bounds()
	tPos := target.Pos()
	///tType := reflect.TypeOf(target).String()
	tType := tBounds.Type

	var hit bool
	var dir mgl32.Vec3
	var ep mgl32.Vec3
	var eb shapes.Shape
	for i := range *group {
		if target == (*group)[i] {
			// Don't match if target is in group
			continue
		}
		hit = false
		ep = (*group)[i].Pos()
		eb = (*group)[i].Bounds()

		if tType == eb.Type && tType == "rect" {
			hit, dir = rectTest(tPos, ep, tBounds, eb, true)
			//} else if tType == reflect.TypeOf(eb).String() && tType == "circle" {
		} else if tType == eb.Type && tType == "circle" {
			hit, dir = circleTest(tPos, ep, tBounds, eb, true)
		} else {
			hit, dir = mixedTest(tPos, ep, tBounds, eb, true)
		}

		if hit {
			hits = append(hits, Collision{Hit: (*group)[i], Dir: dir})
		}
	}
	return hits
}

/**
 * points for Rectangles are probably lower left.
 * points for circles are probably center.
 *
 * @param ap point for shape a
 * @param bp point for shape b
 * @param ab shapeA
 * @param bb shapeB
 * @param ignoreZ if z should be ignored
 * @return a tuple of boolean (collision if true), Vec3 (direction of collision?)
 */
func rectTest(ap, bp mgl32.Vec3, ab, bb shapes.Shape, ignoreZ bool) (bool, mgl32.Vec3) {
	apt := mgl32.Vec3{
		ap[0] + ab.Data[0],
		ap[1] + ab.Data[2],
	}
	aw := float32(math.Abs(float64(ab.Data[0]) + float64(ab.Data[1])))
	ah := float32(math.Abs(float64(ab.Data[2]) + float64(ab.Data[3])))

	bpt := mgl32.Vec3{
		bp[0] + bb.Data[0],
		bp[1] + bb.Data[2],
	}

	bw := float32(math.Abs(float64(bb.Data[0]) + float64(bb.Data[1])))
	bh := float32(math.Abs(float64(bb.Data[2]) + float64(bb.Data[3])))

	if apt[0] < bpt[0]+bw &&
		apt[0]+aw > bpt[0] &&
		apt[1] < bpt[1]+bh &&
		ah+apt[1] > bpt[1] {
		return true, getDir(apt, bpt)
	}
	return false, mgl32.Vec3{}
}

func circleTest(ap, bp mgl32.Vec3, ab, bb shapes.Shape, ignoreZ bool) (bool, mgl32.Vec3) {
	if ap[0] == bp[0] && ap[1] == bp[1] && ab.Data[0] == bb.Data[0] {
		return false, mgl32.Vec3{}
	}

	apt := mgl32.Vec3{ap[0] + ab.Data[0], ap[1] + ab.Data[1]}
	bpt := mgl32.Vec3{bp[0] + bb.Data[0], bp[1] + bb.Data[1]}

	// distance between centers
	d := math.Sqrt(
		math.Pow(float64(apt[0])-float64(bpt[0]), 2) +
			math.Pow(float64(apt[1])-float64(bpt[1]), 2))
	// Sum of radiuses
	if d <= float64(ab.Data[2])+float64(bb.Data[2]) {
		return true, getDir(apt, bpt)
	}
	return false, mgl32.Vec3{}
}

func mixedTest(ap, bp mgl32.Vec3, ab, bb shapes.Shape, ignoreZ bool) (bool, mgl32.Vec3) {
	var rp mgl32.Vec3
	var rs shapes.Shape
	var cp mgl32.Vec3
	var cs shapes.Shape
	mod := float32(1.0)
	if ab.Type == "rect" {
		rp = ap
		rs = ab
		cp = mgl32.Vec3{bp[0] + bb.Data[0] - bb.Data[2], bp[1] + bb.Data[1] - bb.Data[2]}
		cs = bb
	} else {
		rp = bp
		rs = bb
		cp = mgl32.Vec3{ap[0] + ab.Data[0] - ab.Data[2], ap[1] + ab.Data[1] - ab.Data[2]}
		cs = ab
		mod = -1.0
	}
	w := float32(cs.Data[2]) * 2
	result, dir := rectTest(rp, cp, rs, *shapes.NewRect(0, w, 0, w), ignoreZ)
	dir[0] *= mod
	dir[1] *= mod
	dir[2] *= mod
	return result, dir
}

/**
 * Test obtaining collision vector direction.
 * See: http://www.vobarian.com/collisions/2dcollisions2.pdf
 *
 * Find unit normal and unit tangent vectors. The unit normal vector is a
 * vector which has a magnitude of 1 and a direction that is normal 
 * (perpendicular) to the surfaces of the objects at the point of collision.
 * The unit tangent vector is a vector with a magnitude of 1 which is
 * tangent to the circles' surfaces at the point of collision.
 *
 * First find a normal vector. This is done by taking a vector whose 
 * components are the difference between the coordinates of the centers 
 * of the circles. Let x1, x2, y1, and y2 be the x and y coordinates of 
 * the centers of the circles. (It does not matter which circle is 
 * labeled 1 or 2; the end result will be the same.) Then the normal 
 * vector n is:
 *
 * n→  =〈 x_2 − x_1, y_2 − y_1 〉
 *
 * Next, find the unit vector of n→, which we will call un→.
 * This is done by dividing by the magnitude of n→:
 * 
 * un→ =  n→ / (sqrt(n_x^2 + n_y^2))
 */
func getDir(aPos, bPos mgl32.Vec3) mgl32.Vec3 {
	xdiff := bPos[0] - aPos[0]
	ydiff := bPos[1] - aPos[1]
	//zdiff := bPos[2] - aPos[2]
	// Note:
	// The 3rd value, which is often magnitude, here is maybe or maybe not
	// a z value. If it IS a z value, one probably must include the zdiff^2
	// in the denominator sum, I believe.
	denominator := float32(math.Sqrt(float64((xdiff * xdiff) + (ydiff * ydiff))))
	if denominator == 0 {
		// Is this a reasonable answer? I realize divide by zero is bad
		// but what are the side effects of doing this?
		return mgl32.Vec3{}
	}
	uv := mgl32.Vec3{xdiff / denominator, 
		ydiff / denominator,
		1 /*zdiff / denominator*/}
	return uv
}
