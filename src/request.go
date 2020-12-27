package mooze

import (
	"io/ioutil"
	"net/http"
)

type methodtype int

const (
	GET methodtype = iota
	POST
	PUT
	DELETE
	PATCH
)

type Request struct {
	url    string
	method methodtype
	header string // temp
	body   string // temp
}

func NewRequest(url string, method methodtype, header, body string) *Request {
	return &Request{url, method, header, body}
}

func (r *Request) Send() *http.Response {
	res, err := http.Get(r.url)
	if err != nil {
		panic(err)
	}

	return res
}

func (r *Request) Body(res *http.Response) []byte {
	b, _ := ioutil.ReadAll(res.Body)
	return b
}
