package mooze

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Response *http.Response

type methodtype int

const (
	GET methodtype = iota
	POST
	PUT
	PATCH
	DELETE
)

func methodTypeToString(m methodtype) string {
	return []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
	}[m]
}

type MoozeRequest struct {
	url    string
	method methodtype
	header string        // temp
	body   *bytes.Buffer // temp

	resStatus string
	resCode   int
	data      []string
}

func NewMoozeRequest() *MoozeRequest {
	return &MoozeRequest{
		url:       "",
		method:    GET,
		header:    "",
		body:      nil,
		resStatus: "",
		resCode:   -1,
	}
}

func (r *MoozeRequest) Url(u string) {
	r.url = u
}

func (r *MoozeRequest) Method(m int) {
	r.method = methodtype(m)
}

type ReqArgs struct {
	h   string
	buf *bytes.Buffer
}

func (r *MoozeRequest) Send(m string, args ReqArgs) *http.Response {
	switch m {
	case "GET":
		res, err := http.Get(r.url)
		if err != nil {
			panic(err)
		}
		return res

	case "POST":
		res, err := http.Post(r.url, args.h, args.buf)
		if err != nil {
			panic(err)
		}
		return res

	default:
		res, err := http.Get(r.url)
		if err != nil {
			panic(err)
		}
		return res
	}
}

func (r *MoozeRequest) Body(res *http.Response) []byte {
	b, _ := ioutil.ReadAll(res.Body)
	return b
}

func (r *MoozeRequest) Prettify(data []byte) []string {
	j := &bytes.Buffer{}
	err := json.Indent(j, data, "", "  ")
	// response is not a valid json
	if err != nil {
		return []string{string(data)}
	}

	str := []string{}
	buf := ""
	brd := bytes.NewReader(j.Bytes())
	rrd := bufio.NewReader(brd)
	for {
		b, _, err := rrd.ReadRune()
		if err == io.EOF {
			str = append(str, buf)
			break
		}
		buf += string(b)
		if b == '\n' {
			str = append(str, buf)
			buf = ""
		}
	}
	return str
}

func (r *MoozeRequest) ParseJson(s string) *bytes.Buffer {
	b := []byte(s)
	j := &bytes.Buffer{}
	err := json.Indent(j, b, "", "  ")
	if err != nil {
		panic(err)
	}
	return j
}
