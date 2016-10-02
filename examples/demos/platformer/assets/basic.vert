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

#version 410

uniform mat4 modelMatrix;
uniform mat4 projectionMatrix;

layout (location = 0) in vec4 mcVertex;
layout (location = 1) in vec4 mcColor;

out vec4 vsColor;

void main(void)
{
    vsColor = mcColor;
    gl_Position = projectionMatrix * (modelMatrix * mcVertex);
}
