package mock

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HTTPResponseMock struct {
	StatusCode int
	Body       interface{}
}

type RoundTripperMock struct {
	URLMap map[string]HTTPResponseMock
}

func NewRoundTripperMock() *RoundTripperMock {
	return &RoundTripperMock{
		URLMap: make(map[string]HTTPResponseMock),
	}
}

func (m *RoundTripperMock) SetResponse(url string, response HTTPResponseMock) {
	m.URLMap[url] = response
}

func (m *RoundTripperMock) RoundTrip(req *http.Request) (*http.Response, error) {
	if mockResponse, found := m.URLMap[req.URL.Path]; found {
		bodyBytes, _ := json.Marshal(mockResponse.Body)
		return &http.Response{
			StatusCode: mockResponse.StatusCode,
			Body:       ioutil.NopCloser(bytes.NewBuffer(bodyBytes)),
		}, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte("not found"))),
	}, nil
}
