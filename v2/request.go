package v2

import (
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
	Url    string
	Method string
	Client *http.Client
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

func (r *Request) Send() []byte {
	req, err := http.NewRequest(HttpMethod.GET, r.Url, nil)
	// req, err := http.NewRequest("GET", "https://api.github.com/users/RudyPark3091", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := r.Client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return body
}
