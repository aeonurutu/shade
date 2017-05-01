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

// package shader TODO: doc
package shader

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Info representing a shader.
type Info struct {
	// Type of shader: gl.VERTEX_SHADER, gl.FRAGMENT_SHADER, ...
	Type uint32
	// Filename of shader source file.
	Filename string
	// shader ID.
	shader uint32
}

// Load the shaders, returning the ID of the resulting program.  Any problems
// compiling or linking will result in an error.
func Load(shaders *[]Info) (uint32, error) {
	return load(shaders, false)
}

// LoadSeparable is the same as Load with the exception that before the link stage
// GL_PROGRAM_SEPARABLE is set to GL_TRUE.
func LoadSeparable(shaders *[]Info) (uint32, error) {
	return load(shaders, true)
}

// load the shaders
func load(shaders *[]Info, separable bool) (uint32, error) {
	program := gl.CreateProgram()

	for _, s := range *shaders {
		if err := s.Compile(program); err != nil {
			cleanup(shaders)
			gl.DeleteProgram(program)
			return 0, err
		}
	}

	if separable {
		gl.ProgramParameteri(program, gl.PROGRAM_SEPARABLE, gl.TRUE)
	}

	gl.LinkProgram(program)
	cleanup(shaders)
	var linked int32
	if gl.GetProgramiv(program, gl.LINK_STATUS, &linked); linked == gl.FALSE {
		msg := getErrorMsg(false, program)
		gl.DeleteProgram(program)
		return 0, fmt.Errorf("failed to link program: %s", msg)
	}

	return program, nil
}

// Compile the shader using the info provided.
func (i *Info) Compile(program uint32) error {
	i.shader = gl.CreateShader(i.Type)
	if i.shader == 0 {
		return fmt.Errorf("could not create shader")
	}

	source, err := readShader(i.Filename)
	if err != nil {
		return err
	}

	csrc, free := gl.Strs(source)
	gl.ShaderSource(i.shader, 1, csrc, nil)
	free()
	gl.CompileShader(i.shader)

	var compiled int32
	if gl.GetShaderiv(i.shader, gl.COMPILE_STATUS, &compiled); compiled == gl.FALSE {
		return fmt.Errorf("failed to compile %s: %s", i.Filename, getErrorMsg(true, i.shader))
	}
	gl.AttachShader(program, i.shader)
	return nil
}

func readShader(filename string) (string, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	source := fmt.Sprintf("%s\x00", buf.String())

	return source, nil
}

// Delete the shader
func (i *Info) Delete() {
	gl.DeleteShader(i.shader)
}

// cleanup all shaders by calling Delete on any non-zero shader in the slice.
func cleanup(shaders *[]Info) {
	for _, s := range *shaders {
		if s.shader != 0 {
			s.Delete()
		}
	}
}

// getErrorsMsg helps to return an error message when compiling/linking goes
// wrong.  If shader is true, check for logs relating to a shader failing to
// compile.  In this context id should be the ID of the shader that failed to
// compile.  If shader is false, check for logs relating to linking shaters
// to a program.  in this context, id is the ID of the program that was
// attempting to link the shaders.
func getErrorMsg(shader bool, id uint32) string {
	var l int32
	if shader {
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &l)
	} else {
		gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &l)
	}

	msg := strings.Repeat("\x00", int(l+1))

	if shader {
		gl.GetShaderInfoLog(id, l, nil, gl.Str(msg))
	} else {
		gl.GetProgramInfoLog(id, l, nil, gl.Str(msg))
	}

	return msg
}
