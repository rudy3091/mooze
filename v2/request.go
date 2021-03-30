package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type httpMethodStruct struct {
	GET  string
	POST string
}

var HttpMethod = httpMethodStruct{
	GET:  "GET",
	POST: "POST",
}

type Request struct {
	Url        string
	Method     string
	Body       string
	BodyBuffer *bytes.Buffer
	Client     *http.Client
}

func NewRequest() *Request {
	r := &Request{
		Url:    "",
		Method: HttpMethod.GET,
	}
	r.SetHttpClient()
	return r
}

func (r *Request) SetHttpClient() {
	r.Client = &http.Client{}
}

// Json() prettifies json byte data and returns string
func (r *Request) Json(data []byte) string {
	j := &bytes.Buffer{}
	err := json.Indent(j, data, "", "  ")
	// invalid json format
	if err != nil {
		return string(data)
	}
	return string(j.Bytes())
}

func (r *Request) ParseJson(s string) *bytes.Buffer {
	b := []byte(s)
	j := &bytes.Buffer{}
	err := json.Indent(j, b, "", "  ")
	// not valid json
	if err != nil {
		return bytes.NewBufferString(s)
	}
	return j
}

func (r *Request) Send() ([]byte, string, error) {
	req, err := (func() (*http.Request, error) {
		if r.Method == HttpMethod.GET {
			req, err := http.NewRequest(r.Method, r.Url, nil)
			return req, err
		} else {
			req, err := http.NewRequest(r.Method, r.Url, bytes.NewBufferString(r.Body))
			return req, err
		}
	})()
	if err != nil {
		return nil, "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := r.Client.Do(req)

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
