package server

import (
	"application_metadata_api_server/cache/mocks"
	"application_metadata_api_server/server/api"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpServerImpl_PutHandler(t *testing.T) {
	mockStore := &mocks.Store{}
	fakeServer := &httpServerImpl{
		store:     mockStore,
		validator: newAppValidator(),
	}
	// mock on
	mockStore.On("Add", mock.Anything, mock.Anything).Return(api.Id("1"), nil)

	testCases := []struct {
		name                 string
		filePath             string
		expectedResponseCode int
	}{
		{
			name:                 "expect 400 on invalid input 1",
			filePath:             "../testdata/invalid-payload1.yaml",
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "expect 400 on invalid input 2",
			filePath:             "../testdata/invalid-payload2.yaml",
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "expect 200 on valid input",
			filePath:             "../testdata/valid-payload1.yaml",
			expectedResponseCode: http.StatusOK,
		},
		{
			name:                 "expect 200 on valid input",
			filePath:             "../testdata/valid-payload2.yaml",
			expectedResponseCode: http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(test.filePath)
			assert.Nil(t, err)
			req := httptest.NewRequest("POST", "/put", bytes.NewReader(data))
			w := httptest.NewRecorder()
			fakeServer.PutHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, test.expectedResponseCode, resp.StatusCode)
			t.Log(resp.Header)
			t.Log(string(body))
		})
	}
}

func TestHttpServerImpl_GetHandler(t *testing.T) {
	mockStore := &mocks.Store{}
	fakeServer := &httpServerImpl{
		store:     mockStore,
		validator: newAppValidator(),
	}
	data, err := ioutil.ReadFile("../testdata/valid-payload1.yaml")
	assert.Nil(t, err)
	mockStore.On("Get", api.Id("1")).Return(data, nil)
	mockStore.On("Get", api.Id("2")).Return(nil, fmt.Errorf("not found"))

	req := httptest.NewRequest("GET", "/get", strings.NewReader("1"))
	w := httptest.NewRecorder()
	fakeServer.GetHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	t.Log(string(body))

	req = httptest.NewRequest("GET", "/get", strings.NewReader("2"))
	w = httptest.NewRecorder()
	fakeServer.GetHandler(w, req)

	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	t.Log(string(body))
}

func TestHttpServerImpl_SearchHandler(t *testing.T) {
	mockStore := &mocks.Store{}
	fakeServer := &httpServerImpl{
		store:     mockStore,
		validator: newAppValidator(),
	}
	mockStore.On("SearchStruct", mock.Anything).Return([]api.Id{"1"}, nil)

	testCases := []struct {
		name                 string
		filePath             string
		expectedResponseCode int
	}{
		{
			name:                 "expect 200 on valid input 1",
			filePath:             "../testdata/query1.yaml",
			expectedResponseCode: http.StatusOK,
		},
		{
			name:                 "expect 200 on valid input 2",
			filePath:             "../testdata/query2.yaml",
			expectedResponseCode: http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(test.filePath)
			assert.Nil(t, err)
			req := httptest.NewRequest("POST", "/query", bytes.NewReader(data))
			w := httptest.NewRecorder()
			fakeServer.SearchHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, test.expectedResponseCode, resp.StatusCode)
			t.Log(string(body))
		})
	}

}
