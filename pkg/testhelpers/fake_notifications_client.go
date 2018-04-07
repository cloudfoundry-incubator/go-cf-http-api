package testhelpers

import (
	"net/http"
	"sync"

	notificationsapi "github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"
)

type FakeNotificationsClient struct {
	mu           sync.Mutex
	sendError    error
	token        string
	address      string
	status       int
	notification notificationsapi.Notify
	responses    []notificationsapi.NotifyResponse
}

func NewFakeNotificationsClient() *FakeNotificationsClient {
	return &FakeNotificationsClient{
		responses: make([]notificationsapi.NotifyResponse, 0),
		status:    http.StatusAccepted,
	}
}

func (f *FakeNotificationsClient) Email(token, address string, notify notificationsapi.Notify) (int, []notificationsapi.NotifyResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.token = token
	f.address = address
	f.notification = notify

	return f.status, f.responses, f.sendError
}

func (f *FakeNotificationsClient) LastAddress() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.address
}

func (f *FakeNotificationsClient) LastToken() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.token
}

func (f *FakeNotificationsClient) LastNotification() notificationsapi.Notify {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.notification
}

func (f *FakeNotificationsClient) SetSendStatus(status int) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.status = status
}

func (f *FakeNotificationsClient) SetSendError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.sendError = err
}

func (f *FakeNotificationsClient) AddNotificationResponse(response notificationsapi.NotifyResponse) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.responses = append(f.responses, response)
}
