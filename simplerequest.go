package simplerequest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const InternalServerError = 500
const OK = 200

var lastErr error

type BaseResponse interface{}

func BaseRequestWithHeader(method, u string, h map[string][]string, payload string) (int, map[string][]string, string) {
	var hr map[string][]string
	request, _ := http.NewRequest(method, u, strings.NewReader(payload))
	request.Header = h
	client := http.Client{}

	resp, e := client.Do(request)
	if e != nil {
		lastErr = e
		return InternalServerError, hr, "Fail to do the requisition"
	}
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		lastErr = e
	}
	return resp.StatusCode, resp.Header, string(body)
}

func Get(u string) (int, string) {
	resp, e := http.Get(u)
	if e != nil {
		lastErr = e
		return InternalServerError, "Fail to do the requisition"
	}
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		lastErr = e
	}
	return resp.StatusCode, string(body)
}

func GetH(u string, h map[string]string) (int, string, map[string]string) {
	var hr map[string]string
	hToS := headerToMulti(h)
	request, _ := http.NewRequest(http.MethodGet, u, nil)
	request.Header = hToS
	client := http.Client{}

	resp, e := client.Do(request)
	if e != nil {
		lastErr = e
		return InternalServerError, "Fail to do the requisition", hr
	}
	hr = headerToUnified(resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		lastErr = e
	}
	return resp.StatusCode, string(body), hr
}

func PostH(u string, h map[string][]string, payload string) (int, map[string][]string, string) {
	return BaseRequestWithHeader(http.MethodPost, u, h, payload)
}

func PutH(u string, h map[string][]string, payload string) (int, map[string][]string, string) {
	return BaseRequestWithHeader(http.MethodPut, u, h, payload)
}

func Put(u string, payload string) (int, string) {
	request, _ := http.NewRequest(http.MethodPut, u, strings.NewReader(payload))
	client := http.Client{}

	resp, e := client.Do(request)
	if e != nil {
		lastErr = e
		return InternalServerError, "Fail to do the requisition"
	}
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		lastErr = e
	}
	return resp.StatusCode, string(body)
}

func Post(u string, body map[string][]string) (int, string) {
	data := url.Values(body)
	resp, err := http.PostForm(u, data)
	if err != nil {
		lastErr = err
	}
	if resp == nil {
		return 500, ""
	}
	res, e := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if e != nil {
		lastErr = e
	}
	return resp.StatusCode, string(res)
}

//url and the struct with return pattern
func GetToJson(u string, response BaseResponse) int {
	c, r := Get(u)
	b := []byte(r)
	e := json.Unmarshal(b, &response)
	if e != nil {
		lastErr = e
	}
	return c
}

func headerToUnified(h map[string][]string) map[string]string {
	hr := make(map[string]string, len(h))
	for key, header := range h {
		hr[key] = header[0]
	}
	return hr
}

func headerToMulti(h map[string]string) map[string][]string {
	hr := make(map[string][]string, len(h))
	for key, header := range h {
		hr[key][0] = header
	}
	return hr
}
