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
// Package display TODO doc

package display

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/convexbit/shade/events"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Context TODO doc
type Context struct {
	Window  *glfw.Window
	Width   float32
	Height  float32
	Program uint32
}

// Signal to close the window
func (c *Context) Close() {
	c.Window.SetShouldClose(true)
	events.WindowCloseCallback(c.Window)
}

func createWindow(major, minor int, title string) (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)
	if major != 2 && minor != 1 {
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}
	return glfw.CreateWindow(640, 480, title, nil, nil)
}

// SetMode TODO doc
func SetMode(title string, width, height int) (*Context, error) {
	c := Context{
		Width:  float32(width),
		Height: float32(height),
	}
	if err := glfw.Init(); err != nil {
		return &c, fmt.Errorf("failed to initialize glfw: %v", err)
	}
	// TODO: move this to a terminate function
	//defer glfw.Terminate()

	supportedVersions := [][2]int{
		// TODO: create matching shaders for these versions
		//[2]int{4, 1},
		//[2]int{3, 3},
		[2]int{2, 1},
	}

	var err error
	var window *glfw.Window
	for _, v := range supportedVersions {
		window, err = createWindow(v[0], v[1], title)
		if err == nil {
			// Successfully created window, break out of loop
			break
		}
		// TODO: maybe only print this w/ a verbose logging option
		fmt.Println("Warning:", err)
	}

	c.Window = window

	c.Window.MakeContextCurrent()
	c.Window.SetKeyCallback(events.KeyCallback)
	c.Window.SetMouseButtonCallback(events.MouseButtonCallback)
	c.Window.SetCursorPosCallback(events.CursorPositionCallback)
	c.Window.SetCloseCallback(events.WindowCloseCallback)

	if err := gl.Init(); err != nil {
		return &c, fmt.Errorf("failed to init glow: %v", err)
	}

	fmt.Println("OpenGL vendor", gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println("OpenGL renderer", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("GLSL version", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	c.Program, err = newProgram(vertexShader, fragmentShader)
	if err != nil {
		return &c, fmt.Errorf("error loading program: %v", err)
	}

	gl.UseProgram(c.Program)

	gl.BindFragDataLocation(c.Program, 0, gl.Str("FragColor\x00"))

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	// TODO: Figure out why "layering" using z-buffer does not work.
	// gl.DepthFunc(gl.LESS)
	// gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.NEVER)
	gl.Enable(gl.CULL_FACE)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	//window.SwapBuffers()
	return &c, nil
}

// Fill TODO doc
func (c *Context) Fill(r, g, b float32) {
	gl.ClearColor(r, g, b, 1.0)
}

// Flip TODO doc
func (c *Context) Flip() {
	c.Window.SwapBuffers()
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("can not create vert shader: %s", err)
	}

	fragShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("can not create frag shader: %s", err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertShader)
	gl.DeleteShader(fragShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var vertexShader = `
#version 120

uniform mat4 ProjMatrix;
uniform mat4 ViewMatrix;
uniform mat4 ModelMatrix;
uniform mat3 TexMatrix;
uniform vec3 LightPos;

attribute vec3 MCVertex;
attribute vec3 MCNormal;
attribute vec3 MCTangent;
attribute vec2 TexCoord0;

varying vec2 TexCoord;
varying vec3 Pos;
varying vec3 LightDir;
varying vec3 EyeDir;


mat3 inv(mat3 m) {
	// TDOO implement inverse func
	return m;
}

void main() {
  mat4 mvMatrix = ViewMatrix * ModelMatrix;
  vec4 ccVertex = mvMatrix * vec4(MCVertex, 1.0);
  gl_Position = ProjMatrix * ccVertex;
  Pos = vec4(ModelMatrix * vec4(MCVertex, 1.0)).xyz;

  TexCoord = vec3(TexMatrix * vec3(TexCoord0, 1.0)).st;

  mat3 normalMatrix = mat3x3(mvMatrix);
  normalMatrix = inv(normalMatrix);
  normalMatrix = transpose(normalMatrix);

  mat3 mv3Matrix = mat3(mvMatrix);
  vec3 n = normalize(MCNormal); // TODO: fix normalize(mv3Matrix * MCNormal);
  vec3 t = normalize(mv3Matrix * MCTangent);
  vec3 b = normalize(mv3Matrix * cross(n, t));

  LightDir = vec3(ViewMatrix * vec4(LightPos, 0.0)) - vec3(ccVertex);
  vec3 v;
  v.x = dot(LightDir, t);
  v.y = dot(LightDir, b);
  v.z = dot(LightDir, n);
  LightDir = v;

  EyeDir = vec3(-ccVertex);
  v.x = dot(EyeDir, t);
  v.y = dot(EyeDir, b);
  v.z = dot(EyeDir, n);
  EyeDir = v;
}
` + "\x00"

var fragmentShader = `
#version 120

uniform int AddColor;
uniform vec4 AColor;
uniform int SubColor;
uniform vec4 SColor;
uniform sampler2D ColorMap;
uniform sampler2D NormalMap;
uniform vec4 AmbientColor;
uniform vec3 LightPos;
uniform vec4 LightColor;
uniform float LightPower;

varying vec3 Pos;
varying vec3 LightDir;
varying vec3 EyeDir;
varying vec2 TexCoord;

// out vec4 FragColor;

void main() {
  float alpha = texture2D(ColorMap, TexCoord.st).a;
  vec3 diffuse = texture2D(ColorMap, TexCoord.st).rgb;
  if (AddColor == 1) {
    diffuse = clamp(diffuse + AColor.rgb, 0.0, 1.0);
  }
  if (SubColor == 1) {
    diffuse = clamp(diffuse - SColor.rgb, 0.0, 1.0);
  }
  vec3 ambient = AmbientColor.rgb * diffuse;
  vec3 specular = diffuse/8;

  vec3 normal = texture2D(NormalMap, TexCoord.st).rgb * 2 - 1;
  float distance = length(LightPos - Pos);

  vec3 n = normalize(normal);
  vec3 l = normalize(LightDir);

  float cosTheta = clamp(dot(n, l), 0.0, 1.0);

  vec3 e = normalize(EyeDir);
  vec3 r = reflect(-l, n);

  float cosAlpha = clamp(dot(e, r), 0.0, 1.0);

  gl_FragColor = vec4(
    ambient +
    diffuse * LightColor.rgb * LightPower * cosTheta /
      (distance * distance) +
    specular * LightColor.rgb * LightPower * pow(cosAlpha, 5) /
      (distance * distance), alpha);
}
` + "\x00"
