package api

import (
	"errors"
	"net/http"
	"reflect"

	"encoding/json"

	"io/ioutil"

	"../uaaclient"
	"github.com/gorilla/mux"
)

type Request interface {
	GetParam(name string) string
	CurrentUser() *uaaclient.User
	Decode(value interface{}) error
	RawBody() []byte
	Path() string
}

type realRequest struct {
	httpRequest *http.Request
	bodyRead    bool
	body        []byte
	currentUser *uaaclient.User
}

func (r *realRequest) GetParam(n string) string {
	v, ok := mux.Vars(r.httpRequest)[n]
	if ok {
		return v
	}

	return r.httpRequest.URL.Query().Get(n)
}

func (r *realRequest) CurrentUser() *uaaclient.User {
	return r.currentUser
}

func (r *realRequest) Decode(target interface{}) error {
	err := json.Unmarshal(r.RawBody(), target)
	if err != nil {
		return err
	}

	r.httpRequest.Body.Close()

	return nil
}

func (r *realRequest) RawBody() []byte {
	if r.bodyRead {
		return r.body
	} else {
		r.bodyRead = true

		b, err := ioutil.ReadAll(r.httpRequest.Body)
		if err != nil {
			return []byte{}
		}

		r.body = b
		r.httpRequest.Body.Close()
		return r.body
	}
}

func (r *realRequest) Path() string {
	return r.httpRequest.URL.Path
}

type FakeRequest struct {
	User          uaaclient.User
	Params        map[string]string
	Body          interface{}
	ErrorOnDecode bool
}

func (f *FakeRequest) GetParam(n string) string {
	return f.Params[n]
}

func (f *FakeRequest) CurrentUser() *uaaclient.User {
	return &f.User
}

func (f *FakeRequest) Decode(target interface{}) error {
	if f.ErrorOnDecode || f.Body == nil {
		return errors.New("decode error")
	}

	v := reflect.ValueOf(target).Elem()
	v.Set(reflect.ValueOf(f.Body))

	return nil
}

func (f *FakeRequest) RawBody() []byte {
	return []byte("{}")
}

func (f *FakeRequest) Path() string {
	return "/"
}
