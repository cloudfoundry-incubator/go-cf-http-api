package testhelpers

import (
	"sync"
)

type FakeBuffer struct {
	lock sync.Mutex
	logs []string
}

func NewFakeBuffer() *FakeBuffer {
	return &FakeBuffer{
		logs: make([]string, 0, 10),
	}
}

func (f *FakeBuffer) Write(p []byte) (n int, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.logs = append(f.logs, string(p))
	return len(p), nil
}

func (f *FakeBuffer) GetContent() []string {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.logs
}
