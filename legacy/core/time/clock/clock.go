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
// Package clock TODO doc

package clock

import "time"

type Clock struct {
	lastTime  time.Time
	firstCall bool
	count     int
	tmp       float32
}

func New() (*Clock, error) {
	c := Clock{
		firstCall: true,
	}

	return &c, nil
}

// Tick updates the clock, ensuring that it is called at most limit times per second and returns the number of milliseconds since the last call.
//
// If limit is set to 0, then there should be no restriction on the number of calls per second.
func (c *Clock) Tick(limit int) float32 {
	if c.count >= limit || c.tmp >= float32(time.Second/time.Millisecond) {
		if c.tmp <= float32(time.Second/time.Millisecond) {
			// TODO: fix sleep
			//time.Sleep(time.Second - (time.Millisecond * time.Duration(c.tmp)))
		}
		c.tmp = 0
		c.count = 0
	}

	t := time.Now()
	if c.firstCall {
		c.firstCall = false
		c.lastTime = t
	}

	d := float32(t.Sub(c.lastTime)) / float32(time.Millisecond)

	if limit >= 0 {
		c.count += 1
		c.tmp += d
	}

	c.lastTime = t
	return d
}
