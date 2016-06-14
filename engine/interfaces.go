package engine

import (
	"github.com/aeonurutu/shade/entity"
	"github.com/go-gl/mathgl/mgl32"
)

type Scene interface {
	Setup() error
	Cleanup()

	ViewMatrix() *mgl32.Mat4
	ProjMatrix() *mgl32.Mat4
	Entities() []entity.Entity

	SubScenes() []Scene

	ShouldStop() bool
}
