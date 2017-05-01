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

package engine

import (
	"flag"
	"strconv"
)

// CLI flags
var (
	fullscreen bool
	version    bool
)

// Dev mode CLI flags (dev builds only)
var ldDevBuild string
var (
	devBuild    bool
	fps         float64
	debugServer bool
)

func init() {
	flag.BoolVar(&fullscreen, "fullscreen", false, "Launch fullscreen")
	var err error // To keep err from shadowing devBuild
	devBuild, err = strconv.ParseBool(ldDevBuild)

	if err != nil {
		devBuild = false
	}

	if devBuild {
		// Parse dev build specific flags.
		flag.Float64Var(&fps, "fps", -1, "frames per second")
		flag.BoolVar(&debugServer, "debug-server", false, "Launch debug server")
	}
}
