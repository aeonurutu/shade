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

import "fmt"
import "time"
import "math/rand"
import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/sprite"
	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/examples/ex1-pong/player"
)

// Ball state
type Ball struct {
	startPos	mgl32.Vec3
	pos    mgl32.Vec3
	Sprite *sprite.Context
	Shape    *shapes.Shape
	velocity	int
	moveX	int
	moveY	int
	colliding bool
	player1	*player.Player
	player2	*player.Player
}

func New(pos, dir mgl32.Vec3, s *sprite.Context, p1 *player.Player, p2 *player.Player) *Ball {
	randX, randY := randXY()
	
	b := Ball{
		startPos: pos,
		pos: pos,
		Sprite: s,
		Shape: shapes.NewCircle(mgl32.Vec2{float32(s.Width / 2), float32(s.Height / 2)}, float32(s.Width / 2)),
		velocity: 3,
		moveX: randX,
		moveY: randY,
		player1: p1,
		player2: p2,
	}
	fmt.Println("Ball created.")
	return &b
}

func randXY() (int, int) {
	randSource := rand.NewSource(time.Now().UnixNano())
    randomNum := rand.New(randSource)

    var x, y int
    for x = randomNum.Intn(3); x == 0; {
    	x = randomNum.Intn(3)
    }
    if x % 2 != 0 { x = 1 } else { x = -1 }

    for y = randomNum.Intn(3); y == 0; {
    	y = randomNum.Intn(3)
    }
    if y % 2 != 0 { y = 1 } else { y = -1 }
    
	return x, y
}

func (b Ball) Pos() mgl32.Vec3 {
	return b.pos
}

// Bind TODO doc
func (b *Ball) Bind(program uint32) error {
	return b.Sprite.Bind(program)
}

func (b Ball) Bounds() shapes.Shape {
	return *b.Shape
}

func (b *Ball) Update(dt float32, group *[]entity.Entity) {
	// reverse y direction if ball contact top or bottom of screen
	if b.pos[1] >= player.TopY - float32(b.Sprite.Height / 2) {
		b.moveY *= -1
	}
	if b.pos[1] <= player.BottomY {
		b.moveY *= -1
	}

	// reverse x direction if ball contact paddles
	var collided bool
	var cgroup []entity.Collider
	for i := range *group {
		if c, ok := (*group)[i].(entity.Collider); ok {
			cgroup = append(cgroup, c)
		}
	}
	for _, c := range entity.Collide(b, &cgroup, false) {
		if c.Hit.(*player.Player) != nil {
			collided = true
		}
	}
	
	// if collided, reverse direction once until collision is ended
	if !b.colliding && collided {
		//fmt.Println("BOOM!!!")
		b.moveX *= -1
		b.colliding = true
	} else if !collided {
		// end colliding once no collision is detected
		b.colliding = false
	}

	b.pos[0] += float32(b.moveX * b.velocity)
	b.pos[1] += float32(b.moveY * b.velocity)	

	var resetPos bool
	if b.pos[0] < -75 {
		//fmt.Println("Player 1 lose!")
		b.player2.Score += 1
		resetPos = true
	} else if b.pos[0] > 650 {
		//fmt.Println("Player 2 lose!")
		b.player1.Score += 1
		resetPos = true
	}

	if resetPos {
		b.pos[0] = b.startPos[0]
		b.pos[1] = b.startPos[1]
		randX, randY := randXY()
		b.moveX = randX
		b.moveY = randY
	}
}

func (b Ball) Draw() {
	b.Sprite.DrawFrame(mgl32.Vec2{0, 0}, b.pos, nil)
}
