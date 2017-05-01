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

package manager

import (
	"github.com/hurricanerix/shade/core/manager/video"
	"github.com/veandco/go-sdl2/sdl"
)

// Video TODO: doc
type Video interface {
	Init(width, height int) error
	WindowSize() (int, int)
	Swap()
	Quit()
}

// NewVideo TODO: doc
func NewVideo(width, height int) (Video, error) {
	if err := initSDL(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}
	mgr := video.Context{
		QuitSDL: quitSDL,
	}

	if err := mgr.Init(width, height); err != nil {
		return nil, err
	}

	return &mgr, nil
}
