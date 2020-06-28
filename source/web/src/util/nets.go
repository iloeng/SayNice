package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var myClient *http.Client

func init() {
	myClient = &http.Client{Timeout: 10 * time.Second}
}

// GetHTTPClient 获取 http.Client 对象
func GetHTTPClient() *http.Client {
	return myClient
}

// GetJSON Http Get 请求，返回一个 json 对象
func GetJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// err = json.NewDecoder(r.Body).Decode(target)

	// see https://stackoverflow.com/questions/24111888/decoding-a-request-body-in-go-why-am-i-getting-an-eof
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	// fmt.Println(buf.String())
	err = json.Unmarshal(buf.Bytes(), target)

	return err
}

// PostJSON http post 请求，请求体为 json 对象，返回一个 json 对象
func PostJSON(url string, body interface{}, target interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	// fmt.Println(string(b))
	r, err := myClient.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), target)

	return err
}

// Post http post 请求，请求体为 From 表单，返回字符串
func Post(url string, form url.Values) (string, error) {
	body := strings.NewReader(form.Encode())
	resp, err := myClient.Post(url, "application/x-www-form-urlencoded", body)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return string(buf.String()), nil
}
