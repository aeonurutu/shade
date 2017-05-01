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
// package debug TODO: doc
package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

func pushState(ws *websocket.Conn, eng *Engine) {
	for {
		data := eng.State()
		ws.Write(data)
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func DebugServer(ws *websocket.Conn, eng *Engine) {
	//fmt.Println("[DebugServer]")
	defer ws.Close()

	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		panic(err)
	}

	type Message struct {
		Method string `json:"method"`
	}

	dec := json.NewDecoder(strings.NewReader(string(msg[:n])))

	var m Message
	for {

		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	pushState(ws, eng)
}

// SpawnServer to serve debug client and websocket to relay engine state.
func SpawnServer(eng *Engine) {
	fmt.Printf("Spawning debug server on 10888. Go to http://127.0.0.1:10888/\n")

	// Spawn HTTP server
	go func() {
		handler := http.FileServer(http.Dir("assets/debug_server"))
		http.Handle("/", handler)
		err := http.ListenAndServe(":10888", nil)
		if err != nil {
			panic(err)
		}
	}()

	// Spawn websocket server
	go func() {
		http.Handle("/debug", websocket.Handler(func(ws *websocket.Conn) {
			DebugServer(ws, eng)
		}))
		err := http.ListenAndServe(":10889", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
}
