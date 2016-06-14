package scene

import "github.com/aeonurutu/shade/engine"

type Single struct {
}

func (ctx *Single) SubScenes() []engine.Scene {
	return nil
}
