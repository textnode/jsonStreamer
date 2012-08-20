// Copyright 2012 Darren Elwood <darren@textnode.com> http://www.textnode.com @textnode
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonStreamer

import (
	"io"
	"strings"
)

// Ideally these would be consts, but I can't define them as consts when using chars for readability
var START_OBJECT []byte = []byte{'{'}
var END_OBJECT []byte = []byte{'}'}
var START_ARRAY []byte = []byte{'['}
var END_ARRAY []byte = []byte{']'}
var COLON []byte = []byte{':'}
var DOUBLE_QUOTE []byte = []byte{'"'}
var FIELD_SEPARATOR []byte = []byte{','}

// Could just be a bool, but there might be edge cases I haven't discovered and hence extra fields could be needed in future
type Frame struct {
	havePrevious bool
}

func NewFrame() *Frame {
	return &Frame{}
}

type JsonStreamer struct {
	writer     io.Writer
	frames     []*Frame
	lastWasKey bool
}

func NewJsonStreamer(writer io.Writer) *JsonStreamer {
	var newJsonStreamer *JsonStreamer = &JsonStreamer{writer: writer, frames: make([]*Frame, 1, 10)}
	newJsonStreamer.frames[0] = NewFrame()
	return newJsonStreamer
}

// When a field separator is needed write one.
func (self *JsonStreamer) separate() (err error) {
	var index int = len(self.frames) - 1
	if self.frames[index].havePrevious && !self.lastWasKey {
		_, err = self.writer.Write(FIELD_SEPARATOR)
		self.lastWasKey = false
	} else {
		self.frames[index].havePrevious = true
	}
	return
}

func (self *JsonStreamer) StartObject() (err error) {
	err = self.separate()
	self.lastWasKey = false
	_, err = self.writer.Write(START_OBJECT)
	self.frames = append(self.frames, NewFrame())
	return
}

func (self *JsonStreamer) EndObject() (err error) {
	self.lastWasKey = false
	_, err = self.writer.Write(END_OBJECT)
	self.frames = self.frames[:len(self.frames)-1]
	return
}

func (self *JsonStreamer) StartArray() (err error) {
	err = self.separate()
	self.lastWasKey = false
	_, err = self.writer.Write(START_ARRAY)
	self.frames = append(self.frames, NewFrame())
	return
}

func (self *JsonStreamer) EndArray() (err error) {
	self.lastWasKey = false
	_, err = self.writer.Write(END_ARRAY)
	self.frames = self.frames[:len(self.frames)-1]
	return
}

func (self *JsonStreamer) WriteKey(name string) (err error) {
	err = self.separate()
	self.lastWasKey = true
	_, err = self.writer.Write(DOUBLE_QUOTE)
	if nil != err {
		return
	}
	_, err = self.writer.Write([]byte(name))
	if nil != err {
		return
	}
	_, err = self.writer.Write(DOUBLE_QUOTE)
	if nil != err {
		return
	}
	_, err = self.writer.Write(COLON)
	return
}

func (self *JsonStreamer) WriteStringValue(value string) (err error) {
	self.lastWasKey = false
	_, err = self.writer.Write(DOUBLE_QUOTE)
	if nil != err {
		return
	}

	value = self.escapeStringValue(value)

	_, err = self.writer.Write([]byte(value))
	if nil != err {
		return
	}
	_, err = self.writer.Write(DOUBLE_QUOTE)
	return
}

/*
func (self *JsonStreamer) WriteStringValueBytes(value []byte) (err error) {
	self.lastWasKey = false
	_, err = self.writer.Write(DOUBLE_QUOTE)
	if nil != err {
		return
	}

	value = self.escapeStringValueBytes(value)

	_, err = self.writer.Write(value)
	if nil != err {
		return
	}
	_, err = self.writer.Write(DOUBLE_QUOTE)
	return
}
*/

func (self *JsonStreamer) escapeStringValue(value string) (result string) {
	result = value
	result = strings.Replace(result, "<", "\u003c", -1)
	result = strings.Replace(result, ">", "\u003e", -1)
	result = strings.Replace(result, "\b", "\\\b", -1)
	result = strings.Replace(result, "\f", "\\\f", -1)
	result = strings.Replace(result, "\n", "\\\n", -1)
	result = strings.Replace(result, "\r", "\\\r", -1)
	result = strings.Replace(result, "\t", "\\\t", -1)
	result = strings.Replace(result, "\v", "\\\v", -1)
	result = strings.Replace(result, "'", "\\'", -1)
	result = strings.Replace(result, "\"", "\\\"", -1)
	result = strings.Replace(result, "\\", "\\\\", -1)
	return
}

/*
func (self *JsonStreamer) escapeStringValueBytes(value []byte) (result []byte) {
    value = bytes.Replace(value, []byte{"<"}, []byte{"\u003c"}, -1)
    value = bytes.Replace(value, []byte{">"}, []byte{"\u003e"}, -1)
	value = bytes.Replace(value, []byte{"\b"}, []byte{"\\b"}, -1)
	value = bytes.Replace(value, []byte{"\f"}, []byte{"\\f"}, -1)
	value = bytes.Replace(value, []byte{"\n"}, []byte{"\\n"}, -1)
	value = bytes.Replace(value, []byte{"\r"}, []byte{"\\r"}, -1)
	value = bytes.Replace(value, []byte{"\t"}, []byte{"\\t"}, -1)
	value = bytes.Replace(value, []byte{"\v"}, []byte{"\\v"}, -1)
	value = bytes.Replace(value, []byte{"\'"},  []byte{"\\'"}, -1)
	value = bytes.Replace(value, []byte{"\'"}, []byte{"\\\""}, -1)
	value = bytes.Replace(value, []byte{"\\"}, []byte{"\\\\"}, -1)
	return
}
*/


