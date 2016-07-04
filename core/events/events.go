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
// Package events manages events

package events

import (
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Handler for events
type Handler interface {
	Handle(event Event)
}

const (
	CursorPosition  = iota
	KeyDown         = iota
	KeyUp           = iota
	KeyRepeat       = iota
	WindowClose     = iota
	MouseButtonDown = iota
	MouseButtonUp   = iota
)

type Event struct {
	Type        int
	Window      *glfw.Window
	Key         glfw.Key
	Scancode    int
	Mods        glfw.ModifierKey
	X           float32
	Y           float32
	MouseButton glfw.MouseButton
}

var events []Event

func Get() []Event {
	// TODO: This might cause a lot of garbage collection, which is prob bad.
	var elist []Event
	for i := range events {
		elist = append(elist, Event{
			Type:        events[i].Type,
			Window:      events[i].Window,
			Key:         events[i].Key,
			Scancode:    events[i].Scancode,
			Mods:        events[i].Mods,
			X:           events[i].X,
			Y:           events[i].Y,
			MouseButton: events[i].MouseButton,
		})
	}
	events = nil
	return elist
}

// CursorPositionCallback TODO doc
func CursorPositionCallback(w *glfw.Window, x, y float64) {
	_, h := w.GetSize()
	// TODO: these are from the top/left should be bottom/left to match sprite drawing
	events = append(events, Event{
		Type:   CursorPosition,
		Window: w,
		X:      float32(x),
		Y:      float32(h) - float32(y),
	})
}

// KeyCallback TODO doc
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	e := Event{
		Window:   w,
		Key:      key,
		Scancode: scancode,
		Mods:     mods,
	}
	switch action {
	case glfw.Press:
		e.Type = KeyDown
	case glfw.Release:
		e.Type = KeyUp
	case glfw.Repeat:
		e.Type = KeyRepeat
	}
	events = append(events, e)
}

// WindowCloseCallback TODO doc
func WindowCloseCallback(w *glfw.Window) {
	events = append(events, Event{
		Type:   WindowClose,
		Window: w,
	})
}

// MouseButtonCallback  TODO: doc
func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	e := Event{
		Window:      w,
		MouseButton: button,
		Mods:        mods,
	}
	switch action {
	case glfw.Press:
		e.Type = MouseButtonUp
	case glfw.Release:
		e.Type = MouseButtonDown
	}
	events = append(events, e)
}

func Poll() {
	glfw.PollEvents()
}
