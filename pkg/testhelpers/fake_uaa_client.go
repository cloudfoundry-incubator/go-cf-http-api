package testhelpers

import (
	"context"

	"sync"

	"../uaaclient"
)

type fakeUAAClient struct {
	mu   sync.Mutex
	user *uaaclient.User
	err  error
}

func NewFakeUAAClient() *fakeUAAClient {
	return &fakeUAAClient{}
}

func (f *fakeUAAClient) CheckToken(token string, ctx context.Context) (*uaaclient.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.user, f.err
}

func (f *fakeUAAClient) SetUser(user *uaaclient.User) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.user = user
}

func (f *fakeUAAClient) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.err = err
}
