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

// Package ball manages a ball's state

package ball

import "time"
import "math/rand"
import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/sprite"
)

const TopY = 455
const BottomY = 75

// Ball state
type Ball struct {
	pos    mgl32.Vec3
	Sprite *sprite.Context
	velocity	int
	moveX	int
	moveY	int
}

func randXY() (int, int) {
	randSource := rand.NewSource(time.Now().UnixNano())
    randomNum := rand.New(randSource)

    var x, y int
    for x = randomNum.Intn(3); x == 0; {
    	x = randomNum.Intn(3)
    }
    
    for y = randomNum.Intn(3); y == 0; {
    	y = randomNum.Intn(3)
    }
	
	return x, y
}

func New(pos, dir mgl32.Vec3, s *sprite.Context) *Ball {
	randX, randY := randXY()
	
	b := Ball{
		pos:    pos,
		Sprite: s,
		velocity: 3,
		moveX: randX,
		moveY: randY,
	}
	return &b
}

func (b Ball) Pos() mgl32.Vec3 {
	return b.pos
}

func (b *Ball) Update(dt float32, group *[]entity.Entity) {
	// reverse y direction if ball contact top or bottom of screen
	if b.pos[1] > TopY {
		b.moveY *= -1
	}
	if b.pos[1] < BottomY {
		b.moveY *= -1
	}

	/*
	// get collision
	var cgroup []entity.Collider
	for i := range *group {
		if c, ok := (*group)[i].(entity.Collider); ok {
			cgroup = append(cgroup, c)
		}
	}

	for _, c := range entity.Collide(b, &cgroup, false) {
		p.Collision = &c
	}
	*/

	// reverse x direction if ball contact paddles
	//if collideLeft {
	if b.pos[0] > 500 {
		b.moveX *= -1
	}
	if b.pos[0] < 0 {
		b.moveX *= -1
	}

	b.pos[0] += float32(b.moveX * b.velocity)
	b.pos[1] += float32(b.moveY * b.velocity)
}

func (b Ball) Draw() {
	b.Sprite.DrawFrame(mgl32.Vec2{0, 0}, b.pos, nil)
}
