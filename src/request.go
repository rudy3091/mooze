// TODO: single json object will be parsed,
//       but not array of json objects. fix this
package mooze

import (
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
	header string // temp
	body   string // temp

	resStatus string
	resCode   int
}

func NewMoozeRequest() *MoozeRequest {
	return &MoozeRequest{
		url:       "",
		method:    GET,
		header:    "",
		body:      "",
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

func (r *MoozeRequest) Send() *http.Response {
	res, err := http.Get(r.url)
	if err != nil {
		panic(err)
	}

	return res
}

func (r *MoozeRequest) Body(res *http.Response) []byte {
	b, _ := ioutil.ReadAll(res.Body)
	return b
}

func (r *MoozeRequest) Prettify(data []byte) []string {
	j := &bytes.Buffer{}
	json.Indent(j, data, "", "  ")

	str := []string{}
	buf := ""
	rd := bytes.NewReader(j.Bytes())
	for {
		b, err := rd.ReadByte()
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
