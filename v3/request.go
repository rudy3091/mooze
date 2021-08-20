package v3

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Request struct {
	url        string
	method     string
	body       string
	bodyBuffer *bytes.Buffer
	headers    map[string]string
	client     *http.Client
}

var RequestInfo *Request

func NewRequest() *Request {
	RequestInfo = &Request{
		url:     "https://api.github.com/users/rudy3091",
		method:  "GET",
		headers: map[string]string{},
	}
	RequestInfo.SetHttpClient()
	return RequestInfo
}

func (r *Request) SetHttpClient() {
	r.client = &http.Client{}
}

func (r *Request) Json(data []byte) string {
	j := &bytes.Buffer{}
	err := json.Indent(j, data, "", "  ")
	if err != nil {
		return string(data)
	}
	return string(j.Bytes())
}

func (r *Request) ParseJson(s string) *bytes.Buffer {
	b := []byte(s)
	j := &bytes.Buffer{}
	err := json.Indent(j, b, "", "  ")
	if err != nil {
		return bytes.NewBufferString(s)
	}
	return j
}

func (r *Request) ParseHeaders() string {
	s := ""
	if len(r.headers) == 0 {
		return "No Headers"
	} else {
		for k, v := range r.headers {
			s += k + " : " + v
		}
		return s
	}
}

func (r *Request) ParseHeadersOptions() []string {
	s := []string{"+ Add new"}
	if len(r.headers) != 0 {
		for k, v := range r.headers {
			s = append(s, "- "+k+": "+v)
		}
	}
	return s
}

// func (r *Request) Send() ([]byte, string, error) {
// 	req, err := (func() (*http.Request, error) {
// 		if r.method == "GET" {
// 			req, err := http.NewRequest(r.method, r.url, nil)
// 			return req, err
// 		} else {
// 			req, err := http.NewRequest(r.method, r.url, bytes.NewBufferString(r.body))
// 			return req, err
// 		}
// 	})()
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	for k, v := range r.headers {
// 		req.Header.Add(k, v)
// 	}

// 	res, err := r.client.Do(req)

// 	if err != nil {
// 		return nil, "", err
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	return body, res.Status, nil
// }

func (r *Request) Send() ([]byte, string, error) {
	req, err := (func() (*http.Request, error) {
		if RequestInfo.method == "GET" {
			req, err := http.NewRequest(RequestInfo.method, RequestInfo.url, nil)
			return req, err
		} else {
			req, err := http.NewRequest(RequestInfo.method, RequestInfo.url, bytes.NewBufferString(RequestInfo.body))
			return req, err
		}
	})()
	if err != nil {
		return nil, "", err
	}
	for k, v := range RequestInfo.headers {
		req.Header.Add(k, v)
	}

	res, err := RequestInfo.client.Do(req)

	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}
	return body, res.Status, nil
}
