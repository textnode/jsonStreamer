package jsonStreamer

import (
	"io"
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
	_, err = self.writer.Write([]byte(value))
	if nil != err {
		return
	}
	_, err = self.writer.Write(DOUBLE_QUOTE)
	return
}

func (self *JsonStreamer) WriteStringValueBytes(value []byte) (err error) {
	self.lastWasKey = false
	_, err = self.writer.Write(DOUBLE_QUOTE)
	if nil != err {
		return
	}
	_, err = self.writer.Write(value)
	if nil != err {
		return
	}
	_, err = self.writer.Write(DOUBLE_QUOTE)
	return
}
